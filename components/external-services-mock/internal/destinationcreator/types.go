package destinationcreator

import (
	"encoding/json"
	"regexp"

	destinationcreatorpkg "github.com/kyma-incubator/compass/components/director/pkg/destinationcreator"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// Config holds destination creator service API test configuration
type Config struct {
	CorrelationIDsKey string `envconfig:"APP_DESTINATION_CREATOR_CORRELATION_IDS_KEY"`
	*DestinationAPIConfig
	*CertificateAPIConfig
}

// DestinationAPIConfig holds a test configuration specific for the destination API of the destination creator service
type DestinationAPIConfig struct {
	Path                 string `envconfig:"APP_DESTINATION_CREATOR_DESTINATION_PATH"`
	RegionParam          string `envconfig:"APP_DESTINATION_CREATOR_DESTINATION_REGION_PARAMETER"`
	SubaccountIDParam    string `envconfig:"APP_DESTINATION_CREATOR_DESTINATION_SUBACCOUNT_ID_PARAMETER"`
	DestinationNameParam string `envconfig:"APP_DESTINATION_CREATOR_DESTINATION_NAME_PARAMETER"`
}

// CertificateAPIConfig holds a test configuration specific for the certificate API of the destination creator service
type CertificateAPIConfig struct {
	Path                 string `envconfig:"APP_DESTINATION_CREATOR_CERTIFICATE_PATH"`
	RegionParam          string `envconfig:"APP_DESTINATION_CREATOR_CERTIFICATE_REGION_PARAMETER"`
	SubaccountIDParam    string `envconfig:"APP_DESTINATION_CREATOR_CERTIFICATE_SUBACCOUNT_ID_PARAMETER"`
	CertificateNameParam string `envconfig:"APP_DESTINATION_CREATOR_CERTIFICATE_NAME_PARAMETER"`
}

// BaseDestinationRequestBody contains the base fields needed in the destination request body
type BaseDestinationRequestBody struct {
	Name                 string                          `json:"name"`
	URL                  string                          `json:"url"`
	Type                 destinationcreatorpkg.Type      `json:"type"`
	ProxyType            destinationcreatorpkg.ProxyType `json:"proxyType"`
	AuthenticationType   destinationcreatorpkg.AuthType  `json:"authenticationType"`
	AdditionalProperties json.RawMessage                 `json:"additionalProperties,omitempty"`
}

// DesignTimeDestRequestBody contains the necessary fields for the destination request body with authentication type AuthTypeNoAuth
type DesignTimeDestRequestBody struct {
	BaseDestinationRequestBody
}

// BasicDestRequestBody contains the necessary fields for the destination request body with authentication type AuthTypeBasic
type BasicDestRequestBody struct {
	BaseDestinationRequestBody
	User     string `json:"user"`
	Password string `json:"password"`
}

// SAMLAssertionDestRequestBody contains the necessary fields for the destination request body with authentication type AuthTypeSAMLAssertion
type SAMLAssertionDestRequestBody struct {
	BaseDestinationRequestBody
	Audience         string `json:"audience"`
	KeyStoreLocation string `json:"keyStoreLocation"`
}

// ClientCertificateAuthDestRequestBody contains the necessary fields for the destination request body with authentication type AuthTypeClientCertificate
type ClientCertificateAuthDestRequestBody struct {
	BaseDestinationRequestBody
	KeyStoreLocation string `json:"keyStoreLocation"`
}

// CertificateRequestBody contains the necessary fields for the destination creator certificate request body
type CertificateRequestBody struct {
	Name string `json:"name"`
}

// CertificateResponseBody contains the response data from the destination creator certificate request
type CertificateResponseBody struct {
	FileName         string `json:"fileName"`
	CommonName       string `json:"commonName"`
	CertificateChain string `json:"certificateChain"`
}

// reqBodyNameRegex is a regex defined by the destination creator API specifying what destination names are allowed
var reqBodyNameRegex = "[a-zA-Z0-9_-]{1,64}"

// Validate validates that the AuthTypeNoAuth request body contains the required fields and they are valid
func (n *DesignTimeDestRequestBody) Validate() error {
	return validation.ValidateStruct(n,
		validation.Field(&n.Name, validation.Required, validation.Length(1, destinationcreatorpkg.MaxDestinationNameLength), validation.Match(regexp.MustCompile(reqBodyNameRegex))),
		validation.Field(&n.URL, validation.Required),
		validation.Field(&n.Type, validation.In(destinationcreatorpkg.TypeHTTP, destinationcreatorpkg.TypeRFC, destinationcreatorpkg.TypeLDAP, destinationcreatorpkg.TypeMAIL)),
		validation.Field(&n.ProxyType, validation.In(destinationcreatorpkg.ProxyTypeInternet, destinationcreatorpkg.ProxyTypeOnPremise, destinationcreatorpkg.ProxyTypePrivateLink)),
		validation.Field(&n.AuthenticationType, validation.In(destinationcreatorpkg.AuthTypeNoAuth)),
	)
}

// Validate validates that the AuthTypeBasic request body contains the required fields and they are valid
func (b *BasicDestRequestBody) Validate(destinationCreatorCfg *Config) error {
	areAdditionalPropertiesValid := newDestinationDetailsAdditionalPropertiesValidator(destinationCreatorCfg)

	return validation.ValidateStruct(b,
		validation.Field(&b.Name, validation.Required, validation.Length(1, destinationcreatorpkg.MaxDestinationNameLength), validation.Match(regexp.MustCompile(reqBodyNameRegex))),
		validation.Field(&b.URL, validation.Required),
		validation.Field(&b.Type, validation.In(destinationcreatorpkg.TypeHTTP, destinationcreatorpkg.TypeRFC, destinationcreatorpkg.TypeLDAP, destinationcreatorpkg.TypeMAIL)),
		validation.Field(&b.ProxyType, validation.In(destinationcreatorpkg.ProxyTypeInternet, destinationcreatorpkg.ProxyTypeOnPremise, destinationcreatorpkg.ProxyTypePrivateLink)),
		validation.Field(&b.AuthenticationType, validation.In(destinationcreatorpkg.AuthTypeBasic)),
		validation.Field(&b.User, validation.Required, validation.Length(1, 256)),
		validation.Field(&b.AdditionalProperties, areAdditionalPropertiesValid),
	)
}

// Validate validates that the AuthTypeSAMLAssertion request body contains the required fields and they are valid
func (s *SAMLAssertionDestRequestBody) Validate(destinationCreatorCfg *Config) error {
	areAdditionalPropertiesValid := newDestinationDetailsAdditionalPropertiesValidator(destinationCreatorCfg)

	return validation.ValidateStruct(s,
		validation.Field(&s.Name, validation.Required, validation.Length(1, destinationcreatorpkg.MaxDestinationNameLength), validation.Match(regexp.MustCompile(reqBodyNameRegex))),
		validation.Field(&s.URL, validation.Required),
		validation.Field(&s.Type, validation.In(destinationcreatorpkg.TypeHTTP, destinationcreatorpkg.TypeRFC, destinationcreatorpkg.TypeLDAP, destinationcreatorpkg.TypeMAIL)),
		validation.Field(&s.ProxyType, validation.In(destinationcreatorpkg.ProxyTypeInternet, destinationcreatorpkg.ProxyTypeOnPremise, destinationcreatorpkg.ProxyTypePrivateLink)),
		validation.Field(&s.AuthenticationType, validation.In(destinationcreatorpkg.AuthTypeSAMLAssertion)),
		validation.Field(&s.Audience, validation.Required),
		validation.Field(&s.KeyStoreLocation, validation.Required),
		validation.Field(&s.AdditionalProperties, areAdditionalPropertiesValid),
	)
}

// Validate validates that the AuthTypeClientCertificate request body contains the required fields and they are valid
func (s *ClientCertificateAuthDestRequestBody) Validate(destinationCreatorCfg *Config) error {
	areAdditionalPropertiesValid := newDestinationDetailsAdditionalPropertiesValidator(destinationCreatorCfg)

	return validation.ValidateStruct(s,
		validation.Field(&s.Name, validation.Required, validation.Length(1, destinationcreatorpkg.MaxDestinationNameLength), validation.Match(regexp.MustCompile(reqBodyNameRegex))),
		validation.Field(&s.URL, validation.Required),
		validation.Field(&s.Type, validation.In(destinationcreatorpkg.TypeHTTP, destinationcreatorpkg.TypeRFC, destinationcreatorpkg.TypeLDAP, destinationcreatorpkg.TypeMAIL)),
		validation.Field(&s.ProxyType, validation.In(destinationcreatorpkg.ProxyTypeInternet, destinationcreatorpkg.ProxyTypeOnPremise, destinationcreatorpkg.ProxyTypePrivateLink)),
		validation.Field(&s.AuthenticationType, validation.In(destinationcreatorpkg.AuthTypeClientCertificate)),
		validation.Field(&s.KeyStoreLocation, validation.Required),
		validation.Field(&s.AdditionalProperties, areAdditionalPropertiesValid),
	)
}

// Validate validates that the SAML assertion certificate request body contains the required fields and they are valid
func (c *CertificateRequestBody) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.Length(1, destinationcreatorpkg.MaxDestinationNameLength), validation.Match(regexp.MustCompile(reqBodyNameRegex))),
	)
}

type destinationDetailsAdditionalPropertiesValidator struct {
	destinationCreatorCfg *Config
}

func newDestinationDetailsAdditionalPropertiesValidator(destinationCreatorCfg *Config) *destinationDetailsAdditionalPropertiesValidator {
	return &destinationDetailsAdditionalPropertiesValidator{
		destinationCreatorCfg: destinationCreatorCfg,
	}
}

// Validate is a custom method that validates the correlation IDs, as part of the arbitrary destination additional attributes, are in the expected format - string divided by comma
func (d *destinationDetailsAdditionalPropertiesValidator) Validate(value interface{}) error {
	j, ok := value.(json.RawMessage)
	if !ok {
		return errors.Errorf("Invalid type: %T, expected: %T", value, json.RawMessage{})
	}

	if valid := json.Valid(j); !valid {
		return errors.New("The additional properties json is not valid")
	}

	if d.destinationCreatorCfg == nil {
		return errors.New("The destination creator config could not be empty")
	}

	if d.destinationCreatorCfg.CorrelationIDsKey == "" {
		return errors.New("The correlation IDs key in the destination creator config could not be empty")
	}

	correlationIDsResult := gjson.Get(string(j), d.destinationCreatorCfg.CorrelationIDsKey)
	if !correlationIDsResult.Exists() {
		return errors.New("The correlationIds property part of the additional properties json is required but it does not exist.")
	}

	if correlationIDs := correlationIDsResult.String(); correlationIDs == "" {
		return errors.New("The correlationIds property part of the additional properties could not be empty")
	}

	return nil
}
