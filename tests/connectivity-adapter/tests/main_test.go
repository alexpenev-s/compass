package tests

import (
	"context"
	"os"
	"testing"

	"github.com/machinebox/graphql"

	"github.com/kyma-incubator/compass/tests/pkg/clients"

	"github.com/kyma-incubator/compass/components/director/pkg/certloader"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/tests/pkg/gql"
	"github.com/kyma-incubator/compass/tests/pkg/util"
	"github.com/pkg/errors"
	"github.com/vrischmann/envconfig"
)

var (
	directorClient clients.Client
)

type config struct {
	ConnectivityAdapterMtlsUrl     string `envconfig:"default=https://adapter-gateway-mtls.kyma.local"`
	DirectorExternalCertSecuredURL string `envconfig:"default=http://compass-director-external-mtls.compass-system.svc.cluster.local:3000/graphql"`
	SkipSSLValidation              bool   `envconfig:"default=true"`
	EventsBaseURL                  string `envconfig:"default=https://events.com"`
	Tenant                         string `envconfig:"default=3e64ebae-38b5-46a0-b1ed-9ccee153a0ae"`
	DirectorReadyzUrl              string `envconfig:"default=http://compass-director.compass-system.svc.cluster.local:3000/readyz"`
	ApplicationTypeLabelKey        string `envconfig:"APP_APPLICATION_TYPE_LABEL_KEY,default=applicationType"`
	GatewayOauth                   string `envconfig:"APP_GATEWAY_OAUTH"`
	CertLoaderConfig               certloader.Config

	ExternalClientCertSecretName string `envconfig:"APP_EXTERNAL_CLIENT_CERT_SECRET_NAME"`
}

var (
	testConfig               config
	certSecuredGraphQLClient *graphql.Client
)

func TestMain(m *testing.M) {
	err := envconfig.InitWithPrefix(&testConfig, "APP")
	if err != nil {
		log.D().Fatal(errors.Wrap(err, "while initializing envconfig"))
	}

	ctx := context.Background()
	cc, err := certloader.StartCertLoader(ctx, testConfig.CertLoaderConfig)
	if err != nil {
		log.D().Fatal(errors.Wrap(err, "while starting cert cache"))
	}

	if err := util.WaitForCache(cc); err != nil {
		log.D().Fatal(err)
	}

	certSecuredGraphQLClient = gql.NewCertAuthorizedGraphQLClientWithCustomURL(testConfig.DirectorExternalCertSecuredURL, cc.Get()[testConfig.ExternalClientCertSecretName].PrivateKey, cc.Get()[testConfig.ExternalClientCertSecretName].Certificate, testConfig.SkipSSLValidation)
	directorClient, err = clients.NewDirectorClient(certSecuredGraphQLClient, testConfig.Tenant, testConfig.DirectorReadyzUrl)
	if err != nil {
		log.D().Fatal(err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}
