/*
 * Copyright 2020 The Compass Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/pkg/fixtures"
	"github.com/kyma-incubator/compass/tests/pkg/ptr"
	tenant_utils "github.com/kyma-incubator/compass/tests/pkg/tenant"
	"github.com/kyma-incubator/compass/tests/pkg/testctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	tenantPathParamValue       = "tenant"
	regionPathParamValue       = "eu-1"
	defaultSubdomain           = "default-subdomain"
	defaultSubaccountSubdomain = "default-subaccount-subdomain"
)

type Tenant struct {
	TenantID               string
	SubaccountID           string
	CustomerID             string
	Subdomain              string
	SubscriptionProviderID string
}

func TestOnboardingHandler(t *testing.T) {
	t.Run("Success with tenant and customerID", func(t *testing.T) {
		tenantWithCustomer := Tenant{
			TenantID:               uuid.New().String(),
			CustomerID:             uuid.New().String(),
			Subdomain:              defaultSubdomain,
			SubscriptionProviderID: uuid.New().String(),
		}
		// WHEN
		addTenantExpectStatusCode(t, tenantWithCustomer, http.StatusOK)

		tenant, err := fixtures.GetTenantByExternalID(dexGraphQLClient, tenantWithCustomer.TenantID)
		require.NoError(t, err)

		parent, err := fixtures.GetTenantByExternalID(dexGraphQLClient, tenantWithCustomer.CustomerID)
		require.NoError(t, err)

		// THEN
		assertTenant(t, tenant, tenantWithCustomer.TenantID, tenantWithCustomer.Subdomain)
		assertTenant(t, parent, tenantWithCustomer.CustomerID, "")
	})

	t.Run("Success with only tenant", func(t *testing.T) {
		tenant := Tenant{
			TenantID:               uuid.New().String(),
			Subdomain:              defaultSubdomain,
			SubscriptionProviderID: uuid.New().String(),
		}

		addTenantExpectStatusCode(t, tenant, http.StatusOK)

		tnt, err := fixtures.GetTenantByExternalID(dexGraphQLClient, tenant.TenantID)
		require.NoError(t, err)

		// THEN
		assertTenant(t, tnt, tenant.TenantID, tenant.Subdomain)
	})

	t.Run("Successful account tenant creation with matching account and subaccount tenant IDs", func(t *testing.T) {
		id := uuid.New().String()
		tenant := Tenant{
			TenantID:               id,
			SubaccountID:           id,
			Subdomain:              defaultSubdomain,
			SubscriptionProviderID: uuid.New().String(),
		}

		addTenantExpectStatusCode(t, tenant, http.StatusOK)

		tnt, err := fixtures.GetTenantByExternalID(dexGraphQLClient, tenant.TenantID)
		require.NoError(t, err)

		// THEN
		assertTenant(t, tnt, tenant.TenantID, tenant.Subdomain)
	})

	t.Run("Successful account tenant creation with matching customer and account tenant IDs", func(t *testing.T) {
		id := uuid.New().String()
		tenant := Tenant{
			CustomerID:             id,
			TenantID:               id,
			Subdomain:              defaultSubdomain,
			SubscriptionProviderID: uuid.New().String(),
		}

		addTenantExpectStatusCode(t, tenant, http.StatusOK)

		tnt, err := fixtures.GetTenantByExternalID(dexGraphQLClient, tenant.TenantID)
		require.NoError(t, err)

		// THEN
		assertTenant(t, tnt, tenant.TenantID, tenant.Subdomain)
	})

	t.Run("Should not add already existing tenants", func(t *testing.T) {
		tenantWithCustomer := Tenant{
			TenantID:               uuid.New().String(),
			CustomerID:             uuid.New().String(),
			Subdomain:              defaultSubdomain,
			SubscriptionProviderID: uuid.New().String(),
		}
		//GIVEN
		oldTenantState, err := fixtures.GetTenants(dexGraphQLClient)
		require.NoError(t, err)

		//WHEN
		for i := 0; i < 10; i++ {
			addTenantExpectStatusCode(t, tenantWithCustomer, http.StatusOK)
		}

		tenants, err := fixtures.GetTenants(dexGraphQLClient)
		require.NoError(t, err)

		// THEN
		assert.Equal(t, len(oldTenantState)+2, len(tenants))
		assertTenantExists(t, tenants, tenantWithCustomer.TenantID)
		assertTenantExists(t, tenants, tenantWithCustomer.CustomerID)
	})

	t.Run("Should fail when no tenantID is provided", func(t *testing.T) {
		providedTenant := Tenant{
			CustomerID:             uuid.New().String(),
			SubscriptionProviderID: uuid.New().String(),
		}

		oldTenantState, err := fixtures.GetTenants(dexGraphQLClient)
		require.NoError(t, err)

		addTenantExpectStatusCode(t, providedTenant, http.StatusBadRequest)

		tenants, err := fixtures.GetTenants(dexGraphQLClient)
		require.NoError(t, err)

		// THEN
		assert.Equal(t, len(oldTenantState), len(tenants))
	})

	t.Run("Should fail when no subdomain is provided", func(t *testing.T) {
		providedTenant := Tenant{
			TenantID:               uuid.New().String(),
			CustomerID:             uuid.New().String(),
			SubscriptionProviderID: uuid.New().String(),
		}

		oldTenantState, err := fixtures.GetTenants(dexGraphQLClient)
		require.NoError(t, err)

		addTenantExpectStatusCode(t, providedTenant, http.StatusBadRequest)

		tenants, err := fixtures.GetTenants(dexGraphQLClient)
		require.NoError(t, err)

		// THEN
		assert.Equal(t, len(oldTenantState), len(tenants))
	})

	t.Run("Should fail when no SubscriptionProviderID is provided", func(t *testing.T) {
		providedTenant := Tenant{
			TenantID:   uuid.New().String(),
			CustomerID: uuid.New().String(),
			Subdomain:  defaultSubdomain,
		}

		oldTenantState, err := fixtures.GetTenants(dexGraphQLClient)
		require.NoError(t, err)

		addTenantExpectStatusCode(t, providedTenant, http.StatusBadRequest)

		tenants, err := fixtures.GetTenants(dexGraphQLClient)
		require.NoError(t, err)

		// THEN
		assert.Equal(t, len(oldTenantState), len(tenants))
	})

	t.Run("Should fail with subaccount tenant", func(t *testing.T) {
		// GIVEN
		parentTenant := Tenant{
			TenantID:               uuid.New().String(),
			Subdomain:              defaultSubdomain,
			SubscriptionProviderID: uuid.New().String(),
		}
		childTenant := Tenant{
			SubaccountID:           uuid.New().String(),
			TenantID:               parentTenant.TenantID,
			Subdomain:              defaultSubdomain,
			SubscriptionProviderID: uuid.New().String(),
		}

		addTenantExpectStatusCode(t, parentTenant, http.StatusOK)

		parent, err := fixtures.GetTenantByExternalID(dexGraphQLClient, parentTenant.TenantID)
		require.NoError(t, err)
		assertTenant(t, parent, parentTenant.TenantID, parentTenant.Subdomain)

		// THEN
		addTenantExpectStatusCode(t, childTenant, http.StatusInternalServerError)
	})
}

func TestDecommissioningHandler(t *testing.T) {
	t.Run("Success noop", func(t *testing.T) {
		providedTenant := Tenant{
			TenantID:               uuid.New().String(),
			Subdomain:              defaultSubdomain,
			SubscriptionProviderID: uuid.New().String(),
		}

		addTenantExpectStatusCode(t, providedTenant, http.StatusOK)

		oldTenantState, err := fixtures.GetTenants(dexGraphQLClient)
		require.NoError(t, err)

		removeTenantExpectStatusCode(t, providedTenant, http.StatusOK)

		newTenantState, err := fixtures.GetTenants(dexGraphQLClient)
		require.NoError(t, err)

		// THEN
		assert.Equal(t, len(oldTenantState), len(newTenantState))
	})
}

func TestRegionalOnboardingHandler(t *testing.T) {
	t.Run("Regional account tenant creation", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// GIVEN
			providedTenant := Tenant{
				TenantID:               uuid.New().String(),
				Subdomain:              defaultSubdomain,
				SubscriptionProviderID: uuid.New().String(),
			}

			// WHEN
			addRegionalTenantExpectStatusCode(t, providedTenant, http.StatusOK)

			// THEN
			tenant, err := fixtures.GetTenantByExternalID(dexGraphQLClient, providedTenant.TenantID)
			require.NoError(t, err)
			assertTenant(t, tenant, providedTenant.TenantID, providedTenant.Subdomain)
			require.Equal(t, regionPathParamValue, tenant.Labels["region"])
		})
	})

	t.Run("Regional subaccount tenant creation", func(t *testing.T) {
		t.Run("Success when parent account tenant is pre-existing", func(t *testing.T) {
			// GIVEN
			parentTenant := Tenant{
				TenantID:               uuid.New().String(),
				Subdomain:              defaultSubdomain,
				SubscriptionProviderID: uuid.New().String(),
			}
			childTenant := Tenant{
				SubaccountID:           uuid.New().String(),
				TenantID:               parentTenant.TenantID,
				Subdomain:              defaultSubaccountSubdomain,
				SubscriptionProviderID: uuid.New().String(),
			}

			addRegionalTenantExpectStatusCode(t, parentTenant, http.StatusOK)

			parent, err := fixtures.GetTenantByExternalID(dexGraphQLClient, parentTenant.TenantID)
			require.NoError(t, err)
			assertTenant(t, parent, parentTenant.TenantID, parentTenant.Subdomain)
			require.Equal(t, regionPathParamValue, parent.Labels["region"])

			// WHEN
			addRegionalTenantExpectStatusCode(t, childTenant, http.StatusOK)

			// THEN
			tenant, err := fixtures.GetTenantByExternalID(dexGraphQLClient, childTenant.SubaccountID)
			require.NoError(t, err)
			assertTenant(t, tenant, childTenant.SubaccountID, childTenant.Subdomain)
			require.Equal(t, regionPathParamValue, tenant.Labels["region"])

			parentTenantAfterInsert, err := fixtures.GetTenantByExternalID(dexGraphQLClient, parentTenant.TenantID)
			require.NoError(t, err)
			assertTenant(t, parentTenantAfterInsert, parentTenant.TenantID, parentTenant.Subdomain)
			require.Equal(t, regionPathParamValue, parentTenantAfterInsert.Labels["region"])
		})

		t.Run("Success when parent account tenant does not exist", func(t *testing.T) {
			// GIVEN
			providedTenant := Tenant{
				TenantID:               uuid.New().String(),
				CustomerID:             uuid.New().String(),
				SubaccountID:           uuid.New().String(),
				Subdomain:              defaultSubaccountSubdomain,
				SubscriptionProviderID: uuid.New().String(),
			}

			// THEN
			addRegionalTenantExpectStatusCode(t, providedTenant, http.StatusOK)

			// THEN
			childTenant, err := fixtures.GetTenantByExternalID(dexGraphQLClient, providedTenant.SubaccountID)
			require.NoError(t, err)
			assertTenant(t, childTenant, providedTenant.SubaccountID, providedTenant.Subdomain)
			require.Equal(t, regionPathParamValue, childTenant.Labels["region"])

			parentTenant, err := fixtures.GetTenantByExternalID(dexGraphQLClient, providedTenant.TenantID)
			require.NoError(t, err)
			assertTenant(t, parentTenant, providedTenant.TenantID, "")
			require.Empty(t, parentTenant.Labels)

			customerTenant, err := fixtures.GetTenantByExternalID(dexGraphQLClient, providedTenant.CustomerID)
			require.NoError(t, err)
			assertTenant(t, customerTenant, providedTenant.CustomerID, "")
			require.Empty(t, customerTenant.Labels)
		})

		t.Run("Should not fail when tenant already exists", func(t *testing.T) {
			// GIVEN
			parentTenantId := uuid.New().String()
			parentTenant := Tenant{
				TenantID:               parentTenantId,
				Subdomain:              defaultSubaccountSubdomain,
				SubscriptionProviderID: uuid.New().String(),
			}
			childTenant := Tenant{
				TenantID:               parentTenantId,
				SubaccountID:           uuid.New().String(),
				Subdomain:              defaultSubaccountSubdomain,
				SubscriptionProviderID: uuid.New().String(),
			}
			oldTenantState, err := fixtures.GetTenants(dexGraphQLClient)
			require.NoError(t, err)

			addTenantExpectStatusCode(t, parentTenant, http.StatusOK)
			parent, err := fixtures.GetTenantByExternalID(dexGraphQLClient, parentTenant.TenantID)
			require.NoError(t, err)
			assertTenant(t, parent, parentTenant.TenantID, parentTenant.Subdomain)

			// WHEN
			for i := 0; i < 10; i++ {
				addRegionalTenantExpectStatusCode(t, childTenant, http.StatusOK)
			}

			tenant, err := fixtures.GetTenantByExternalID(dexGraphQLClient, childTenant.SubaccountID)
			require.NoError(t, err)

			tenants, err := fixtures.GetTenants(dexGraphQLClient)
			require.NoError(t, err)

			// THEN
			assertTenant(t, tenant, childTenant.SubaccountID, childTenant.Subdomain)
			assert.Equal(t, len(oldTenantState)+2, len(tenants))
		})

		t.Run("Should fail when parent tenantID is not provided", func(t *testing.T) {
			// GIVEN
			providedTenant := Tenant{
				CustomerID:             uuid.New().String(),
				SubaccountID:           uuid.New().String(),
				Subdomain:              defaultSubaccountSubdomain,
				SubscriptionProviderID: uuid.New().String(),
			}
			oldTenantState, err := fixtures.GetTenants(dexGraphQLClient)
			require.NoError(t, err)

			// WHEN
			addRegionalTenantExpectStatusCode(t, providedTenant, http.StatusBadRequest)

			// THEN
			tenants, err := fixtures.GetTenants(dexGraphQLClient)
			require.NoError(t, err)
			assert.Equal(t, len(oldTenantState), len(tenants))
		})

		t.Run("Should fail when subdomain is not provided", func(t *testing.T) {
			// GIVEN
			providedTenant := Tenant{
				TenantID:               uuid.New().String(),
				SubaccountID:           uuid.New().String(),
				CustomerID:             uuid.New().String(),
				SubscriptionProviderID: uuid.New().String(),
			}
			oldTenantState, err := fixtures.GetTenants(dexGraphQLClient)
			require.NoError(t, err)

			// WHEN
			addRegionalTenantExpectStatusCode(t, providedTenant, http.StatusBadRequest)

			// THEN
			tenants, err := fixtures.GetTenants(dexGraphQLClient)
			require.NoError(t, err)
			assert.Equal(t, len(oldTenantState), len(tenants))
		})

		t.Run("Should fail when SubscriptionProviderID is not provided", func(t *testing.T) {
			// GIVEN
			providedTenant := Tenant{
				TenantID:     uuid.New().String(),
				SubaccountID: uuid.New().String(),
				CustomerID:   uuid.New().String(),
			}
			oldTenantState, err := fixtures.GetTenants(dexGraphQLClient)
			require.NoError(t, err)

			// WHEN
			addRegionalTenantExpectStatusCode(t, providedTenant, http.StatusBadRequest)

			// THEN
			tenants, err := fixtures.GetTenants(dexGraphQLClient)
			require.NoError(t, err)
			assert.Equal(t, len(oldTenantState), len(tenants))
		})
	})

	t.Run("Regional subaccount tenant subscription flows", func(t *testing.T) {
		t.Run("Subscribe tenant for correct consumer", func(t *testing.T) {
			// GIVEN
			ctx := context.Background()
			subscriptionConsumerID := uuid.New().String()
			subaccountID := uuid.New().String()
			accountID := uuid.New().String()

			tenantId := tenant_utils.TestTenants.GetDefaultTenantID()

			subscribedRuntime := registerRuntime(t, ctx, "subscribed-runtime", subscriptionConsumerID)
			defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, tenantId, &subscribedRuntime)

			notSubscribedRuntime := registerRuntime(t, ctx, "not-subscribed-runtime", "fake_subscrioption_id")
			defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, tenantId, &notSubscribedRuntime)

			// WHEN
			tenant := Tenant{
				TenantID:               accountID,
				SubaccountID:           subaccountID,
				Subdomain:              defaultSubdomain,
				SubscriptionProviderID: subscriptionConsumerID,
			}

			addRegionalTenantExpectStatusCode(t, tenant, http.StatusOK)

			// THEN
			assertRuntimeSubscription(t, ctx, subscribedRuntime.ID, tenant)

			assertNoRuntimeSubscription(t, ctx, notSubscribedRuntime.ID)
		})

		t.Run("Unsubscribe tenant from correct consumer", func(t *testing.T) {
			// GIVEN
			ctx := context.Background()
			subscriptionConsumerID := uuid.New().String()
			secondSubscriptionConsumerID := uuid.New().String()
			subaccountID := uuid.New().String()
			accountID := uuid.New().String()

			tenantId := tenant_utils.TestTenants.GetDefaultTenantID()

			subscribedRuntime := registerRuntime(t, ctx, "subscribed-runtime", subscriptionConsumerID)
			defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, tenantId, &subscribedRuntime)

			secondSubscribedRuntime := registerRuntime(t, ctx, "second-subscribed-runtime", secondSubscriptionConsumerID)
			defer fixtures.CleanupRuntime(t, ctx, dexGraphQLClient, tenantId, &secondSubscribedRuntime)

			tenant := Tenant{
				TenantID:               accountID,
				SubaccountID:           subaccountID,
				Subdomain:              defaultSubdomain,
				SubscriptionProviderID: subscriptionConsumerID,
			}
			addRegionalTenantExpectStatusCode(t, tenant, http.StatusOK)
			assertRuntimeSubscription(t, ctx, subscribedRuntime.ID, tenant)

			tenantSecondSubscription := Tenant{
				TenantID:               accountID,
				SubaccountID:           subaccountID,
				Subdomain:              defaultSubdomain,
				SubscriptionProviderID: secondSubscriptionConsumerID,
			}
			addRegionalTenantExpectStatusCode(t, tenantSecondSubscription, http.StatusOK)
			assertRuntimeSubscription(t, ctx, secondSubscribedRuntime.ID, tenantSecondSubscription)

			// WHEN
			removeRegionalTenantExpectStatusCode(t, tenant, http.StatusOK)

			// THEN
			assertNoRuntimeSubscription(t, ctx, subscribedRuntime.ID)
			assertRuntimeSubscription(t, ctx, secondSubscribedRuntime.ID, tenantSecondSubscription)
		})
	})
}

func TestGetDependenciesHandler(t *testing.T) {
	t.Run("Returns empty body", func(t *testing.T) {
		// GIVEN
		request, err := http.NewRequest(http.MethodGet, config.TenantFetcherFullDependenciesURL, nil)
		require.NoError(t, err)
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", fetchToken(t)))

		// WHEN
		response, err := httpClient.Do(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		responseBody, err := ioutil.ReadAll(response.Body)
		require.NoError(t, err)
		responseBodyJson := make(map[string]interface{}, 0)

		// THEN
		err = json.Unmarshal(responseBody, &responseBodyJson)
		require.NoError(t, err)
		require.Empty(t, responseBodyJson)
	})
}

func addTenantExpectStatusCode(t *testing.T, providedTenant Tenant, expectedStatusCode int) {
	makeTenantRequestExpectStatusCode(t, providedTenant, http.MethodPut, config.TenantFetcherFullURL, expectedStatusCode)
}

func addRegionalTenantExpectStatusCode(t *testing.T, providedTenant Tenant, expectedStatusCode int) {
	makeTenantRequestExpectStatusCode(t, providedTenant, http.MethodPut, config.TenantFetcherFullRegionalURL, expectedStatusCode)
}

func removeTenantExpectStatusCode(t *testing.T, providedTenant Tenant, expectedStatusCode int) {
	makeTenantRequestExpectStatusCode(t, providedTenant, http.MethodDelete, config.TenantFetcherFullURL, expectedStatusCode)
}

func removeRegionalTenantExpectStatusCode(t *testing.T, providedTenant Tenant, expectedStatusCode int) {
	makeTenantRequestExpectStatusCode(t, providedTenant, http.MethodDelete, config.TenantFetcherFullRegionalURL, expectedStatusCode)
}

func makeTenantRequestExpectStatusCode(t *testing.T, providedTenant Tenant, httpMethod, url string, expectedStatusCode int) {
	request := createTenantRequest(t, providedTenant, httpMethod, url)

	t.Log(fmt.Sprintf("Provisioning tenant with ID %s", actualTenantID(providedTenant)))
	response, err := httpClient.Do(request)
	require.NoError(t, err)
	require.Equal(t, expectedStatusCode, response.StatusCode)
}

func actualTenantID(tenant Tenant) string {
	if len(tenant.SubaccountID) > 0 {
		return tenant.SubaccountID
	}

	return tenant.TenantID
}

func createTenantRequest(t *testing.T, tenant Tenant, httpMethod string, url string) *http.Request {
	var (
		body = "{}"
		err  error
	)

	if len(tenant.TenantID) > 0 {
		body, err = sjson.Set(body, config.TenantIDProperty, tenant.TenantID)
		require.NoError(t, err)
	}
	if len(tenant.SubaccountID) > 0 {
		body, err = sjson.Set(body, config.SubaccountTenantIDProperty, tenant.SubaccountID)
		require.NoError(t, err)
	}
	if len(tenant.CustomerID) > 0 {
		body, err = sjson.Set(body, config.CustomerIDProperty, tenant.CustomerID)
		require.NoError(t, err)
	}
	if len(tenant.Subdomain) > 0 {
		body, err = sjson.Set(body, config.SubdomainProperty, tenant.Subdomain)
		require.NoError(t, err)
	}
	if len(tenant.SubscriptionProviderID) > 0 {
		body, err = sjson.Set(body, config.SubscriptionProviderIDProperty, tenant.SubscriptionProviderID)
		require.NoError(t, err)
	}

	request, err := http.NewRequest(httpMethod, url, bytes.NewBuffer([]byte(body)))
	require.NoError(t, err)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", fetchToken(t)))

	return request
}

func fetchToken(t *testing.T) string {
	claims := map[string]interface{}{
		"test": "tenant-fetcher",
		"scope": []string{
			"prefix.Callback",
		},
		"tenant":   "tenant",
		"identity": "tenant-fetcher-tests",
		"iss":      config.ExternalServicesMockURL,
		"exp":      time.Now().Unix() + int64(time.Minute.Seconds()),
	}

	data, err := json.Marshal(claims)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, config.ExternalServicesMockURL+"/oauth/token", bytes.NewBuffer(data))
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	token := gjson.GetBytes(body, "access_token")
	require.True(t, token.Exists())

	return token.String()
}

func assertTenant(t *testing.T, tenant *graphql.Tenant, tenantID, subdomain string) {
	require.Equal(t, tenantID, tenant.ID)
	if len(subdomain) > 0 {
		require.Equal(t, subdomain, tenant.Labels["subdomain"])
	}
}

func assertTenantExists(t *testing.T, tenants []*graphql.Tenant, tenantID string) {
	for _, tenant := range tenants {
		if tenant.ID == tenantID {
			return
		}
	}

	require.Fail(t, fmt.Sprintf("Tenant with ID %q not found in %v", tenantID, tenants))
}

func registerRuntime(t *testing.T, ctx context.Context, name, subscriptionConsumerID string) graphql.RuntimeExt {
	runtimeInput := graphql.RuntimeInput{
		Name:        name,
		Description: ptr.String(name),
		Labels:      graphql.Labels{"region": regionPathParamValue, config.SubscriptionProviderLabelKey: subscriptionConsumerID},
	}
	runtimeInGQL, err := testctx.Tc.Graphqlizer.RuntimeInputToGQL(runtimeInput)
	require.NoError(t, err)
	actualRuntime := graphql.RuntimeExt{}

	registerReq := fixtures.FixRegisterRuntimeRequest(runtimeInGQL)
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, registerReq, &actualRuntime)
	require.NoError(t, err)
	return actualRuntime
}

func assertRuntimeSubscription(t *testing.T, ctx context.Context, runtimeID string, tenant Tenant) {
	actualTenant, err := fixtures.GetTenantByExternalID(dexGraphQLClient, tenant.SubaccountID)
	require.NoError(t, err)
	assertTenant(t, actualTenant, tenant.SubaccountID, tenant.Subdomain)
	require.Equal(t, regionPathParamValue, actualTenant.Labels["region"])

	subscribedRuntime := graphql.RuntimeExt{}
	getSubscribedReq := fixtures.FixGetRuntimeRequest(runtimeID)
	err = testctx.Tc.RunOperation(ctx, dexGraphQLClient, getSubscribedReq, &subscribedRuntime)
	require.NoError(t, err)
	consumerSubaccountIDsLabel, ok := subscribedRuntime.Labels[config.ConsumerSubaccountIDsLabelKey].([]interface{})
	require.Equal(t, true, ok)
	assert.Len(t, consumerSubaccountIDsLabel, 1)
	labelValue, ok := consumerSubaccountIDsLabel[0].(string)
	require.Equal(t, true, ok)
	require.Equal(t, tenant.SubaccountID, labelValue)
}

func assertNoRuntimeSubscription(t *testing.T, ctx context.Context, runtimeID string) {
	notSubscribedRuntime := graphql.RuntimeExt{}
	getNotSubscribedReq := fixtures.FixGetRuntimeRequest(runtimeID)
	err := testctx.Tc.RunOperation(ctx, dexGraphQLClient, getNotSubscribedReq, &notSubscribedRuntime)
	require.NoError(t, err)

	if _, ok := notSubscribedRuntime.Labels[config.ConsumerSubaccountIDsLabelKey]; !ok {
		return
	}
	require.Len(t, notSubscribedRuntime.Labels[config.ConsumerSubaccountIDsLabelKey], 0)
}
