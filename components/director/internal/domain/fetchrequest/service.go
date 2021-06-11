package fetchrequest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/internal/timestamp"

	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/kyma-incubator/compass/components/director/internal/model"
)

type service struct {
	repo         FetchRequestRepository
	client       *http.Client
	timestampGen timestamp.Generator
}

//go:generate mockery --name=FetchRequestRepository --output=automock --outpkg=automock --case=underscore
type FetchRequestRepository interface {
	Update(ctx context.Context, item *model.FetchRequest) error
}

func NewService(repo FetchRequestRepository, client *http.Client) *service {
	return &service{
		repo:         repo,
		client:       client,
		timestampGen: timestamp.DefaultGenerator(),
	}
}

func (s *service) HandleSpec(ctx context.Context, fr *model.FetchRequest) *string {
	var data *string
	data, fr.Status = s.fetchSpec(ctx, fr)

	err := s.repo.Update(ctx, fr)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error has occurred while updating fetch request status: %v", err)
		return nil
	}

	return data
}

func (s *service) fetchSpec(ctx context.Context, fr *model.FetchRequest) (*string, *model.FetchRequestStatus) {
	err := s.validateFetchRequest(fr)
	if err != nil {
		log.C(ctx).WithError(err).Error()
		return nil, FixStatus(model.FetchRequestStatusConditionInitial, str.Ptr(err.Error()), s.timestampGen())
	}

	var resp *http.Response
	if fr.Auth != nil {
		resp, err = s.requestWithCredentials(ctx, fr)
	} else {
		resp, err = s.requestWithoutCredentials(fr)
	}

	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error has occurred while fetching Spec: %v", err)
		return nil, FixStatus(model.FetchRequestStatusConditionFailed, str.Ptr(fmt.Sprintf("While fetching Spec: %s", err.Error())), s.timestampGen())
	}

	defer func() {
		if resp.Body != nil {
			err := resp.Body.Close()
			if err != nil {
				log.C(ctx).WithError(err).Errorf("An error has occurred while closing response body: %v", err)
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.C(ctx).WithError(err).Errorf("Failed to execute fetch request for %s with id %q: status code: %d: %v", fr.ObjectType, fr.ObjectID, resp.StatusCode, err)
		return nil, FixStatus(model.FetchRequestStatusConditionFailed, str.Ptr(fmt.Sprintf("While fetching Spec status code: %d", resp.StatusCode)), s.timestampGen())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("An error has occurred while reading Spec: %v", err)
		return nil, FixStatus(model.FetchRequestStatusConditionFailed, str.Ptr(fmt.Sprintf("While reading Spec: %s", err.Error())), s.timestampGen())
	}

	spec := string(body)
	return &spec, FixStatus(model.FetchRequestStatusConditionSucceeded, nil, s.timestampGen())
}

func (s *service) validateFetchRequest(fr *model.FetchRequest) error {
	if fr.Mode != model.FetchModeSingle {
		return apperrors.NewInvalidDataError("Unsupported fetch mode: %s", fr.Mode)
	}

	if fr.Filter != nil {
		return apperrors.NewInvalidDataError("Filter for Fetch Request was provided, currently it's unsupported")
	}

	return nil
}

func (s *service) requestWithCredentials(ctx context.Context, fr *model.FetchRequest) (*http.Response, error) {
	if fr.Auth.Credential.Basic == nil && fr.Auth.Credential.Oauth == nil {
		return nil, apperrors.NewInvalidDataError("Credentials not provided")
	}

	req, err := http.NewRequest(http.MethodGet, fr.URL, nil)
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	if fr.Auth.Credential.Basic != nil {
		req.SetBasicAuth(fr.Auth.Credential.Basic.Username, fr.Auth.Credential.Basic.Password)

		if fr.ProxyURL != "" {
			if err := s.setProxy(ctx, fr.ProxyURL); err != nil {
				return nil, err
			}
			defer func() {
				if err := s.removeProxy(); err != nil {
					log.C(ctx).Errorf("Error occurred while reverting proxy transport configuration: %v", err.Error())
				}
			}()
		}

		resp, err = s.client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}
	}

	if fr.Auth.Credential.Oauth != nil {
		resp, err = s.secureClient(ctx, fr).Do(req)
	}

	return resp, err
}

func (s *service) secureClient(ctx context.Context, fr *model.FetchRequest) *http.Client {
	conf := &clientcredentials.Config{
		ClientID:     fr.Auth.Credential.Oauth.ClientID,
		ClientSecret: fr.Auth.Credential.Oauth.ClientSecret,
		TokenURL:     fr.Auth.Credential.Oauth.URL,
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, s.client)
	securedClient := conf.Client(ctx)
	securedClient.Timeout = s.client.Timeout
	return securedClient
}

func (s *service) requestWithoutCredentials(fr *model.FetchRequest) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, fr.URL, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req)
}

func FixStatus(condition model.FetchRequestStatusCondition, message *string, timestamp time.Time) *model.FetchRequestStatus {
	return &model.FetchRequestStatus{
		Condition: condition,
		Message:   message,
		Timestamp: timestamp,
	}
}

func (s *service) setProxy(ctx context.Context, proxyURL string) error {
	proxyUrl, err := url.Parse(proxyURL)
	if err != nil {
		log.C(ctx).WithError(err).Warnf("Got error parsing proxy url: %s", proxyUrl)
		return err
	}

	transport := s.client.Transport.(*http.Transport)
	transport.Proxy = http.ProxyURL(proxyUrl)

	return nil
}

func (s *service) removeProxy() error {
	transport := s.client.Transport.(*http.Transport)
	transport.Proxy = nil

	return nil
}
