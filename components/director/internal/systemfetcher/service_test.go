package systemfetcher_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/tidwall/gjson"

	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/systemfetcher"
	"github.com/kyma-incubator/compass/components/director/internal/systemfetcher/automock"
	pAutomock "github.com/kyma-incubator/compass/components/director/pkg/persistence/automock"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence/txtest"
	tenantEntity "github.com/kyma-incubator/compass/components/director/pkg/tenant"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	appType    = "type-1"
	mainURLKey = "mainUrl"
	testID     = "testID"
)

func TestSyncSystems(t *testing.T) {
	const appTemplateID = "appTmp1"
	testErr := errors.New("testErr")
	setupSuccessfulTemplateRenderer := func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInput) *automock.TemplateRenderer {
		svc := &automock.TemplateRenderer{}
		for i := range appsInputs {
			svc.On("ApplicationRegisterInputFromTemplate", txtest.CtxWithDBMatcher(), systems[i]).Return(&appsInputs[i], nil).Once()
		}
		return svc
	}

	var mutex sync.Mutex
	type testCase struct {
		name                     string
		mockTransactioner        func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner)
		fixTestSystems           func() []systemfetcher.System
		fixAppInputs             func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate
		setupTenantSvc           func() *automock.TenantService
		setupTbtSvc              func() *automock.TenantBusinessTypeService
		setupTemplateRendererSvc func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInput) *automock.TemplateRenderer
		setupSystemSvc           func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService
		setupSystemsSyncSvc      func() *automock.SystemsSyncService
		setupSysAPIClient        func(testSystems []systemfetcher.System) *automock.SystemsAPIClient
		setupDirectorClient      func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient
		verificationTenant       string
		expectedErr              error
	}
	tests := []testCase{
		{
			name: "Success with one tenant and one system",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(3)
				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystems()
				systems[0].TemplateID = appTemplateID
				systems[0].SystemPayload["productId"] = "TEST"
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: setupSuccessfulTemplateRenderer,
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), appsInputs[0].ApplicationRegisterInput, mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(testSystems, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Success with one tenant and one system which has tbt and tbt does not exist in the db",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(3)
				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystemsWithTbt()
				systems[0].TemplateID = appTemplateID
				systems[0].SystemPayload["productId"] = "TEST"
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil).Once()
				tbtSvc.On("Create", txtest.CtxWithDBMatcher(), &model.TenantBusinessTypeInput{Code: fixSystemsWithTbt()[0].SystemPayload["businessTypeId"].(string), Name: fixSystemsWithTbt()[0].SystemPayload["businessTypeDescription"].(string)}).Return(testID, nil).Once()
				tbtSvc.On("GetByID", txtest.CtxWithDBMatcher(), testID).Return(&model.TenantBusinessType{ID: testID, Code: fixSystemsWithTbt()[0].SystemPayload["businessTypeId"].(string), Name: fixSystemsWithTbt()[0].SystemPayload["businessTypeDescription"].(string)}, nil).Once()
				return tbtSvc
			},
			setupTemplateRendererSvc: setupSuccessfulTemplateRenderer,
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				appsInputs[0].ApplicationRegisterInput.TenantBusinessTypeID = str.Ptr(testID)
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), appsInputs[0].ApplicationRegisterInput, mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(testSystems, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Success with one tenant and one system which has tbt and tbt exists in the db",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(3)
				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystemsWithTbt()
				systems[0].TemplateID = appTemplateID
				systems[0].SystemPayload["productId"] = "TEST"
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				systemPayload, err := json.Marshal(fixSystemsWithTbt()[0].SystemPayload)
				require.NoError(t, err)

				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{{ID: testID, Code: gjson.GetBytes(systemPayload, "businessTypeId").String(), Name: gjson.GetBytes(systemPayload, "BusinessTypeDescription").String()}}, nil).Once()
				return tbtSvc
			},
			setupTemplateRendererSvc: setupSuccessfulTemplateRenderer,
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				appsInputs[0].ApplicationRegisterInput.TenantBusinessTypeID = str.Ptr(testID)
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), appsInputs[0].ApplicationRegisterInput, mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(testSystems, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Success when in verification mode",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(3)
				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystems()
				systems[0].TemplateID = appTemplateID
				systems[0].SystemPayload["productId"] = "TEST"
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("GetTenantByExternalID", txtest.CtxWithDBMatcher(), "t1").Return(tenants[0], nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: setupSuccessfulTemplateRenderer,
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), appsInputs[0].ApplicationRegisterInput, mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(testSystems, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
			verificationTenant: "t1",
		},
		{
			name: "Success with one tenant and one system that has already been in the database and will not have it's status condition changed",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(3)
				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystems()
				systems[0].TemplateID = appTemplateID
				systems[0].SystemPayload["productId"] = "TEST"
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				svc := &automock.TemplateRenderer{}
				connectedStatus := model.ApplicationStatusConditionConnected

				for i := range appsInputs {
					input := systems[i]
					input.StatusCondition = connectedStatus

					result := appsInputs[i]
					result.StatusCondition = &connectedStatus

					svc.On("ApplicationRegisterInputFromTemplate", txtest.CtxWithDBMatcher(), input).Return(&result, nil).Once()
				}
				return svc
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}

				connectedStatus := model.ApplicationStatusConditionConnected
				appInput := appsInputs[0].ApplicationRegisterInput
				appInput.StatusCondition = &connectedStatus

				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), appInput, mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
					Status: &model.ApplicationStatus{
						Condition: model.ApplicationStatusConditionConnected,
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(testSystems, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Success with one tenant and one system with null base url",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(3)
				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := []systemfetcher.System{
					{
						SystemPayload: map[string]interface{}{
							"displayName":            "System1",
							"productDescription":     "System1 description",
							"infrastructureProvider": "test",
						},
						StatusCondition: model.ApplicationStatusConditionInitial,
					},
				}
				systems[0].TemplateID = "type1"
				systems[0].SystemPayload["productId"] = "TEST"
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: setupSuccessfulTemplateRenderer,
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), mock.AnythingOfType("model.ApplicationRegisterInput"), mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(testSystems, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Success with one tenant and one system without template",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(3)
				return mockedTx, transactioner
			},
			fixTestSystems: fixSystems,
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: func(_ []systemfetcher.System, _ []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				systemSvc.On("TrustedUpsert", txtest.CtxWithDBMatcher(), appsInputs[0].ApplicationRegisterInput, mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(testSystems, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Success with one tenant and multiple systems",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(4)
				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystems()
				systems[0].TemplateID = appTemplateID
				systems = append(systems, systemfetcher.System{
					SystemPayload: map[string]interface{}{
						"displayName":            "System2",
						"productDescription":     "System2 description",
						"baseUrl":                "http://example2.com",
						"infrastructureProvider": "test",
					},
					TemplateID:      "appTmp2",
					StatusCondition: model.ApplicationStatusConditionInitial,
				})
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: setupSuccessfulTemplateRenderer,
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), appsInputs[0].ApplicationRegisterInput, mock.Anything).Return(nil).Once()
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), appsInputs[1].ApplicationRegisterInput, mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[1].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(testSystems, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Success with multiple tenants with one system",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(4)
				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystems()
				systems[0].TemplateID = appTemplateID
				systems = append(systems, systemfetcher.System{
					SystemPayload: map[string]interface{}{
						"displayName":            "System2",
						"productDescription":     "System2 description",
						"baseUrl":                "http://example2.com",
						"infrastructureProvider": "test",
					},
					TemplateID:      "appTmp2",
					StatusCondition: model.ApplicationStatusConditionInitial,
				})
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				firstTenant := newModelBusinessTenantMapping("t1", "tenant1")
				firstTenant.ExternalTenant = "t1"
				secondTenant := newModelBusinessTenantMapping("t2", "tenant2")
				secondTenant.ExternalTenant = "t2"
				tenants := []*model.BusinessTenantMapping{firstTenant, secondTenant}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: setupSuccessfulTemplateRenderer,
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), appsInputs[0].ApplicationRegisterInput, mock.Anything).Return(nil).Once()
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), appsInputs[1].ApplicationRegisterInput, mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[1].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "t1", &mutex).Return([]systemfetcher.System{testSystems[0]}, nil).Once()
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "t2", &mutex).Return([]systemfetcher.System{testSystems[1]}, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Fail when listing all tenant business types",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				persistTx := &pAutomock.PersistenceTx{}
				persistTx.On("Commit").Return(nil).Once()

				transact := &pAutomock.Transactioner{}
				transact.On("Begin").Return(persistTx, nil).Twice()
				transact.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(false).Once()
				transact.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(true).Once()

				return persistTx, transact
			},
			fixTestSystems: func() []systemfetcher.System {
				return fixSystemsWithTbt()
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return(nil, testErr).Once()
				return tbtSvc
			},
			setupTemplateRendererSvc: func(_ []systemfetcher.System, _ []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				return &automock.SystemsService{}
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				return &automock.SystemsAPIClient{}
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
			expectedErr: testErr,
		},
		{
			name: "Begin transaction fails when listing all tenant business types",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				persistTx := &pAutomock.PersistenceTx{}
				persistTx.On("Commit").Return(nil).Once()

				transact := &pAutomock.Transactioner{}
				transact.On("Begin").Return(persistTx, nil).Once()
				transact.On("Begin").Return(persistTx, testErr).Once()
				transact.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(false).Once()

				return persistTx, transact
			},
			fixTestSystems: func() []systemfetcher.System {
				return fixSystemsWithTbt()
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				return &automock.TenantBusinessTypeService{}
			},
			setupTemplateRendererSvc: func(_ []systemfetcher.System, _ []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				return &automock.SystemsService{}
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				return &automock.SystemsAPIClient{}
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
			expectedErr: testErr,
		},
		{
			name: "Commit transaction fails when listing all tenant business types",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				persistTx := &pAutomock.PersistenceTx{}
				persistTx.On("Commit").Return(nil).Once()
				persistTx.On("Commit").Return(testErr).Once()

				transact := &pAutomock.Transactioner{}
				transact.On("Begin").Return(persistTx, nil).Twice()
				transact.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(false).Once()
				transact.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(true).Once()

				return persistTx, transact
			},
			fixTestSystems: func() []systemfetcher.System {
				return fixSystemsWithTbt()
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{{ID: testID, Code: fixSystemsWithTbt()[0].SystemPayload["businessTypeId"].(string), Name: fixSystemsWithTbt()[0].SystemPayload["businessTypeDescription"].(string)}}, nil).Once()
				return tbtSvc
			},
			setupTemplateRendererSvc: func(_ []systemfetcher.System, _ []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				return &automock.SystemsService{}
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				return &automock.SystemsAPIClient{}
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
			expectedErr: testErr,
		},
		{
			name: "Fails when creating new tenant business type",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				persistTx := &pAutomock.PersistenceTx{}
				persistTx.On("Commit").Return(nil).Twice()

				transact := &pAutomock.Transactioner{}
				transact.On("Begin").Return(persistTx, nil).Times(3)
				transact.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(false).Twice()
				transact.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(true).Once()

				return persistTx, transact
			},
			fixTestSystems: func() []systemfetcher.System {
				return fixSystemsWithTbt()
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil).Once()
				tbtSvc.On("Create", txtest.CtxWithDBMatcher(), &model.TenantBusinessTypeInput{Code: fixSystemsWithTbt()[0].SystemPayload["businessTypeId"].(string), Name: fixSystemsWithTbt()[0].SystemPayload["businessTypeDescription"].(string)}).Return("", testErr).Once()
				return tbtSvc
			},
			setupTemplateRendererSvc: func(_ []systemfetcher.System, _ []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(testSystems, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Fails when getting by id newly created tenant business type",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				persistTx := &pAutomock.PersistenceTx{}
				persistTx.On("Commit").Return(nil).Twice()

				transact := &pAutomock.Transactioner{}
				transact.On("Begin").Return(persistTx, nil).Times(3)
				transact.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(false).Twice()
				transact.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(true).Once()

				return persistTx, transact
			},
			fixTestSystems: func() []systemfetcher.System {
				return fixSystemsWithTbt()
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil).Once()
				tbtSvc.On("Create", txtest.CtxWithDBMatcher(), &model.TenantBusinessTypeInput{Code: fixSystemsWithTbt()[0].SystemPayload["businessTypeId"].(string), Name: fixSystemsWithTbt()[0].SystemPayload["businessTypeDescription"].(string)}).Return(testID, nil).Once()
				tbtSvc.On("GetByID", txtest.CtxWithDBMatcher(), testID).Return(nil, testErr).Once()
				return tbtSvc
			},
			setupTemplateRendererSvc: func(_ []systemfetcher.System, _ []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(testSystems, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Fail when tenant fetching fails",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatDoesntExpectCommit()

				return mockedTx, transactioner
			},
			fixTestSystems: fixSystems,
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(nil, testErr).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				return &automock.TenantBusinessTypeService{}
			},
			setupTemplateRendererSvc: func(_ []systemfetcher.System, _ []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				return &automock.SystemsService{}
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				return &automock.SystemsAPIClient{}
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
			expectedErr: testErr,
		},
		{
			name: "Fail when transaction cannot be started",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(testErr).ThatFailsOnBegin()

				return mockedTx, transactioner
			},
			fixTestSystems: fixSystems,
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				return &automock.TenantService{}
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				return &automock.TenantBusinessTypeService{}
			},
			setupTemplateRendererSvc: func(_ []systemfetcher.System, _ []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				return &automock.SystemsService{}
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				return &automock.SystemsAPIClient{}
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
			expectedErr: testErr,
		},
		{
			name: "Fail when commit fails",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(testErr).ThatFailsOnCommit()
				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystems()
				systems[0].TemplateID = appTemplateID
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				return &automock.TenantBusinessTypeService{}
			},
			setupTemplateRendererSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				return &automock.SystemsService{}
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				return &automock.SystemsAPIClient{}
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
			expectedErr: testErr,
		},
		{
			name: "Fail when client fails to fetch systems",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(2)

				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				return fixSystems()
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: func(_ []systemfetcher.System, _ []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				return &automock.SystemsService{}
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(nil, errors.New("expected")).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Fail when service fails to save systems",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(2)
				persistTx := &pAutomock.PersistenceTx{}

				transactioner.On("Begin").Return(persistTx, nil).Once()
				persistTx.On("Commit").Return(nil).Once()
				transactioner.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(true).Once()

				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystems()
				systems[0].TemplateID = appTemplateID
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: setupSuccessfulTemplateRenderer,
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), appsInputs[0].ApplicationRegisterInput, mock.Anything).Return(errors.New("expected")).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return([]systemfetcher.System{testSystems[0]}, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Fail when application from template cannot be rendered",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(2)
				persistTx := &pAutomock.PersistenceTx{}

				transactioner.On("Begin").Return(persistTx, nil).Once()
				transactioner.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(true).Once()

				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystems()
				systems[0].TemplateID = appTemplateID
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				svc := &automock.TemplateRenderer{}
				for i := range appsInputs {
					svc.On("ApplicationRegisterInputFromTemplate", txtest.CtxWithDBMatcher(), systems[i]).Return(nil, testErr).Once()
				}
				return svc
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)

				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return([]systemfetcher.System{testSystems[0]}, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Succeed when client fails to fetch systems only for some tenants",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(4)

				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystems()
				systems[0].TemplateID = "type1"
				systems = append(systems, systemfetcher.System{
					SystemPayload: map[string]interface{}{
						"displayName":            "System2",
						"productDescription":     "System2 description",
						"baseUrl":                "http://example2.com",
						"infrastructureProvider": "test",
					},
					TemplateID:      "type2",
					StatusCondition: model.ApplicationStatusConditionInitial,
				})
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				tenants := []*model.BusinessTenantMapping{
					newModelBusinessTenantMapping("t1", "tenant1"),
					newModelBusinessTenantMapping("t2", "tenant2"),
					newModelBusinessTenantMapping("t3", "tenant3"),
				}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: setupSuccessfulTemplateRenderer,
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), mock.AnythingOfType("model.ApplicationRegisterInput"), mock.Anything).Return(nil).Once()
				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), mock.AnythingOfType("model.ApplicationRegisterInput"), mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[1].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return([]systemfetcher.System{testSystems[0]}, nil).Once()
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return(nil, errors.New("expected")).Once()
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "external", &mutex).Return([]systemfetcher.System{testSystems[1]}, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Do nothing if system is already being deleted",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(2)
				persistTx := &pAutomock.PersistenceTx{}

				transactioner.On("Begin").Return(persistTx, nil).Twice()
				persistTx.On("Commit").Return(nil).Once()
				transactioner.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(true).Twice()

				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystems()
				systems[0].TemplateID = "type1"
				systems = append(systems, systemfetcher.System{
					SystemPayload: map[string]interface{}{
						"displayName":            "System2",
						"productDescription":     "System2 description",
						"baseUrl":                "http://example2.com",
						"infrastructureProvider": "test",
						"additionalAttributes": map[string]string{
							systemfetcher.LifecycleAttributeName: systemfetcher.LifecycleDeleted,
						},
						"systemNumber": "sysNumber1",
					},
					TemplateID:      "type2",
					StatusCondition: model.ApplicationStatusConditionInitial,
				})
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				firstTenant := newModelBusinessTenantMapping("t1", "tenant1")
				firstTenant.ExternalTenant = "t1"
				secondTenant := newModelBusinessTenantMapping("t2", "tenant2")
				secondTenant.ExternalTenant = "t2"
				tenants := []*model.BusinessTenantMapping{firstTenant, secondTenant}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				svc := &automock.TemplateRenderer{}
				appInput := appsInputs[0] // appsInputs[1] belongs to a system with status "DELETED"
				svc.On("ApplicationRegisterInputFromTemplate", txtest.CtxWithDBMatcher(), systems[0]).Return(&appInput, nil)
				return svc
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}

				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), mock.AnythingOfType("model.ApplicationRegisterInput"), mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[1].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "t1", &mutex).Return([]systemfetcher.System{testSystems[0]}, nil).Once()
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "t2", &mutex).Return([]systemfetcher.System{testSystems[1]}, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				directorClient := &automock.DirectorClient{}
				directorClient.On("DeleteSystemAsync", mock.Anything, "id", "t2").Return(nil).Once()
				return directorClient
			},
		},
		{
			name: "Do nothing if system has already been deleted",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(2)
				persistTx := &pAutomock.PersistenceTx{}

				transactioner.On("Begin").Return(persistTx, nil).Twice()
				persistTx.On("Commit").Return(nil).Once()
				transactioner.On("RollbackUnlessCommitted", mock.Anything, persistTx).Return(true).Twice()

				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				systems := fixSystems()
				systems[0].TemplateID = "type1"
				systems = append(systems, systemfetcher.System{
					SystemPayload: map[string]interface{}{
						"displayName":            "System2",
						"productDescription":     "System2 description",
						"baseUrl":                "http://example2.com",
						"infrastructureProvider": "test",
						"additionalAttributes": map[string]string{
							systemfetcher.LifecycleAttributeName: systemfetcher.LifecycleDeleted,
						},
						"systemNumber": "sysNumber1",
					},
					TemplateID:      "type2",
					StatusCondition: model.ApplicationStatusConditionInitial,
				})
				return systems
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return fixAppsInputsWithTemplatesBySystems(t, systems)
			},
			setupTenantSvc: func() *automock.TenantService {
				firstTenant := newModelBusinessTenantMapping("t1", "tenant1")
				firstTenant.ExternalTenant = "t1"
				secondTenant := newModelBusinessTenantMapping("t2", "tenant2")
				secondTenant.ExternalTenant = "t2"
				tenants := []*model.BusinessTenantMapping{firstTenant, secondTenant}
				tenantSvc := &automock.TenantService{}
				tenantSvc.On("ListByType", txtest.CtxWithDBMatcher(), tenantEntity.Account).Return(tenants, nil).Once()
				return tenantSvc
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				tbtSvc := &automock.TenantBusinessTypeService{}
				tbtSvc.On("ListAll", txtest.CtxWithDBMatcher()).Return([]*model.TenantBusinessType{}, nil)
				return tbtSvc
			},
			setupTemplateRendererSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				svc := &automock.TemplateRenderer{}
				appInput := appsInputs[0] // appsInputs[1] belongs to a system with status "DELETED"
				svc.On("ApplicationRegisterInputFromTemplate", txtest.CtxWithDBMatcher(), systems[0]).Return(&appInput, nil)
				return svc
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				systemSvc := &automock.SystemsService{}

				systemSvc.On("TrustedUpsertFromTemplate", txtest.CtxWithDBMatcher(), mock.AnythingOfType("model.ApplicationRegisterInput"), mock.Anything).Return(nil).Once()
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[0].SystemNumber).Return(&model.Application{
					BaseEntity: &model.BaseEntity{
						ID: "id",
					},
				}, nil)
				systemSvc.On("GetBySystemNumber", txtest.CtxWithDBMatcher(), *appsInputs[1].SystemNumber).Return(nil, nil)
				return systemSvc
			},
			setupSystemsSyncSvc: emptySystemsSyncSvc,
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				sysAPIClient := &automock.SystemsAPIClient{}
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "t1", &mutex).Return([]systemfetcher.System{testSystems[0]}, nil).Once()
				sysAPIClient.On("FetchSystemsForTenant", mock.Anything, "t2", &mutex).Return([]systemfetcher.System{testSystems[1]}, nil).Once()
				return sysAPIClient
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				directorClient := &automock.DirectorClient{}
				directorClient.On("DeleteSystemAsync", mock.Anything, "id", "t2").Return(nil).Once()
				return directorClient
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockedTx, transactioner := testCase.mockTransactioner()
			tenantSvc := testCase.setupTenantSvc()
			tbtSvc := testCase.setupTbtSvc()
			testSystems := testCase.fixTestSystems()
			appsInputs := testCase.fixAppInputs(testSystems)
			systemSvc := testCase.setupSystemSvc(testSystems, appsInputs)
			systemsSyncSvc := testCase.setupSystemsSyncSvc()
			appInputsWithoutTemplates := make([]model.ApplicationRegisterInput, 0)
			for _, in := range appsInputs {
				appInputsWithoutTemplates = append(appInputsWithoutTemplates, in.ApplicationRegisterInput)
			}
			templateAppResolver := testCase.setupTemplateRendererSvc(testSystems, appInputsWithoutTemplates)
			sysAPIClient := testCase.setupSysAPIClient(testSystems)
			directorClient := testCase.setupDirectorClient(testSystems, appsInputs)
			defer mock.AssertExpectationsForObjects(t, tenantSvc, tbtSvc, sysAPIClient, systemSvc, templateAppResolver, mockedTx, transactioner)

			svc := systemfetcher.NewSystemFetcher(transactioner, tenantSvc, systemSvc, systemsSyncSvc, tbtSvc, templateAppResolver, sysAPIClient, directorClient, systemfetcher.Config{
				SystemsQueueSize:     100,
				FetcherParallellism:  30,
				EnableSystemDeletion: true,
				VerifyTenant:         testCase.verificationTenant,
			})

			err := svc.SyncSystems(context.TODO())
			if testCase.expectedErr != nil {
				require.ErrorIs(t, err, testCase.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpsertSystemsSyncTimestamps(t *testing.T) {
	testError := errors.New("testError")

	systemfetcher.SystemSynchronizationTimestamps = map[string]map[string]systemfetcher.SystemSynchronizationTimestamp{
		"t": {
			"type1": {
				ID:                "time",
				LastSyncTimestamp: time.Date(2023, 5, 2, 20, 30, 0, 0, time.UTC).UTC(),
			},
		},
	}

	type testCase struct {
		name                     string
		mockTransactioner        func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner)
		fixTestSystems           func() []systemfetcher.System
		fixAppInputs             func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate
		setupTenantSvc           func() *automock.TenantService
		setupTbtSvc              func() *automock.TenantBusinessTypeService
		setupTemplateRendererSvc func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInput) *automock.TemplateRenderer
		setupSystemSvc           func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService
		setupSystemsSyncSvc      func() *automock.SystemsSyncService
		setupSysAPIClient        func(testSystems []systemfetcher.System) *automock.SystemsAPIClient
		setupDirectorClient      func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient
		verificationTenant       string
		expectedErr              error
	}

	tests := []testCase{
		{
			name: "Success",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatSucceedsMultipleTimes(1)
				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				return []systemfetcher.System{}
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return []model.ApplicationRegisterInputWithTemplate{}
			},
			setupTenantSvc: func() *automock.TenantService {
				return &automock.TenantService{}
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				return &automock.TenantBusinessTypeService{}
			},
			setupTemplateRendererSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				return &automock.SystemsService{}
			},
			setupSystemsSyncSvc: func() *automock.SystemsSyncService {
				syncMock := &automock.SystemsSyncService{}
				syncMock.On("Upsert", mock.Anything, mock.Anything).Return(nil)
				return syncMock
			},
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				return &automock.SystemsAPIClient{}
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
		},
		{
			name: "Error while upserting",
			mockTransactioner: func() (*pAutomock.PersistenceTx, *pAutomock.Transactioner) {
				mockedTx, transactioner := txtest.NewTransactionContextGenerator(nil).ThatDoesntExpectCommit()
				return mockedTx, transactioner
			},
			fixTestSystems: func() []systemfetcher.System {
				return []systemfetcher.System{}
			},
			fixAppInputs: func(systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
				return []model.ApplicationRegisterInputWithTemplate{}
			},
			setupTenantSvc: func() *automock.TenantService {
				return &automock.TenantService{}
			},
			setupTbtSvc: func() *automock.TenantBusinessTypeService {
				return &automock.TenantBusinessTypeService{}
			},
			setupTemplateRendererSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInput) *automock.TemplateRenderer {
				return &automock.TemplateRenderer{}
			},
			setupSystemSvc: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.SystemsService {
				return &automock.SystemsService{}
			},
			setupSystemsSyncSvc: func() *automock.SystemsSyncService {
				syncMock := &automock.SystemsSyncService{}
				syncMock.On("Upsert", mock.Anything, mock.Anything).Return(testError)
				return syncMock
			},
			setupSysAPIClient: func(testSystems []systemfetcher.System) *automock.SystemsAPIClient {
				return &automock.SystemsAPIClient{}
			},
			setupDirectorClient: func(systems []systemfetcher.System, appsInputs []model.ApplicationRegisterInputWithTemplate) *automock.DirectorClient {
				return &automock.DirectorClient{}
			},
			expectedErr: testError,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mockedTx, transactioner := testCase.mockTransactioner()
			tenantSvc := testCase.setupTenantSvc()
			tbtSvc := testCase.setupTbtSvc()
			testSystems := testCase.fixTestSystems()
			appsInputs := testCase.fixAppInputs(testSystems)
			systemSvc := testCase.setupSystemSvc(testSystems, appsInputs)
			systemsSyncSvc := testCase.setupSystemsSyncSvc()
			appInputsWithoutTemplates := make([]model.ApplicationRegisterInput, 0)
			for _, in := range appsInputs {
				appInputsWithoutTemplates = append(appInputsWithoutTemplates, in.ApplicationRegisterInput)
			}
			templateAppResolver := testCase.setupTemplateRendererSvc(testSystems, appInputsWithoutTemplates)
			sysAPIClient := testCase.setupSysAPIClient(testSystems)
			directorClient := testCase.setupDirectorClient(testSystems, appsInputs)
			defer mock.AssertExpectationsForObjects(t, tenantSvc, sysAPIClient, systemSvc, templateAppResolver, mockedTx, transactioner)

			svc := systemfetcher.NewSystemFetcher(transactioner, tenantSvc, systemSvc, systemsSyncSvc, tbtSvc, templateAppResolver, sysAPIClient, directorClient, systemfetcher.Config{
				SystemsQueueSize:     100,
				FetcherParallellism:  30,
				EnableSystemDeletion: true,
				VerifyTenant:         testCase.verificationTenant,
			})

			err := svc.UpsertSystemsSyncTimestamps(context.TODO(), transactioner)
			if testCase.expectedErr != nil {
				require.ErrorIs(t, err, testCase.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}

	systemfetcher.SystemSynchronizationTimestamps = nil
}

func fixAppsInputsWithTemplatesBySystems(t *testing.T, systems []systemfetcher.System) []model.ApplicationRegisterInputWithTemplate {
	initStatusCond := model.ApplicationStatusConditionInitial
	result := make([]model.ApplicationRegisterInputWithTemplate, 0, len(systems))
	for _, s := range systems {
		systemPayload, err := json.Marshal(s.SystemPayload)
		require.NoError(t, err)
		input := model.ApplicationRegisterInputWithTemplate{
			ApplicationRegisterInput: model.ApplicationRegisterInput{
				Name:            gjson.GetBytes(systemPayload, "displayName").String(),
				Description:     str.Ptr(gjson.GetBytes(systemPayload, "productDescription").String()),
				BaseURL:         str.Ptr(gjson.GetBytes(systemPayload, "additionalUrls"+"."+mainURLKey).String()),
				ProviderName:    str.Ptr(gjson.GetBytes(systemPayload, "infrastructureProvider").String()),
				SystemNumber:    str.Ptr(gjson.GetBytes(systemPayload, "systemNumber").String()),
				StatusCondition: &initStatusCond,
				Labels: map[string]interface{}{
					"managed":              "true",
					"productId":            str.Ptr(gjson.GetBytes(systemPayload, "productId").String()),
					"ppmsProductVersionId": str.Ptr(gjson.GetBytes(systemPayload, "ppmsProductVersionId").String()),
				},
			},
			TemplateID: s.TemplateID,
		}
		if len(input.TemplateID) > 0 {
			input.Labels["applicationType"] = appType
		}
		result = append(result, input)
	}
	return result
}

func emptySystemsSyncSvc() *automock.SystemsSyncService {
	return &automock.SystemsSyncService{}
}
