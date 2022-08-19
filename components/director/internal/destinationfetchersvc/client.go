package destinationfetchersvc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/config"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// DestinationServiceAPIConfig destination service api configuration
type DestinationServiceAPIConfig struct {
	GoroutineLimit                int64         `envconfig:"APP_DESTINATIONS_SENSITIVE_GOROUTINE_LIMIT,default=10"`
	RetryInterval                 time.Duration `envconfig:"APP_DESTINATIONS_RETRY_INTERVAL,default=100ms"`
	RetryAttempts                 uint          `envconfig:"APP_DESTINATIONS_RETRY_ATTEMPTS,default=3"`
	EndpointGetTenantDestinations string        `envconfig:"APP_ENDPOINT_GET_TENANT_DESTINATIONS,default=/destination-configuration/v1/subaccountDestinations"`
	EndpointFindDestination       string        `envconfig:"APP_ENDPOINT_FIND_DESTINATION,default=/destination-configuration/v1/destinations"`
	Timeout                       time.Duration `envconfig:"APP_DESTINATIONS_TIMEOUT,default=30s"`
	PageSize                      int           `envconfig:"APP_DESTINATIONS_PAGE_SIZE,default=100"`
	PagingPageParam               string        `envconfig:"APP_DESTINATIONS_PAGE_PARAM,default=$page"`
	PagingSizeParam               string        `envconfig:"APP_DESTINATIONS_PAGE_SIZE_PARAM,default=$pageSize"`
	PagingCountParam              string        `envconfig:"APP_DESTINATIONS_PAGE_COUNT_PARAM,default=$pageCount"`
	PagingCountHeader             string        `envconfig:"APP_DESTINATIONS_PAGE_COUNT_HEADER,default=Page-Count"`
}

// Client destination client
type Client struct {
	httpClient *http.Client
	apiConfig  DestinationServiceAPIConfig
	apiURL     string
}

// destinationFromService destination received from destination service
type destinationFromService struct {
	Name                    string `json:"Name"`
	Type                    string `json:"Type"`
	URL                     string `json:"URL"`
	Authentication          string `json:"Authentication"`
	XFSystemName            string `json:"XFSystemName"`
	CommunicationScenarioID string `json:"communicationScenarioId"`
	ProductName             string `json:"product.name"`
	XCorrelationID          string `json:"x-correlation-id"`
	XSystemTenantID         string `json:"x-system-id"`
	XSystemTenantName       string `json:"x-system-name"`
	XSystemType             string `json:"x-system-type"`
}

// ToModel missing godoc
func (d *destinationFromService) ToModel() (model.DestinationInput, error) {
	result := model.DestinationInput{
		Name:              d.Name,
		Type:              d.Type,
		URL:               d.URL,
		Authentication:    d.Authentication,
		XCorrelationID:    d.XCorrelationID,
		XSystemTenantID:   d.XSystemTenantID,
		XSystemTenantName: d.XSystemTenantName,
		XSystemType:       d.XSystemType,
	}

	// Set values from custom properties
	if result.XSystemType == "" {
		result.XSystemType = d.ProductName
	}
	if result.XSystemType == "" {
		return model.DestinationInput{}, errors.New("system type not found in destination")
	}
	if result.XCorrelationID == "" {
		if result.XSystemType == s4HANAType && d.CommunicationScenarioID != "" {
			result.XCorrelationID = correlationIDPrefix + d.CommunicationScenarioID
		}
	}
	if result.XCorrelationID == "" {
		return model.DestinationInput{}, errors.New("no correlation id found in destination")
	}
	if result.XSystemTenantName == "" {
		result.XSystemTenantName = d.XFSystemName
	}

	if result.XSystemTenantID == "" && (result.XSystemTenantName == "" || result.URL == "") {
		return model.DestinationInput{}, errors.New("system tenant could not be determined by destination")
	}
	return result, nil
}

// DestinationResponse paged response from destination service
type DestinationResponse struct {
	destinations []destinationFromService
	pageCount    string
}

// NewClient returns new destination client
func NewClient(instanceConfig config.InstanceConfig, apiConfig DestinationServiceAPIConfig,
	tokenPath, subdomain string) (*Client, error) {
	ctx := context.Background()

	baseTokenURL, err := url.Parse(instanceConfig.TokenURL)
	if err != nil {
		return nil, errors.Errorf("failed to parse auth url '%s': %v", instanceConfig.TokenURL, err)
	}
	parts := strings.Split(baseTokenURL.Hostname(), ".")
	if len(parts) < 2 {
		return nil, errors.Errorf("auth url '%s' should have a subdomain", instanceConfig.TokenURL)
	}
	originalSubdomain := parts[0]

	tokenURL := strings.Replace(instanceConfig.TokenURL, originalSubdomain, subdomain, 1) + tokenPath
	cfg := clientcredentials.Config{
		ClientID:  instanceConfig.ClientID,
		TokenURL:  tokenURL,
		AuthStyle: oauth2.AuthStyleInParams,
	}
	cert, err := tls.X509KeyPair([]byte(instanceConfig.Cert), []byte(instanceConfig.Key))
	if err != nil {
		return nil, errors.Errorf("failed to create destinations client x509 pair: %v", err)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	mtlsClient := &http.Client{
		Transport: transport,
		Timeout:   apiConfig.Timeout,
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, mtlsClient)

	httpClient := cfg.Client(ctx)
	httpClient.Timeout = apiConfig.Timeout

	return &Client{
		httpClient: httpClient,
		apiConfig:  apiConfig,
		apiURL:     instanceConfig.URL,
	}, nil
}

// FetchTenantDestinationsPage returns a page of destinations
func (c *Client) FetchTenantDestinationsPage(ctx context.Context, page string) (*DestinationResponse, error) {
	url := c.apiURL + c.apiConfig.EndpointGetTenantDestinations
	req, err := c.buildRequest(ctx, url, page)
	if err != nil {
		return nil, err
	}

	log.C(ctx).Infof("Getting destinations page: %s data from: %s\n", page, url)

	res, err := c.sendRequestWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("received status code %d when trying to fetch destinations", res.StatusCode)
	}

	var destinations []destinationFromService
	if err := json.NewDecoder(res.Body).Decode(&destinations); err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	pageCount := res.Header.Get(c.apiConfig.PagingCountHeader)
	if pageCount == "" {
		return nil, errors.Errorf("failed to extract header '%s' from destinations response", c.apiConfig.PagingCountParam)
	}

	return &DestinationResponse{
		destinations: destinations,
		pageCount:    pageCount,
	}, nil
}

func (c *Client) buildRequest(ctx context.Context, url string, page string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build request")
	}

	query := req.URL.Query()
	query.Add(c.apiConfig.PagingCountParam, "true")
	query.Add(c.apiConfig.PagingPageParam, page)
	query.Add(c.apiConfig.PagingSizeParam, strconv.Itoa(c.apiConfig.PageSize))
	req.URL.RawQuery = query.Encode()
	return req, nil
}

// FetchDestinationSensitiveData returns sensitive data of a destination
func (c *Client) FetchDestinationSensitiveData(ctx context.Context, destinationName string) ([]byte, error) {
	url := fmt.Sprintf("%s%s/%s", c.apiURL, c.apiConfig.EndpointFindDestination, destinationName)
	log.C(ctx).Infof("Getting destination data from: %s \n", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build request")
	}

	res, err := c.sendRequestWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, apperrors.NewNotFoundError(resource.Destination, destinationName)
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("received status code %d when trying to get destination info for %s",
			res.StatusCode, destinationName)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, errors.Wrap(err, "failed to read body of response")
	}

	return body, nil
}

func (c *Client) sendRequestWithRetry(req *http.Request) (*http.Response, error) {
	var response *http.Response
	err := retry.Do(func() error {
		res, err := c.httpClient.Do(req)
		if err != nil {
			return errors.Wrap(err, "failed to execute HTTP request")
		}

		if err == nil && res.StatusCode < http.StatusInternalServerError {
			response = res
			return nil
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}
		return errors.Errorf("request failed with status code %d, error message: %v", res.StatusCode, string(body))
	}, retry.Attempts(c.apiConfig.RetryAttempts), retry.Delay(c.apiConfig.RetryInterval))

	if err != nil {
		return nil, err
	}
	return response, nil
}
