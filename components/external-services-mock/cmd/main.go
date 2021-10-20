package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/kyma-incubator/compass/components/external-services-mock/pkg/health"

	"github.com/kyma-incubator/compass/components/external-services-mock/internal/cert"

	"github.com/kyma-incubator/compass/components/external-services-mock/internal/oauth"

	"github.com/kyma-incubator/compass/components/external-services-mock/pkg/webhook"

	"github.com/kyma-incubator/compass/components/external-services-mock/internal/apispec"
	ord_aggregator "github.com/kyma-incubator/compass/components/external-services-mock/internal/ord-aggregator"
	"github.com/kyma-incubator/compass/components/external-services-mock/internal/systemfetcher"

	"github.com/kyma-incubator/compass/components/external-services-mock/internal/httphelpers"

	"github.com/kyma-incubator/compass/components/external-services-mock/internal/auditlog/configurationchange"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/vrischmann/envconfig"
)

type config struct {
	Port       int `envconfig:"default=8080"`
	ORDServers ORDServers
	BaseURL    string `envconfig:"default=http://compass-external-services-mock.compass-system.svc.cluster.local"`
	JWKSPath   string `envconfig:"default=/jwks.json"`
	OAuthConfig
	BasicCredentialsConfig
	DefaultTenant string `envconfig:"APP_DEFAULT_TENANT"`

	CACert string `envconfig:"APP_CA_CERT"`
	CAKey  string `envconfig:"APP_CA_KEY"`
}

// ORDServers is a configuration for ORD e2e tests. Those tests are more complex and require a dedicated server per application involved.
// This is needed in order to ensure that every call in the context of an application happens in a single server isolated from others.
// Prior to this separation there were cases when tests succeeded (false positive) due to mistakenly configured baseURL resulting in different flow - different access strategy returned.
type ORDServers struct {
	CertPort      int `envconfig:"default=8081"`
	UnsecuredPort int `envconfig:"default=8082"`
	BasicPort     int `envconfig:"default=8083"`
	OauthPort     int `envconfig:"default=8084"`

	CertSecuredBaseURL string
}

type OAuthConfig struct {
	ClientID     string `envconfig:"APP_CLIENT_ID"`
	ClientSecret string `envconfig:"APP_CLIENT_SECRET"`
}

type BasicCredentialsConfig struct {
	Username string `envconfig:"BASIC_USERNAME"`
	Password string `envconfig:"BASIC_PASSWORD"`
}

func main() {
	ctx := context.Background()

	cfg := config{}
	err := envconfig.InitWithOptions(&cfg, envconfig.Options{Prefix: "APP", AllOptional: true})
	exitOnError(err, "while loading configuration")

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	exitOnError(err, "while generating rsa key")

	ordServers := initORDServers(cfg, key)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go startServer(ctx, initDefaultServer(cfg, key), wg)

	for _, server := range ordServers {
		wg.Add(1)
		go startServer(ctx, server, wg)
	}

	wg.Wait()
}

func exitOnError(err error, context string) {
	if err != nil {
		wrappedError := errors.Wrap(err, context)
		log.Fatal(wrappedError)
	}
}

func initDefaultServer(cfg config, key *rsa.PrivateKey) *http.Server {
	logger := logrus.New()
	router := mux.NewRouter()

	router.HandleFunc("/v1/healtz", health.HandleFunc)

	// Oauth server handlers
	tokenHandler := oauth.NewHandlerWithSigningKey(cfg.ClientSecret, cfg.ClientID, key)
	router.HandleFunc("/secured/oauth/token", tokenHandler.Generate).Methods(http.MethodPost)
	openIDConfigHandler := oauth.NewOpenIDConfigHandler(fmt.Sprintf("%s:%d", cfg.BaseURL, cfg.Port), cfg.JWKSPath)
	router.HandleFunc("/.well-known/openid-configuration", openIDConfigHandler.Handle)
	jwksHanlder := oauth.NewJWKSHandler(&key.PublicKey)
	router.HandleFunc(cfg.JWKSPath, jwksHanlder.Handle)

	// CA server handlers
	certHandler := cert.NewHandler(cfg.CACert, cfg.CAKey)
	router.HandleFunc("/cert", certHandler.Generate).Methods(http.MethodPost)

	// AL handlers
	configChangeSvc := configurationchange.NewService()
	configChangeHandler := configurationchange.NewConfigurationHandler(configChangeSvc, logger)
	configChangeRouter := router.PathPrefix("/audit-log/v2/configuration-changes").Subrouter()
	configChangeRouter.Use(oauthMiddleware(&key.PublicKey))
	configurationchange.InitConfigurationChangeHandler(configChangeRouter, configChangeHandler)

	// System fetcher handlers
	systemFetcherHandler := systemfetcher.NewSystemFetcherHandler(cfg.DefaultTenant)
	router.Methods(http.MethodPost).PathPrefix("/systemfetcher/configure").HandlerFunc(systemFetcherHandler.HandleConfigure)
	router.Methods(http.MethodDelete).PathPrefix("/systemfetcher/reset").HandlerFunc(systemFetcherHandler.HandleReset)
	router.HandleFunc("/systemfetcher/systems", systemFetcherHandler.HandleFunc)

	// Fetch request handlers
	router.HandleFunc("/external-api/spec", apispec.HandleFunc)

	oauthRouter := router.PathPrefix("/external-api/secured/oauth").Subrouter()
	oauthRouter.Use(oauthMiddleware(&key.PublicKey))
	oauthRouter.HandleFunc("/spec", apispec.HandleFunc)

	basicAuthRouter := router.PathPrefix("/external-api/secured/basic").Subrouter()
	basicAuthRouter.Use(basicAuthMiddleware(cfg.Username, cfg.Password))
	basicAuthRouter.HandleFunc("/spec", apispec.HandleFunc)

	// Operations controller handlers
	router.HandleFunc(webhook.DeletePath, webhook.NewDeleteHTTPHandler()).Methods(http.MethodDelete)
	router.HandleFunc(webhook.OperationPath, webhook.NewWebHookOperationGetHTTPHandler()).Methods(http.MethodGet)
	router.HandleFunc(webhook.OperationPath, webhook.NewWebHookOperationPostHTTPHandler()).Methods(http.MethodPost)

	// non-isolated and unsecured ORD handlers. NOTE: Do not host document endpoints on this default server in order to ensure testс separation.
	// Unsecured config pointing to cert secured document
	router.HandleFunc("/cert/.well-known/open-resource-discovery", ord_aggregator.HandleFuncOrdConfig(cfg.ORDServers.CertSecuredBaseURL, "sap:cmp-mtls:v1"))

	return &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", cfg.Port),
		Handler: router,
	}
}

func initORDServers(cfg config, key *rsa.PrivateKey) []*http.Server {
	servers := make([]*http.Server, 0, 0)
	servers = append(servers, initCertSecuredORDServer(cfg))
	servers = append(servers, initUnsecuredORDServer(cfg))
	servers = append(servers, initBasicSecuredORDServer(cfg))
	servers = append(servers, initOauthSecuredORDServer(cfg, key))
	return servers
}

func initCertSecuredORDServer(cfg config) *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/.well-known/open-resource-discovery", ord_aggregator.HandleFuncOrdConfig("", "sap:cmp-mtls:v1"))

	router.HandleFunc("/open-resource-discovery/v1/documents/example1", ord_aggregator.HandleFuncOrdDocument(cfg.ORDServers.CertSecuredBaseURL, "sap:cmp-mtls:v1"))

	router.HandleFunc("/external-api/spec", apispec.HandleFunc)
	router.HandleFunc("/external-api/spec/flapping", apispec.FlappingHandleFunc())

	return &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", cfg.ORDServers.CertPort),
		Handler: router,
	}
}

func initUnsecuredORDServer(cfg config) *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/.well-known/open-resource-discovery", ord_aggregator.HandleFuncOrdConfig("", "open"))
	router.HandleFunc("/test/fullPath", ord_aggregator.HandleFuncOrdConfig("", "open"))

	router.HandleFunc("/open-resource-discovery/v1/documents/example1", ord_aggregator.HandleFuncOrdDocument(fmt.Sprintf("%s:%d", cfg.BaseURL, cfg.ORDServers.UnsecuredPort), "open"))

	router.HandleFunc("/external-api/spec", apispec.HandleFunc)
	router.HandleFunc("/external-api/spec/flapping", apispec.FlappingHandleFunc())

	return &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", cfg.ORDServers.UnsecuredPort),
		Handler: router,
	}
}

func initBasicSecuredORDServer(cfg config) *http.Server {
	router := mux.NewRouter()

	configRouter := router.PathPrefix("/.well-known").Subrouter()
	configRouter.Use(basicAuthMiddleware(cfg.Username, cfg.Password))
	configRouter.HandleFunc("/open-resource-discovery", ord_aggregator.HandleFuncOrdConfig("", "open"))

	router.HandleFunc("/open-resource-discovery/v1/documents/example1", ord_aggregator.HandleFuncOrdDocument(fmt.Sprintf("%s:%d", cfg.BaseURL, cfg.ORDServers.BasicPort), "open"))

	router.HandleFunc("/external-api/spec", apispec.HandleFunc)
	router.HandleFunc("/external-api/spec/flapping", apispec.FlappingHandleFunc())

	return &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", cfg.ORDServers.BasicPort),
		Handler: router,
	}
}

func initOauthSecuredORDServer(cfg config, key *rsa.PrivateKey) *http.Server {
	router := mux.NewRouter()

	configRouter := router.PathPrefix("/.well-known").Subrouter()
	configRouter.Use(oauthMiddleware(&key.PublicKey))
	configRouter.HandleFunc("/open-resource-discovery", ord_aggregator.HandleFuncOrdConfig("", "open"))

	router.HandleFunc("/open-resource-discovery/v1/documents/example1", ord_aggregator.HandleFuncOrdDocument(fmt.Sprintf("%s:%d", cfg.BaseURL, cfg.ORDServers.OauthPort), "open"))

	router.HandleFunc("/external-api/spec", apispec.HandleFunc)
	router.HandleFunc("/external-api/spec/flapping", apispec.FlappingHandleFunc())

	return &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", cfg.ORDServers.OauthPort),
		Handler: router,
	}
}

func oauthMiddleware(key *rsa.PublicKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) == 0 {
				httphelpers.WriteError(w, errors.New("No Authorization header"), http.StatusUnauthorized)
				return
			}
			if !strings.Contains(authHeader, "Bearer") {
				httphelpers.WriteError(w, errors.New("No Bearer token"), http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if _, err := jwt.Parse(token, func(_ *jwt.Token) (interface{}, error) {
				return key, nil
			}); err != nil {
				httphelpers.WriteError(w, errors.Wrap(err, "Invalid Bearer token"), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func basicAuthMiddleware(username, password string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()

			if !ok {
				httphelpers.WriteError(w, errors.New("No Basic credentials"), http.StatusUnauthorized)
				return
			}
			if username != u || password != p {
				httphelpers.WriteError(w, errors.New("Bad credentials"), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func startServer(parentCtx context.Context, server *http.Server, wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	go func() {
		defer wg.Done()
		<-ctx.Done()
		stopServer(server)
	}()

	log.Printf("Starting and listening on %s://%s", "http", server.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s://%s: %v\n", "http", server.Addr, err)
	}
}

func stopServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	go func(ctx context.Context) {
		<-ctx.Done()

		if ctx.Err() == context.Canceled {
			return
		} else if ctx.Err() == context.DeadlineExceeded {
			log.Fatal("Timeout while stopping the server, killing instance!")
		}
	}(ctx)

	server.SetKeepAlivesEnabled(false)

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
}
