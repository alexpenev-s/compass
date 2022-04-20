package tenantfetchersvc

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/oauth"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	"github.com/kyma-incubator/compass/components/director/internal/tenantfetcher"

	"github.com/gorilla/mux"
	"github.com/kyma-incubator/compass/components/director/internal/features"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/tidwall/gjson"
)

const (
	// InternalServerError message
	InternalServerError = "Internal Server Error"
	compassURL          = "https://github.com/kyma-incubator/compass"
)

// TenantFetcher is used to fectch tenants for creation;
//go:generate mockery --name=TenantFetcher --output=automock --outpkg=automock --case=underscore
type TenantFetcher interface {
	FetchTenantOnDemand(ctx context.Context, tenantID string) error
}

// TenantSubscriber is used to apply subscription changes for tenants;
//go:generate mockery --name=TenantSubscriber --output=automock --outpkg=automock --case=underscore
type TenantSubscriber interface {
	Subscribe(ctx context.Context, tenantSubscriptionRequest *TenantSubscriptionRequest) error
	Unsubscribe(ctx context.Context, tenantSubscriptionRequest *TenantSubscriptionRequest) error
}

// HandlerConfig is the configuration required by the tenant handler.
// It includes configurable parameters for incoming requests, including different tenant IDs json properties, and path parameters.
type HandlerConfig struct {
	TenantOnDemandHandlerEndpoint string `envconfig:"APP_TENANT_ON_DEMAND_HANDLER_ENDPOINT,default=/v1/fetch/{tenantId}"`
	RegionalHandlerEndpoint       string `envconfig:"APP_REGIONAL_HANDLER_ENDPOINT,default=/v1/regional/{region}/callback/{tenantId}"`
	DependenciesEndpoint          string `envconfig:"APP_DEPENDENCIES_ENDPOINT,default=/v1/dependencies"`
	TenantPathParam               string `envconfig:"APP_TENANT_PATH_PARAM,default=tenantId"`
	RegionPathParam               string `envconfig:"APP_REGION_PATH_PARAM,default=region"`

	Features features.Config

	TenantFetcherJobIntervalMins time.Duration `envconfig:"default=5m,APP_TENANT_FETCH_JOB_INTERVAL_MINS"`
	FullResyncInterval           time.Duration `envconfig:"default=12h"`
	ShouldSyncSubaccounts        bool          `envconfig:"default=false,APP_SYNC_SUBACCOUNTS"`

	Kubernetes tenantfetcher.KubeConfig
	Database   persistence.DatabaseConfig

	DirectorGraphQLEndpoint     string        `envconfig:"APP_DIRECTOR_GRAPHQL_ENDPOINT"`
	ClientTimeout               time.Duration `envconfig:"default=60s"`
	HTTPClientSkipSslValidation bool          `envconfig:"APP_HTTP_CLIENT_SKIP_SSL_VALIDATION,default=false"`

	TenantInsertChunkSize int `envconfig:"default=500,APP_TENANT_INSERT_CHUNK_SIZE"`
	TenantProviderConfig
}

type TenantFetchConfig struct {
}

// TenantProviderConfig includes the configuration for tenant providers - the tenant ID json property names, the subdomain property name, and the tenant provider name.
type TenantProviderConfig struct {
	TenantIDProperty               string `envconfig:"APP_TENANT_PROVIDER_TENANT_ID_PROPERTY,default=tenantId"`
	SubaccountTenantIDProperty     string `envconfig:"APP_TENANT_PROVIDER_SUBACCOUNT_TENANT_ID_PROPERTY,default=subaccountTenantId"`
	CustomerIDProperty             string `envconfig:"APP_TENANT_PROVIDER_CUSTOMER_ID_PROPERTY,default=customerId"`
	SubdomainProperty              string `envconfig:"APP_TENANT_PROVIDER_SUBDOMAIN_PROPERTY,default=subdomain"`
	TenantProvider                 string `envconfig:"APP_TENANT_PROVIDER,default=external-provider"`
	SubscriptionProviderIDProperty string `envconfig:"APP_TENANT_PROVIDER_SUBSCRIPTION_PROVIDER_ID_PROPERTY,default=subscriptionProviderId"`
}

// EventsConfig contains configuration for Events API requests
type EventsConfig struct {
	AccountsRegion    string   `envconfig:"default=central,APP_ACCOUNT_REGION"`
	SubaccountRegions []string `envconfig:"default=central,APP_SUBACCOUNT_REGIONS"`

	AuthMode    oauth.AuthMode `envconfig:"APP_OAUTH_AUTH_MODE,default=standard"`
	OAuthConfig tenantfetcher.OAuth2Config
	APIConfig   tenantfetcher.APIConfig
	QueryConfig tenantfetcher.QueryConfig

	TenantFieldMapping          tenantfetcher.TenantFieldMapping
	MovedSubaccountFieldMapping tenantfetcher.MovedSubaccountsFieldMapping

	MetricsPushEndpoint string `envconfig:"optional,APP_METRICS_PUSH_ENDPOINT"`
}

type handler struct {
	fetcher    TenantFetcher
	subscriber TenantSubscriber
	config     HandlerConfig
}

// NewTenantsHTTPHandler returns a new HTTP handler, responsible for creation and deletion of regional and non-regional tenants.
func NewTenantsHTTPHandler(subscriber TenantSubscriber, config HandlerConfig) *handler {
	return &handler{
		subscriber: subscriber,
		config:     config,
	}
}

// NewTenantFetcherHTTPHandler returns a new HTTP handler, responsible for creation of on-demand tenants.
func NewTenantFetcherHTTPHandler(fetcher TenantFetcher, config HandlerConfig) *handler {
	return &handler{
		fetcher: fetcher,
		config:  config,
	}
}

// FetchTenantOnDemand fetches External tenants registry events for a provided subaccount and creates a subaccount tenant
func (h *handler) FetchTenantOnDemand(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	vars := mux.Vars(request)
	tenantID, ok := vars[h.config.TenantPathParam]
	if !ok {
		log.C(ctx).WithError(errors.New("tenant path parameter is missing from request")).Error()
		http.Error(writer, "Tenant path parameter is missing from request", http.StatusBadRequest)
		return
	}

	log.C(ctx).Infof("Fetching create event for tenant with ID %s", tenantID)

	err := h.fetcher.FetchTenantOnDemand(ctx, tenantID)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Error while processing request for creation of tenant %s: %v", tenantID, err)
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}
	writeCreatedResponse(writer, ctx, tenantID)
}

func writeCreatedResponse(writer http.ResponseWriter, ctx context.Context, tenantID string) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte(compassURL)); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to write response body for request for creation of tenant %s: %v", tenantID, err)
	}
}

// SubscribeTenant handles subscription for tenant. If tenant does not exist, will create it first.
func (h *handler) SubscribeTenant(writer http.ResponseWriter, request *http.Request) {
	h.applySubscriptionChange(writer, request, h.subscriber.Subscribe)
}

// UnSubscribeTenant handles unsubscription for tenant which will remove the tenant id label from the runtime
func (h *handler) UnSubscribeTenant(writer http.ResponseWriter, request *http.Request) {
	h.applySubscriptionChange(writer, request, h.subscriber.Unsubscribe)
}

// Dependencies handler returns all external services where once created in Compass, the tenant should be created as well.
func (h *handler) Dependencies(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write([]byte("{}")); err != nil {
		log.C(request.Context()).WithError(err).Errorf("Failed to write response body for dependencies request")
		return
	}
}

func (h *handler) applySubscriptionChange(writer http.ResponseWriter, request *http.Request, subscriptionFunc subscriptionFunc) {
	ctx := request.Context()

	vars := mux.Vars(request)
	region, ok := vars[h.config.RegionPathParam]
	if !ok {
		log.C(ctx).Error("Region path parameter is missing from request")
		http.Error(writer, "Region path parameter is missing from request", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to read tenant information from request body: %v", err)
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	subscriptionRequest, err := h.getSubscriptionRequest(body, region)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to extract tenant information from request body: %v", err)
		http.Error(writer, fmt.Sprintf("Failed to extract tenant information from request body: %v", err), http.StatusBadRequest)
		return
	}

	mainTenantID := subscriptionRequest.MainTenantID()
	if err := subscriptionFunc(ctx, subscriptionRequest); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to apply subscription change for tenant %s: %v", mainTenantID, err)
		http.Error(writer, InternalServerError, http.StatusInternalServerError)
		return
	}

	respondSuccess(ctx, writer, mainTenantID)
}

func (h *handler) getSubscriptionRequest(body []byte, region string) (*TenantSubscriptionRequest, error) {
	properties, err := getProperties(body, map[string]bool{
		h.config.TenantIDProperty:               true,
		h.config.SubaccountTenantIDProperty:     false,
		h.config.SubdomainProperty:              true,
		h.config.CustomerIDProperty:             false,
		h.config.SubscriptionProviderIDProperty: true,
	})
	if err != nil {
		return nil, err
	}

	req := &TenantSubscriptionRequest{
		AccountTenantID:        properties[h.config.TenantIDProperty],
		SubaccountTenantID:     properties[h.config.SubaccountTenantIDProperty],
		CustomerTenantID:       properties[h.config.CustomerIDProperty],
		Subdomain:              properties[h.config.SubdomainProperty],
		SubscriptionProviderID: properties[h.config.SubscriptionProviderIDProperty],
		Region:                 region,
	}

	if req.AccountTenantID == req.SubaccountTenantID {
		req.SubaccountTenantID = ""
	}

	if req.AccountTenantID == req.CustomerTenantID {
		req.CustomerTenantID = ""
	}

	return req, nil
}

func getProperties(body []byte, props map[string]bool) (map[string]string, error) {
	resultProps := map[string]string{}
	for propName, mandatory := range props {
		result := gjson.GetBytes(body, propName).String()
		if mandatory && len(result) == 0 {
			return nil, fmt.Errorf("mandatory property %q is missing from request body", propName)
		}
		resultProps[propName] = result
	}

	return resultProps, nil
}

func respondSuccess(ctx context.Context, writer http.ResponseWriter, mainTenantID string) {
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte(compassURL)); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to write response body for tenant request creation for tenant %s: %v", mainTenantID, err)
	}
}
