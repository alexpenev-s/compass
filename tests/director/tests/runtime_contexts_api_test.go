package tests

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/kyma-incubator/compass/tests/pkg/tenantfetcher"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/tests/pkg/assertions"
	"github.com/kyma-incubator/compass/tests/pkg/certs/certprovider"
	"github.com/kyma-incubator/compass/tests/pkg/fixtures"
	"github.com/kyma-incubator/compass/tests/pkg/gql"
	"github.com/kyma-incubator/compass/tests/pkg/ptr"
	"github.com/kyma-incubator/compass/tests/pkg/subscription"
	"github.com/kyma-incubator/compass/tests/pkg/tenant"
	"github.com/kyma-incubator/compass/tests/pkg/testctx"
	testingx "github.com/kyma-incubator/compass/tests/pkg/testing"
	"github.com/kyma-incubator/compass/tests/pkg/token"
	"github.com/stretchr/testify/require"
)

func TestAddRuntimeContext(t *testing.T) {
	ctx := context.Background()

	tenantId := tenant.TestTenants.GetDefaultTenantID()

	in := fixRuntimeInput("addRuntimeContext")

	runtime, err := fixtures.RegisterRuntimeFromInputWithinTenant(t, ctx, certSecuredGraphQLClient, tenantId, &in)
	defer fixtures.CleanupRuntime(t, ctx, certSecuredGraphQLClient, tenantId, &runtime)
	require.NoError(t, err)
	require.NotEmpty(t, runtime.ID)

	rtmCtxInput := fixtures.FixRuntimeContextInput("create", "create")
	rtmCtxInputGQL, err := testctx.Tc.Graphqlizer.RuntimeContextInputToGQL(rtmCtxInput)
	require.NoError(t, err)

	addRtmCtxRequest := fixtures.FixAddRuntimeContextRequest(runtime.ID, rtmCtxInputGQL)
	output := graphql.RuntimeContextExt{}

	// WHEN
	t.Log("Create runtimeContext")
	err = testctx.Tc.RunOperation(ctx, certSecuredGraphQLClient, addRtmCtxRequest, &output)
	defer fixtures.DeleteRuntimeContext(t, ctx, certSecuredGraphQLClient, tenantId, output.ID)

	// THEN
	require.NoError(t, err)
	require.NotEmpty(t, output.ID)
	assertions.AssertRuntimeContext(t, &rtmCtxInput, &output)

	saveExample(t, addRtmCtxRequest.Query(), "register runtime context")

	rtmCtxRequest := fixtures.FixRuntimeContextRequest(runtime.ID, output.ID)
	runtimeFromAPI := graphql.RuntimeExt{}

	err = testctx.Tc.RunOperation(ctx, certSecuredGraphQLClient, rtmCtxRequest, &runtimeFromAPI)
	require.NoError(t, err)

	assertions.AssertRuntimeContext(t, &rtmCtxInput, &runtimeFromAPI.RuntimeContext)
	saveExample(t, rtmCtxRequest.Query(), "query runtimeContext")
}

func TestQueryRuntimeContexts(t *testing.T) {
	ctx := context.Background()

	tenantId := tenant.TestTenants.GetDefaultTenantID()

	in := fixRuntimeInput("addRuntimeContext")

	runtime, err := fixtures.RegisterRuntimeFromInputWithinTenant(t, ctx, certSecuredGraphQLClient, tenantId, &in)
	defer fixtures.CleanupRuntime(t, ctx, certSecuredGraphQLClient, tenantId, &runtime)
	require.NoError(t, err)
	require.NotEmpty(t, runtime.ID)

	rtmCtx1 := fixtures.CreateRuntimeContext(t, ctx, certSecuredGraphQLClient, tenantId, runtime.ID, "queryRuntimeContexts1", "queryRuntimeContexts1")
	defer fixtures.DeleteRuntimeContext(t, ctx, certSecuredGraphQLClient, tenantId, rtmCtx1.ID)

	rtmCtx2 := fixtures.CreateRuntimeContext(t, ctx, certSecuredGraphQLClient, tenantId, runtime.ID, "queryRuntimeContexts2", "queryRuntimeContexts2")
	defer fixtures.DeleteRuntimeContext(t, ctx, certSecuredGraphQLClient, tenantId, rtmCtx2.ID)

	rtmCtxsRequest := fixtures.FixGetRuntimeContextsRequest(runtime.ID)
	runtimeGql := graphql.RuntimeExt{}

	err = testctx.Tc.RunOperation(ctx, certSecuredGraphQLClient, rtmCtxsRequest, &runtimeGql)
	require.NoError(t, err)
	require.Equal(t, 2, len(runtimeGql.RuntimeContexts.Data))
	require.ElementsMatch(t, []*graphql.RuntimeContextExt{&rtmCtx1, &rtmCtx2}, runtimeGql.RuntimeContexts.Data)

	saveExample(t, rtmCtxsRequest.Query(), "query runtime contexts")
}

func TestUpdateRuntimeContext(t *testing.T) {
	ctx := context.Background()

	tenantId := tenant.TestTenants.GetDefaultTenantID()

	in := fixRuntimeInput("addRuntimeContext")

	runtime, err := fixtures.RegisterRuntimeFromInputWithinTenant(t, ctx, certSecuredGraphQLClient, tenantId, &in)
	defer fixtures.CleanupRuntime(t, ctx, certSecuredGraphQLClient, tenantId, &runtime)
	require.NoError(t, err)
	require.NotEmpty(t, runtime.ID)

	rtmCtx := fixtures.CreateRuntimeContext(t, ctx, certSecuredGraphQLClient, tenantId, runtime.ID, "runtimeContext", "runtimeContext")
	defer fixtures.DeleteRuntimeContext(t, ctx, certSecuredGraphQLClient, tenantId, rtmCtx.ID)

	rtmCtxUpdateInput := fixtures.FixRuntimeContextInput("updateRuntimeContext", "updateRuntimeContext")
	rtmCtxUpdate, err := testctx.Tc.Graphqlizer.RuntimeContextInputToGQL(rtmCtxUpdateInput)
	require.NoError(t, err)

	updateRtmCtxReq := fixtures.FixUpdateRuntimeContextRequest(rtmCtx.ID, rtmCtxUpdate)
	runtimeContext := graphql.RuntimeContextExt{}

	// WHEN
	t.Log("Update runtime context")
	err = testctx.Tc.RunOperation(ctx, certSecuredGraphQLClient, updateRtmCtxReq, &runtimeContext)

	// THEN
	require.NoError(t, err)
	require.NotEmpty(t, runtimeContext.ID)

	assertions.AssertRuntimeContext(t, &rtmCtxUpdateInput, &runtimeContext)
	saveExample(t, updateRtmCtxReq.Query(), "update runtime context")
}

func TestDeleteRuntimeContext(t *testing.T) {
	ctx := context.Background()

	tenantId := tenant.TestTenants.GetDefaultTenantID()

	in := fixRuntimeInput("addRuntimeContext")

	runtime, err := fixtures.RegisterRuntimeFromInputWithinTenant(t, ctx, certSecuredGraphQLClient, tenantId, &in)
	defer fixtures.CleanupRuntime(t, ctx, certSecuredGraphQLClient, tenantId, &runtime)
	require.NoError(t, err)
	require.NotEmpty(t, runtime.ID)

	rtmCtx := fixtures.CreateRuntimeContext(t, ctx, certSecuredGraphQLClient, tenantId, runtime.ID, "deleteRuntimeContext", "deleteRuntimeContext")

	rtmCtxDeleteReq := fixtures.FixDeleteRuntimeContextRequest(rtmCtx.ID)
	rtmCtxGql := graphql.RuntimeContext{}

	// WHEN
	t.Log("Delete runtimeContext")
	err = testctx.Tc.RunOperation(ctx, certSecuredGraphQLClient, rtmCtxDeleteReq, &rtmCtxGql)

	// THEN
	require.NoError(t, err)
	require.NotEmpty(t, rtmCtxGql.ID)
	require.Equal(t, "deleteRuntimeContext", rtmCtxGql.Key)
	require.Equal(t, "deleteRuntimeContext", rtmCtxGql.Value)

	saveExample(t, rtmCtxDeleteReq.Query(), "delete runtime context")
}

func TestRuntimeContextSubscriptionFlows(stdT *testing.T) {
	testCases := []struct {
		FirstSubscriptionFlow       string // For the first subscription in the test. Used by external services mock to determine to which saas application is the subscription.
		SecondSubscriptionFlow      string // For the second subscription in the tests. Used by external services mock to determine to which saas application is the subscription.
		SubscriptionProviderID      string
		SubscribedToProviderAppName string // this is provider app name which we subscribe for
		ExpectedProviderAppName     string // this is the value of runtimeType label of the runtime from which we create runtime context. For the tests where we are subscribing for an indirect dependency this should be the app name of the direct dependency.
	}{
		{
			FirstSubscriptionFlow:       conf.SubscriptionConfig.StandardFlow,
			SecondSubscriptionFlow:      conf.SubscriptionConfig.StandardFlow,
			SubscriptionProviderID:      conf.SubscriptionConfig.SelfRegDistinguishLabelValue,
			SubscribedToProviderAppName: conf.SubscriptionProviderAppNameValue,
			ExpectedProviderAppName:     conf.SubscriptionProviderAppNameValue,
		},
		{
			FirstSubscriptionFlow:       conf.SubscriptionConfig.IndirectDependencyFlow,
			SecondSubscriptionFlow:      conf.SubscriptionConfig.DirectDependencyFlow,
			SubscriptionProviderID:      conf.SelfRegisterDirectDependencyDistinguishLabelValue,
			SubscribedToProviderAppName: conf.IndirectDependencySubscriptionProviderAppNameValue,
			ExpectedProviderAppName:     conf.DirectDependencySubscriptionProviderAppNameValue,
		},
	}

	t := testingx.NewT(stdT)
	for _, testCase := range testCases {
		t.Run("Runtime Contexts subscription flows", func(t *testing.T) {
			ctx := context.Background()
			subscriptionProviderSubaccountID := conf.TestProviderSubaccountID // in local set up the parent is testDefaultTenant
			subscriptionConsumerAccountID := conf.TestConsumerAccountID
			subscriptionConsumerSubaccountID := conf.TestConsumerSubaccountID // in local set up the parent is ApplicationsForRuntimeTenantName
			subscriptionConsumerTenantID := conf.TestConsumerTenantID

			// Prepare provider external client certificate and secret and Build graphql director client configured with certificate
			providerClientKey, providerRawCertChain := certprovider.NewExternalCertFromConfig(t, ctx, conf.ExternalCertProviderConfig, true)
			directorCertSecuredClient := gql.NewCertAuthorizedGraphQLClientWithCustomURL(conf.DirectorExternalCertSecuredURL, providerClientKey, providerRawCertChain, conf.SkipSSLValidation)

			providerRuntimeInput := graphql.RuntimeRegisterInput{
				Name:        "providerRuntime",
				Description: ptr.String("providerRuntime-description"),
				Labels:      graphql.Labels{conf.SubscriptionConfig.SelfRegDistinguishLabelKey: testCase.SubscriptionProviderID},
			}

			var providerRuntime graphql.RuntimeExt // needed so the 'defer' can be above the runtime registration
			defer fixtures.CleanupRuntimeWithoutTenant(t, ctx, directorCertSecuredClient, &providerRuntime)
			providerRuntime = fixtures.RegisterRuntimeFromInputWithoutTenant(t, ctx, directorCertSecuredClient, &providerRuntimeInput)
			require.NotEmpty(t, providerRuntime.ID)

			selfRegLabelValue, ok := providerRuntime.Labels[conf.SubscriptionConfig.SelfRegisterLabelKey].(string)
			require.True(t, ok)
			require.Contains(t, selfRegLabelValue, conf.SubscriptionConfig.SelfRegisterLabelValuePrefix+providerRuntime.ID)

			regionLbl, ok := providerRuntime.Labels[tenantfetcher.RegionKey].(string)
			require.True(t, ok)
			require.Equal(t, conf.SubscriptionConfig.SelfRegRegion, regionLbl)

			saasAppLbl, ok := providerRuntime.Labels[conf.SaaSAppNameLabelKey].(string)
			require.True(t, ok)
			require.NotEmpty(t, saasAppLbl)

			httpClient := &http.Client{
				Timeout: 10 * time.Second,
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.SkipSSLValidation},
				},
			}

			depConfigureReq, err := http.NewRequest(http.MethodPost, conf.ExternalServicesMockBaseURL+"/v1/dependencies/configure", bytes.NewBuffer([]byte(selfRegLabelValue)))
			require.NoError(t, err)
			response, err := httpClient.Do(depConfigureReq)
			defer func() {
				if err := response.Body.Close(); err != nil {
					t.Logf("Could not close response body %s", err)
				}
			}()
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, response.StatusCode)

			subscriptionToken := token.GetClientCredentialsToken(t, ctx, conf.SubscriptionConfig.TokenURL+conf.TokenPath, conf.SubscriptionConfig.ClientID, conf.SubscriptionConfig.ClientSecret, "tenantFetcherClaims")
			apiPath := fmt.Sprintf("/saas-manager/v1/applications/%s/subscription", testCase.SubscribedToProviderAppName)
			defer subscription.BuildAndExecuteUnsubscribeRequest(t, providerRuntime.ID, providerRuntime.Name, httpClient, conf.SubscriptionConfig.URL, apiPath, subscriptionToken, conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, subscriptionProviderSubaccountID, testCase.FirstSubscriptionFlow, conf.SubscriptionConfig.SubscriptionFlowHeaderKey)
			createRuntimeSubscription(t, httpClient, providerRuntime, subscriptionToken, apiPath, subscriptionConsumerTenantID, subscriptionConsumerSubaccountID, subscriptionProviderSubaccountID, testCase.SubscribedToProviderAppName, true, testCase.FirstSubscriptionFlow)

			t.Log("Assert provider runtime is visible in the consumer's subaccount after successful subscription")
			consumerSubaccountRuntime := fixtures.GetRuntime(t, ctx, certSecuredGraphQLClient, subscriptionConsumerSubaccountID, providerRuntime.ID)
			require.Equal(t, providerRuntime.ID, consumerSubaccountRuntime.ID)

			t.Log("Assert subscription provider application name label of the provider runtime exists and it is the correct one")
			subProviderAppNameValue, ok := consumerSubaccountRuntime.Labels[conf.RuntimeTypeLabelKey].(string)
			require.True(t, ok)
			require.Equal(t, testCase.ExpectedProviderAppName, subProviderAppNameValue)

			t.Log("Assert there is a runtime context(subscription) as part of the provider runtime")
			require.Len(t, consumerSubaccountRuntime.RuntimeContexts.Data, 1)
			require.Equal(t, conf.SubscriptionLabelKey, consumerSubaccountRuntime.RuntimeContexts.Data[0].Key)
			require.Equal(t, subscriptionConsumerTenantID, consumerSubaccountRuntime.RuntimeContexts.Data[0].Value)

			t.Log("Assert the runtime context has label containing consumer subaccount ID")
			consumerSubaccountFromRtmCtxLabel, ok := consumerSubaccountRuntime.RuntimeContexts.Data[0].Labels[conf.GlobalSubaccountIDLabelKey].(string)
			require.True(t, ok)
			require.Equal(t, subscriptionConsumerSubaccountID, consumerSubaccountFromRtmCtxLabel)

			t.Log("Assert provider runtime is visible in the consumer's account after successful subscription")
			consumerAccountRuntime := fixtures.GetRuntime(t, ctx, certSecuredGraphQLClient, subscriptionConsumerAccountID, providerRuntime.ID)
			require.Equal(t, providerRuntime.ID, consumerAccountRuntime.ID)
			require.Len(t, consumerSubaccountRuntime.RuntimeContexts.Data, 1)

			t.Log("Assert the consumer cannot update the provider runtime(owner false check)")
			consumerRuntimeUpdateInput := fixRuntimeUpdateInput("consumerUpdatedRuntime")
			rtm, err := fixtures.UpdateRuntimeWithinTenant(t, ctx, certSecuredGraphQLClient, subscriptionConsumerSubaccountID, providerRuntime.ID, consumerRuntimeUpdateInput)
			require.Empty(t, rtm)
			require.Error(t, err)
			require.Contains(t, err.Error(), "Owner access is needed for resource modification")

			t.Log("Assert the consumer cannot delete the provider runtime(cleanup of self-registered runtime failure)")
			deleteRuntimeReq := fixtures.FixUnregisterRuntimeRequest(providerRuntime.ID)
			rtmExt := graphql.RuntimeExt{}
			err = testctx.Tc.RunOperationWithCustomTenant(ctx, certSecuredGraphQLClient, subscriptionConsumerSubaccountID, deleteRuntimeReq, &rtmExt)
			require.Empty(t, rtmExt)
			require.Error(t, err)
			// We shouldn't be able cleanup the self-registered runtime if there is a subscription
			require.Contains(t, err.Error(), "received unexpected status code 409")

			subscription.BuildAndExecuteUnsubscribeRequest(t, providerRuntime.ID, providerRuntime.Name, httpClient, conf.SubscriptionConfig.URL, apiPath, subscriptionToken, conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, subscriptionProviderSubaccountID, testCase.FirstSubscriptionFlow, conf.SubscriptionConfig.SubscriptionFlowHeaderKey)

			t.Log("List runtimes(and runtime contexts) after successful unsubscribe request")
			consumerSubaccountRtms := fixtures.ListRuntimes(t, ctx, certSecuredGraphQLClient, subscriptionConsumerSubaccountID)
			require.Len(t, consumerSubaccountRtms.Data, 1)
			require.Equal(t, consumerSubaccountRtms.Data[0].ID, providerRuntime.ID)

			t.Log("Assert there is no runtime context(subscription) after successful unsubscribe request")
			require.Len(t, consumerSubaccountRtms.Data[0].RuntimeContexts.Data, 0)
		})

		t.Run("Runtime Contexts subscription flow - validate subscriptions label", func(t *testing.T) {
			ctx := context.Background()
			subscriptionProviderSubaccountID := conf.TestProviderSubaccountID // in local set up the parent is testDefaultTenant
			subscriptionConsumerSubaccountID := conf.TestConsumerSubaccountID // in local set up the parent is ApplicationsForRuntimeTenantName
			subscriptionConsumerTenantID := conf.TestConsumerTenantID

			// Prepare provider external client certificate and secret and Build graphql director client configured with certificate
			providerClientKey, providerRawCertChain := certprovider.NewExternalCertFromConfig(t, ctx, conf.ExternalCertProviderConfig, true)
			directorCertSecuredClient := gql.NewCertAuthorizedGraphQLClientWithCustomURL(conf.DirectorExternalCertSecuredURL, providerClientKey, providerRawCertChain, conf.SkipSSLValidation)

			providerRuntimeInput := graphql.RuntimeRegisterInput{
				Name:        "providerRuntime",
				Description: ptr.String("providerRuntime-description"),
				Labels:      graphql.Labels{conf.SubscriptionConfig.SelfRegDistinguishLabelKey: testCase.SubscriptionProviderID},
			}

			providerRuntime := fixtures.RegisterRuntimeFromInputWithoutTenant(t, ctx, directorCertSecuredClient, &providerRuntimeInput)
			defer fixtures.CleanupRuntimeWithoutTenant(t, ctx, directorCertSecuredClient, &providerRuntime)
			require.NotEmpty(t, providerRuntime.ID)

			selfRegLabelValue, ok := providerRuntime.Labels[conf.SubscriptionConfig.SelfRegisterLabelKey].(string)
			require.True(t, ok)
			require.Contains(t, selfRegLabelValue, conf.SubscriptionConfig.SelfRegisterLabelValuePrefix+providerRuntime.ID)

			regionLbl, ok := providerRuntime.Labels[tenantfetcher.RegionKey].(string)
			require.True(t, ok)
			require.Equal(t, conf.SubscriptionConfig.SelfRegRegion, regionLbl)

			saasAppLbl, ok := providerRuntime.Labels[conf.SaaSAppNameLabelKey].(string)
			require.True(t, ok)
			require.NotEmpty(t, saasAppLbl)

			httpClient := &http.Client{
				Timeout: 10 * time.Second,
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.SkipSSLValidation},
				},
			}

			depConfigureReq, err := http.NewRequest(http.MethodPost, conf.ExternalServicesMockBaseURL+"/v1/dependencies/configure", bytes.NewBuffer([]byte(selfRegLabelValue)))
			require.NoError(t, err)
			response, err := httpClient.Do(depConfigureReq)
			defer func() {
				if err := response.Body.Close(); err != nil {
					t.Logf("Could not close response body %s", err)
				}
			}()
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, response.StatusCode)

			subscriptionToken := token.GetClientCredentialsToken(t, ctx, conf.SubscriptionConfig.TokenURL+conf.TokenPath, conf.SubscriptionConfig.ClientID, conf.SubscriptionConfig.ClientSecret, "tenantFetcherClaims")
			apiPath := fmt.Sprintf("/saas-manager/v1/applications/%s/subscription", testCase.SubscribedToProviderAppName)

			defer subscription.BuildAndExecuteUnsubscribeRequest(t, providerRuntime.ID, providerRuntime.Name, httpClient, conf.SubscriptionConfig.URL, apiPath, subscriptionToken, conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, subscriptionProviderSubaccountID, testCase.FirstSubscriptionFlow, conf.SubscriptionConfig.SubscriptionFlowHeaderKey)
			createRuntimeSubscription(t, httpClient, providerRuntime, subscriptionToken, apiPath, subscriptionConsumerTenantID, subscriptionConsumerSubaccountID, subscriptionProviderSubaccountID, testCase.SubscribedToProviderAppName, true, testCase.FirstSubscriptionFlow)

			t.Log("Assert the runtime context has subscriptions label")
			assertRuntimeContextFromSubscription(t, subscriptionConsumerSubaccountID, providerRuntime.ID, 1)

			t.Logf("Creating a second subscription between consumer with subaccount id: %q and tenant id: %q, and provider with name: %q, id: %q and subaccount id: %q", subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, providerRuntime.Name, providerRuntime.ID, subscriptionProviderSubaccountID)

			defer subscription.BuildAndExecuteUnsubscribeRequest(t, providerRuntime.ID, providerRuntime.Name, httpClient, conf.SubscriptionConfig.URL, apiPath, subscriptionToken, conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, subscriptionProviderSubaccountID, testCase.SecondSubscriptionFlow, conf.SubscriptionConfig.SubscriptionFlowHeaderKey)
			createRuntimeSubscription(t, httpClient, providerRuntime, subscriptionToken, apiPath, subscriptionConsumerTenantID, subscriptionConsumerSubaccountID, subscriptionProviderSubaccountID, testCase.SubscribedToProviderAppName, false, testCase.SecondSubscriptionFlow)
			t.Logf("Successfully created second subscription between consumer with subaccount id: %q and tenant id: %q, and provider with name: %q, id: %q and subaccount id: %q", subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, providerRuntime.Name, providerRuntime.ID, subscriptionProviderSubaccountID)

			t.Log("Assert runtime context subscriptions label has new value added")
			assertRuntimeContextFromSubscription(t, subscriptionConsumerSubaccountID, providerRuntime.ID, 2)

			subscription.BuildAndExecuteUnsubscribeRequest(t, providerRuntime.ID, providerRuntime.Name, httpClient, conf.SubscriptionConfig.URL, apiPath, subscriptionToken, conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, subscriptionProviderSubaccountID, testCase.FirstSubscriptionFlow, conf.SubscriptionConfig.SubscriptionFlowHeaderKey)

			t.Log("Assert runtime context subscriptions label has one value less")
			assertRuntimeContextFromSubscription(t, subscriptionConsumerSubaccountID, providerRuntime.ID, 1)

			subscription.BuildAndExecuteUnsubscribeRequest(t, providerRuntime.ID, providerRuntime.Name, httpClient, conf.SubscriptionConfig.URL, apiPath, subscriptionToken, conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, subscriptionProviderSubaccountID, testCase.SecondSubscriptionFlow, conf.SubscriptionConfig.SubscriptionFlowHeaderKey)

			t.Log("List runtimes(and runtime contexts) after successful unsubscribe request")
			consumerSubaccountRtms := fixtures.ListRuntimes(t, ctx, certSecuredGraphQLClient, subscriptionConsumerSubaccountID)
			require.Len(t, consumerSubaccountRtms.Data, 1)
			require.Equal(t, consumerSubaccountRtms.Data[0].ID, providerRuntime.ID)

			t.Log("Assert there is no runtime context(subscription) after successful unsubscribe request")
			require.Len(t, consumerSubaccountRtms.Data[0].RuntimeContexts.Data, 0)
		})
	}
}

func assertRuntimeContextFromSubscription(t *testing.T, tenantID, runtimeID string, expectedSubscriptionsCount int) {
	consumerSubaccountRuntime := fixtures.GetRuntime(t, ctx, certSecuredGraphQLClient, tenantID, runtimeID)
	require.Len(t, consumerSubaccountRuntime.RuntimeContexts.Data, 1)
	subscriptionsLabelInterface, ok := consumerSubaccountRuntime.RuntimeContexts.Data[0].Labels[conf.SubscriptionConfig.SubscriptionsLabelKey].([]interface{})
	require.True(t, ok)
	subscriptionsLabelValue := make([]string, len(subscriptionsLabelInterface))
	for i, v := range subscriptionsLabelInterface {
		subscriptionsLabelValue[i], ok = v.(string)
		require.True(t, ok)
	}
	require.Len(t, subscriptionsLabelValue, expectedSubscriptionsCount)
}

func createRuntimeSubscription(t *testing.T, httpClient *http.Client, providerRuntime graphql.RuntimeExt, subscriptionToken, apiPath, subscriptionConsumerTenantID, subscriptionConsumerSubaccountID, subscriptionProviderSubaccountID, subscriptionProviderAppNameValue string, shouldUnsubscribeFirst bool, subscriptionFlow string) {
	subscribeReq := subscription.BuildSubscriptionRequest(t, subscriptionToken, conf.SubscriptionConfig.URL, subscriptionProviderSubaccountID, subscriptionProviderAppNameValue, conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionFlow, conf.SubscriptionConfig.SubscriptionFlowHeaderKey)

	if shouldUnsubscribeFirst {
		// unsubscribe request execution to ensure no resources/subscriptions are left unintentionally due to old unsubscribe failures or broken tests in the middle.
		// In case there isn't subscription it will fail-safe without error
		subscription.BuildAndExecuteUnsubscribeRequest(t, providerRuntime.ID, providerRuntime.Name, httpClient, conf.SubscriptionConfig.URL, apiPath, subscriptionToken, conf.SubscriptionConfig.PropagatedProviderSubaccountHeader, subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, subscriptionProviderSubaccountID, conf.SubscriptionConfig.StandardFlow, conf.SubscriptionConfig.SubscriptionFlowHeaderKey)
	}

	t.Logf("Creating a subscription between consumer with subaccount id: %q and tenant id: %q, and provider with name: %q, id: %q and subaccount id: %q", subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, providerRuntime.Name, providerRuntime.ID, subscriptionProviderSubaccountID)
	resp, err := httpClient.Do(subscribeReq)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Could not close response body %s", err)
		}
	}()
	require.NoError(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, resp.StatusCode, fmt.Sprintf("actual status code %d is different from the expected one: %d. Reason: %v", resp.StatusCode, http.StatusAccepted, string(body)))

	subJobStatusPath := resp.Header.Get(subscription.LocationHeader)
	require.NotEmpty(t, subJobStatusPath)
	subJobStatusURL := conf.SubscriptionConfig.URL + subJobStatusPath
	require.Eventually(t, func() bool {
		return subscription.GetSubscriptionJobStatus(t, httpClient, subJobStatusURL, subscriptionToken) == subscription.JobSucceededStatus
	}, subscription.EventuallyTimeout, subscription.EventuallyTick)
	t.Logf("Successfully created subscription between consumer with subaccount id: %q and tenant id: %q, and provider with name: %q, id: %q and subaccount id: %q", subscriptionConsumerSubaccountID, subscriptionConsumerTenantID, providerRuntime.Name, providerRuntime.ID, subscriptionProviderSubaccountID)
}
