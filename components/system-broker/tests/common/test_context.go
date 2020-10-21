package common

import (
	"context"
	"fmt"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/graphqlizer"
	"github.com/kyma-incubator/compass/components/system-broker/internal/director"
	"github.com/kyma-incubator/compass/components/system-broker/internal/osb"
	"github.com/kyma-incubator/compass/components/system-broker/internal/specs"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/graphql"
	gql "github.com/machinebox/graphql"
	"log"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	sblog "github.com/kyma-incubator/compass/components/system-broker/pkg/log"
	"github.com/kyma-incubator/compass/components/system-broker/internal/config"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/env"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/server"
	"github.com/kyma-incubator/compass/components/system-broker/pkg/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
)

const SystemBrokerServer = "system-broker-server"
const DirectorServer = "director-server"

type TestContext struct {
	SystemBroker *httpexpect.Expect

	Servers map[string]FakeServer
}

func (tc *TestContext) CleanUp() {
	for _, server := range tc.Servers {
		server.Close()
	}
}

type TestContextBuilder struct {
	envHooks []func(env env.Environment, servers map[string]FakeServer)

	Environment env.Environment

	Servers    map[string]FakeServer
	HttpClient *http.Client
}

// NewTestContextBuilder sets up a builder with default values
func NewTestContextBuilder() *TestContextBuilder {
	return &TestContextBuilder{
		Environment: TestEnv(),
		envHooks: []func(env env.Environment, servers map[string]FakeServer){
			func(env env.Environment, servers map[string]FakeServer) {
				env.Set("server.shutdown_timeout", "1s")
				port := findFreePort()
				env.Set("server.port", port)
				env.Set("server.self_url", "http://localhost:"+port)
			},
			func(env env.Environment, servers map[string]FakeServer) {
				env.Set("graphql_client.graphql_endpoint", servers[DirectorServer].URL() + "/graphql")
			},
		},
		Servers: map[string]FakeServer{},
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   20 * time.Second,
					KeepAlive: 20 * time.Second,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       30 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}
}

func (tcb *TestContextBuilder) WithDefaultEnv(env env.Environment) *TestContextBuilder {
	tcb.Environment = env

	return tcb
}

func (tcb *TestContextBuilder) WithEnvExtensions(fs ...func(e env.Environment, servers map[string]FakeServer)) *TestContextBuilder {
	tcb.envHooks = append(tcb.envHooks, fs...)

	return tcb
}

func (tcb *TestContextBuilder) WithHttpClient(client *http.Client) *TestContextBuilder {
	tcb.HttpClient = client

	return tcb
}

func (tcb *TestContextBuilder) Build(t *testing.T) *TestContext {
	gqlMockHandler, err := NewGqlFakeRouter("director", "../../../director/pkg/graphql/schema.graphql")
	if err != nil {
		panic(fmt.Errorf("could not build gql mock handler: %s", err))
	}
	gqlMockServer := NewGqlFakeServer(gqlMockHandler.Handler())
	tcb.Servers[DirectorServer] = gqlMockServer

	for _, envPostHook := range tcb.envHooks {
		envPostHook(tcb.Environment, tcb.Servers)
	}

	sbServer := newSystemBrokerServer(tcb.Environment)
	tcb.Servers[SystemBrokerServer] = sbServer

	systemBroker := httpexpect.New(t, sbServer.URL()).Builder(func(request *httpexpect.Request) {
		request.WithClient(tcb.HttpClient)
	})

	testContext := &TestContext{
		SystemBroker: systemBroker,
		Servers:      tcb.Servers,
	}

	return testContext
}

func TestEnv() env.Environment {
	env, err := env.Default(context.TODO())
	if err != nil {
		panic(err)
	}
	return env
}

type testSystemBrokerServer struct {
	url             string
	cancel          context.CancelFunc
	shutdownTimeout time.Duration
}

func (ts *testSystemBrokerServer) URL() string {
	return ts.url
}

func (ts *testSystemBrokerServer) Close() {
	ts.cancel()
	time.Sleep(ts.shutdownTimeout)
}

func newSystemBrokerServer(sbEnv env.Environment) FakeServer {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.New(sbEnv)
	if err != nil {
		panic(err)
	}

	directorGraphQLClient, err := prepareGQLClient(cfg)
	systemBroker := osb.NewSystemBroker(directorGraphQLClient, cfg.Server.SelfURL+cfg.Server.RootAPI)
	osbApi := osb.API(cfg.Server.RootAPI, systemBroker, sblog.NewDefaultLagerAdapter())
	specsApi := specs.API(cfg.Server.RootAPI, directorGraphQLClient)
	sbServer := server.New(cfg.Server, uuid.NewService(), osbApi, specsApi)

	sbServer.Addr = "localhost:" + strconv.Itoa(cfg.Server.Port) // Needed to avoid annoying macOS permissions popup

	go sbServer.Start(ctx)

	err = wait.PollImmediate(time.Millisecond*250, time.Second*5, func() (bool, error) {
		_, err := http.Get(fmt.Sprintf("http://%s", sbServer.Addr))
		if err != nil {
			log.Printf("Waiting for server to start: %v", err)
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		panic(err)
	}

	return &testSystemBrokerServer{
		url:             cfg.Server.SelfURL + cfg.Server.RootAPI,
		cancel:          cancel,
		shutdownTimeout: cfg.Server.ShutdownTimeout,
	}
}

func prepareGQLClient(cfg *config.Config) (*director.GraphQLClient, error) {
	graphClient := gql.NewClient(cfg.GraphQLClient.GraphqlEndpoint, gql.WithHTTPClient(http.DefaultClient))
	gqlClient := graphql.NewClient(cfg.GraphQLClient, graphClient)

	inputGraphqlizer := &graphqlizer.Graphqlizer{}
	outputGraphqlizer := &graphqlizer.GqlFieldsProvider{}

	// prepare director graphql client
	return director.NewGraphQLClient(gqlClient, inputGraphqlizer, outputGraphqlizer), nil
}


func findFreePort() string {
	// Create a new listener without specifying a port which will result in an open port being chosen
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	hostString := listener.Addr().String()
	_, port, err := net.SplitHostPort(hostString)
	if err != nil {
		panic(err)
	}

	return port
}
