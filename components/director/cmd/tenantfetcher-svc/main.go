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

package main

import (
	"context"
	"crypto/tls"
	"github.com/kyma-incubator/compass/components/director/internal/domain/runtime"
	"github.com/kyma-incubator/compass/components/director/internal/domain/scenarioassignment"
	"github.com/kyma-incubator/compass/components/director/internal/metrics"
	"net/http"
	"os"
	"time"

	"github.com/kyma-incubator/compass/components/director/internal/domain/label"
	"github.com/kyma-incubator/compass/components/director/internal/domain/labeldef"
	"github.com/kyma-incubator/compass/components/director/internal/uid"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"

	graphqlclient "github.com/kyma-incubator/compass/components/director/pkg/graphql_client"
	gcli "github.com/machinebox/graphql"

	httputil "github.com/kyma-incubator/compass/components/director/pkg/http"

	auth "github.com/kyma-incubator/compass/components/director/internal/authenticator"
	"github.com/kyma-incubator/compass/components/director/internal/authenticator/claims"
	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/pkg/executor"

	"github.com/kyma-incubator/compass/components/director/pkg/correlation"

	tf "github.com/kyma-incubator/compass/components/director/internal/tenantfetcher"
	tenantfetcher "github.com/kyma-incubator/compass/components/director/internal/tenantfetchersvc"

	"github.com/gorilla/mux"
	timeouthandler "github.com/kyma-incubator/compass/components/director/pkg/handler"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/signal"
	"github.com/pkg/errors"
	"github.com/vrischmann/envconfig"
)

const (
	envPrefix       = "APP"
	syncSubaccounts = true
)

type config struct {
	Address string `envconfig:"default=127.0.0.1:8080"`

	ServerTimeout   time.Duration `envconfig:"default=110s"`
	ShutdownTimeout time.Duration `envconfig:"default=10s"`

	Log log.Config

	TenantsRootAPI string `envconfig:"APP_ROOT_API,default=/tenants"`

	Handler   tenantfetcher.HandlerConfig
	EventsCfg tenantfetcher.EventsConfig

	SecurityConfig securityConfig
}

type securityConfig struct {
	JWKSSyncPeriod            time.Duration `envconfig:"default=5m"`
	AllowJWTSigningNone       bool          `envconfig:"APP_ALLOW_JWT_SIGNING_NONE,default=false"`
	JwksEndpoint              string        `envconfig:"default=file://hack/default-jwks.json,APP_JWKS_ENDPOINT"`
	SubscriptionCallbackScope string        `envconfig:"APP_SUBSCRIPTION_CALLBACK_SCOPE"`
	FetchTenantOnDemandScope  string        `envconfig:"APP_FETCH_TENANT_ON_DEMAND_SCOPE"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	term := make(chan os.Signal)
	signal.HandleInterrupts(ctx, cancel, term)

	cfg := config{}
	err := envconfig.InitWithPrefix(&cfg, envPrefix)
	exitOnError(err, "Error while loading app config")

	ctx, err = log.Configure(ctx, &cfg.Log)
	exitOnError(err, "Failed to configure Logger")

	if cfg.Handler.TenantPathParam == "" {
		exitOnError(errors.New("missing tenant path parameter"), "Error while loading app handler config")
	}

	transact, closeFunc, err := persistence.Configure(ctx, cfg.Handler.Database)
	exitOnError(err, "Error while establishing the connection to the database")

	defer func() {
		err := closeFunc()
		exitOnError(err, "Error while closing the connection to the database")
	}()

	httpClient := &http.Client{
		Transport: httputil.NewCorrelationIDTransport(http.DefaultTransport),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	handler := initAPIHandler(ctx, httpClient, cfg, transact)
	runMainSrv, shutdownMainSrv := createServer(ctx, cfg, handler, "main")

	go func() {
		<-ctx.Done()
		// Interrupt signal received - shut down the servers
		shutdownMainSrv()
	}()

	runMainSrv()

	stopSubaccountsJob := make(chan bool, 1)
	runTenantFetcherJob(ctx, cfg, transact, syncSubaccounts, stopSubaccountsJob)

	stopGlobalAccountsJob := make(chan bool, 1)
	runTenantFetcherJob(ctx, cfg, transact, !syncSubaccounts, stopGlobalAccountsJob)

	<-stopSubaccountsJob
	<-stopGlobalAccountsJob
}

func runTenantFetcherJob(ctx context.Context, cfg config, transact persistence.Transactioner, syncSubaccounts bool, stopJob chan bool) {
	jobInterval := cfg.Handler.TenantFetcherJobIntervalMins
	ticker := time.NewTicker(time.Duration(jobInterval) * time.Minute)
	jobType := jobType(syncSubaccounts)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.C(ctx).Infof("Scheduled %s tenant fetcher job will be executed, job interval is %d minutes", jobType, jobInterval)
				syncTenants(ctx, cfg, transact, syncSubaccounts)
			case <-ctx.Done():
				log.C(ctx).Errorf("Context is canceled and scheduled %s tenant fetcher job will be stopped", jobType)
				stopTenantFetcherJobTicker(ctx, ticker, jobType)
				stopJob <- true
				return
			}
		}
	}()
}

func jobType(syncSubaccounts bool) string {
	jobType := "subaccount"
	if !syncSubaccounts {
		jobType = "global account"
	}
	return jobType
}

func syncTenants(ctx context.Context, cfg config, transact persistence.Transactioner, syncSubaccounts bool) {
	tenantsFetcherSvc, err := createTenantsFetcherSvc(ctx, cfg, transact, syncSubaccounts)
	exitOnError(err, "failed to create tenants fetcher service")

	tenantsFetcherSvc.SyncTenants()
}

func createTenantsFetcherSvc(ctx context.Context, cfg config, transact persistence.Transactioner, syncSubaccounts bool) (tf.TenantSyncService, error) {
	eventsCfg := cfg.EventsCfg
	handlerCfg := cfg.Handler

	uidSvc := uid.NewService()

	labelDefConverter := labeldef.NewConverter()
	labelDefRepository := labeldef.NewRepository(labelDefConverter)

	labelConverter := label.NewConverter()
	labelRepository := label.NewRepository(labelConverter)
	labelService := label.NewLabelService(labelRepository, labelDefRepository, uidSvc)

	tenantStorageConv := tenant.NewConverter()
	tenantStorageRepo := tenant.NewRepository(tenantStorageConv)
	tenantStorageSvc := tenant.NewServiceWithLabels(tenantStorageRepo, uidSvc, labelRepository, labelService)

	runtimeConverter := runtime.NewConverter()
	runtimeRepository := runtime.NewRepository(runtimeConverter)

	scenarioAssignConv := scenarioassignment.NewConverter()
	scenarioAssignRepo := scenarioassignment.NewRepository(scenarioAssignConv)
	scenarioAssignEngine := scenarioassignment.NewEngine(labelService, labelRepository, scenarioAssignRepo, runtimeRepository)

	scenarioAssignmentRepo := scenarioassignment.NewRepository(scenarioAssignConv)
	labelDefService := labeldef.NewService(labelDefRepository, labelRepository, scenarioAssignmentRepo, tenantStorageRepo, uidSvc, handlerCfg.Features.DefaultScenarioEnabled)

	runtimeService := runtime.NewService(runtimeRepository, labelRepository, labelDefService, labelService, uidSvc, scenarioAssignEngine, tenantStorageSvc, handlerCfg.Features.ProtectedLabelPattern, handlerCfg.Features.ImmutableLabelPattern)

	var metricsPusher *metrics.Pusher
	if handlerCfg.MetricsPushEndpoint != "" {
		metricsPusher = metrics.NewPusher(handlerCfg.MetricsPushEndpoint, handlerCfg.ClientTimeout)
	}

	kubeClient, err := tf.NewKubernetesClient(ctx, handlerCfg.Kubernetes)
	exitOnError(err, "Failed to initialize Kubernetes client")

	eventAPIClient, err := tf.NewClient(eventsCfg.OAuthConfig, eventsCfg.AuthMode, eventsCfg.APIConfig, handlerCfg.ClientTimeout)
	if nil != err {
		return nil, err
	}

	if metricsPusher != nil {
		eventAPIClient.SetMetricsPusher(metricsPusher)
	}

	gqlClient := newInternalGraphQLClient(handlerCfg.DirectorGraphQLEndpoint, handlerCfg.ClientTimeout, handlerCfg.HTTPClientSkipSslValidation)
	gqlClient.Log = func(s string) {
		log.D().Debug(s)
	}
	directorClient := graphqlclient.NewDirector(gqlClient)

	if syncSubaccounts {
		return tf.NewSubaccountService(eventsCfg.QueryConfig, transact, kubeClient, eventsCfg.TenantFieldMapping, eventsCfg.MovedSubaccountFieldMapping, handlerCfg.TenantProvider, eventsCfg.SubaccountRegions, eventAPIClient, tenantStorageSvc, runtimeService, labelRepository, handlerCfg.FullResyncInterval, directorClient, handlerCfg.TenantInsertChunkSize, tenantStorageConv), nil
	}
	return tf.NewGlobalAccountService(eventsCfg.QueryConfig, transact, kubeClient, eventsCfg.TenantFieldMapping, handlerCfg.TenantProvider, eventsCfg.AccountsRegion, eventAPIClient, tenantStorageSvc, handlerCfg.FullResyncInterval, directorClient, handlerCfg.TenantInsertChunkSize, tenantStorageConv), nil
}

func stopTenantFetcherJobTicker(ctx context.Context, tenantFetcherJobTicker *time.Ticker, jobType string) {
	tenantFetcherJobTicker.Stop()
	log.C(ctx).Infof("Tenant fetcher job ticker for %ss is stopped", jobType)
}

func initAPIHandler(ctx context.Context, httpClient *http.Client, cfg config, transact persistence.Transactioner) http.Handler {
	logger := log.C(ctx)
	mainRouter := mux.NewRouter()
	mainRouter.Use(correlation.AttachCorrelationIDToContext(), log.RequestLogger())

	tenantsAPIRouter := mainRouter.PathPrefix(cfg.TenantsRootAPI).Subrouter()
	configureAuthMiddleware(ctx, httpClient, tenantsAPIRouter, cfg.SecurityConfig, []string{cfg.SecurityConfig.SubscriptionCallbackScope})
	registerTenantsHandler(ctx, tenantsAPIRouter, cfg.Handler)

	tenantsOnDemandAPIRouter := mainRouter.PathPrefix(cfg.TenantsRootAPI).Subrouter()
	configureAuthMiddleware(ctx, httpClient, tenantsOnDemandAPIRouter, cfg.SecurityConfig, []string{cfg.SecurityConfig.FetchTenantOnDemandScope})
	registerTenantsOnDemandHandler(ctx, tenantsOnDemandAPIRouter, cfg.EventsCfg, cfg.Handler, transact)

	healthCheckRouter := mainRouter.PathPrefix(cfg.TenantsRootAPI).Subrouter()
	logger.Infof("Registering readiness endpoint...")
	healthCheckRouter.HandleFunc("/readyz", newReadinessHandler())
	logger.Infof("Registering liveness endpoint...")
	healthCheckRouter.HandleFunc("/healthz", newReadinessHandler())

	return mainRouter
}

func exitOnError(err error, context string) {
	if err != nil {
		wrappedError := errors.Wrap(err, context)
		log.D().Fatal(wrappedError)
	}
}

func createServer(ctx context.Context, cfg config, handler http.Handler, name string) (func(), func()) {
	logger := log.C(ctx)

	handlerWithTimeout, err := timeouthandler.WithTimeout(handler, cfg.ServerTimeout)
	exitOnError(err, "Error while configuring tenant mapping handler")

	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           handlerWithTimeout,
		ReadHeaderTimeout: cfg.ServerTimeout,
	}

	runFn := func() {
		logger.Infof("Running %s server on %s...", name, cfg.Address)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Errorf("%s HTTP server ListenAndServe: %v", name, err)
		}
	}

	shutdownFn := func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		logger.Infof("Shutting down %s server...", name)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf("%s HTTP server Shutdown: %v", name, err)
		}
	}

	return runFn, shutdownFn
}

func configureAuthMiddleware(ctx context.Context, httpClient *http.Client, router *mux.Router, cfg securityConfig, requiredScopes []string) {
	scopeValidator := claims.NewScopesValidator(requiredScopes)
	middleware := auth.New(httpClient, cfg.JwksEndpoint, cfg.AllowJWTSigningNone, "", scopeValidator)
	router.Use(middleware.Handler())

	log.C(ctx).Infof("JWKS synchronization enabled. Sync period: %v", cfg.JWKSSyncPeriod)
	periodicExecutor := executor.NewPeriodic(cfg.JWKSSyncPeriod, func(ctx context.Context) {
		if err := middleware.SynchronizeJWKS(ctx); err != nil {
			log.C(ctx).WithError(err).Errorf("An error has occurred while synchronizing JWKS: %v", err)
		}
	})
	go periodicExecutor.Run(ctx)
}

func registerTenantsHandler(ctx context.Context, router *mux.Router, cfg tenantfetcher.HandlerConfig) {
	gqlClient := newInternalGraphQLClient(cfg.DirectorGraphQLEndpoint, cfg.ClientTimeout, cfg.HTTPClientSkipSslValidation)
	gqlClient.Log = func(s string) {
		log.C(ctx).Debug(s)
	}
	directorClient := graphqlclient.NewDirector(gqlClient)

	tenantConverter := tenant.NewConverter()

	provisioner := tenantfetcher.NewTenantProvisioner(directorClient, tenantConverter, cfg.TenantProvider)
	subscriber := tenantfetcher.NewSubscriber(directorClient, provisioner)
	tenantHandler := tenantfetcher.NewTenantsHTTPHandler(subscriber, cfg)

	log.C(ctx).Infof("Registering Regional Tenant Onboarding endpoint on %s...", cfg.RegionalHandlerEndpoint)
	router.HandleFunc(cfg.RegionalHandlerEndpoint, tenantHandler.SubscribeTenant).Methods(http.MethodPut)

	log.C(ctx).Infof("Registering Regional Tenant Decommissioning endpoint on %s...", cfg.RegionalHandlerEndpoint)
	router.HandleFunc(cfg.RegionalHandlerEndpoint, tenantHandler.UnSubscribeTenant).Methods(http.MethodDelete)

	log.C(ctx).Infof("Registering service dependencies endpoint on %s...", cfg.DependenciesEndpoint)
	router.HandleFunc(cfg.DependenciesEndpoint, tenantHandler.Dependencies).Methods(http.MethodGet)
}

func registerTenantsOnDemandHandler(ctx context.Context, router *mux.Router, eventsCfg tenantfetcher.EventsConfig, tenantHandlerCfg tenantfetcher.HandlerConfig, transact persistence.Transactioner) {
	onDemandSvc, err := createTenantFetcherOnDemandSvc(eventsCfg, tenantHandlerCfg, transact)
	exitOnError(err, "failed to create tenant fetcher on-demand service")

	fetcher := tenantfetcher.NewTenantFetcher(*onDemandSvc)
	tenantHandler := tenantfetcher.NewTenantFetcherHTTPHandler(fetcher, tenantHandlerCfg)

	log.C(ctx).Infof("Registering fetch tenant on-demand endpoint on %s...", tenantHandlerCfg.TenantOnDemandHandlerEndpoint)
	router.HandleFunc(tenantHandlerCfg.TenantOnDemandHandlerEndpoint, tenantHandler.FetchTenantOnDemand).Methods(http.MethodPost)
}

func createTenantFetcherOnDemandSvc(eventsCfg tenantfetcher.EventsConfig, handlerCfg tenantfetcher.HandlerConfig, transact persistence.Transactioner) (*tf.SubaccountOnDemandService, error) {
	eventAPIClient, err := tf.NewClient(eventsCfg.OAuthConfig, eventsCfg.AuthMode, eventsCfg.APIConfig, handlerCfg.ClientTimeout)
	if nil != err {
		return nil, err
	}

	tenantStorageConv := tenant.NewConverter()
	uidSvc := uid.NewService()

	labelDefConverter := labeldef.NewConverter()
	labelDefRepository := labeldef.NewRepository(labelDefConverter)

	labelConverter := label.NewConverter()
	labelRepository := label.NewRepository(labelConverter)
	labelService := label.NewLabelService(labelRepository, labelDefRepository, uidSvc)

	tenantStorageRepo := tenant.NewRepository(tenantStorageConv)
	tenantStorageSvc := tenant.NewServiceWithLabels(tenantStorageRepo, uidSvc, labelRepository, labelService)

	gqlClient := newInternalGraphQLClient(handlerCfg.DirectorGraphQLEndpoint, handlerCfg.ClientTimeout, handlerCfg.HTTPClientSkipSslValidation)
	gqlClient.Log = func(s string) {
		log.D().Debug(s)
	}
	directorClient := graphqlclient.NewDirector(gqlClient)

	return tf.NewSubaccountOnDemandService(eventsCfg.QueryConfig, eventsCfg.TenantFieldMapping, eventAPIClient, transact, tenantStorageSvc, directorClient, handlerCfg.TenantProvider, tenantStorageConv), nil
}

func newReadinessHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}
}

func newInternalGraphQLClient(url string, timeout time.Duration, skipSSLValidation bool) *gcli.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipSSLValidation,
		},
	}
	client := &http.Client{
		Transport: httputil.NewCorrelationIDTransport(httputil.NewServiceAccountTokenTransportWithHeader(tr, "Authorization")),
		Timeout:   timeout,
	}
	return gcli.NewClient(url, gcli.WithHTTPClient(client))
}
