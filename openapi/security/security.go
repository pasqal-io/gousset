package security

// A list of security requirements for an endpoint.
//
// Each name MUST correspond to a security scheme which is declared in the Security Schemes under the Components Object. If the security scheme is of type "oauth2" or "openIdConnect", then the value is a list of scope names required for the execution, and the list MAY be empty if authorization does not require a specified scope. For other security scheme types, the array MAY contain a list of role names which are required for the execution, but are not otherwise defined or exchanged in-band.
type Requirement = map[string][]string

// A security scheme.
type Scheme interface {
	sealed()
}

// A reference to a security scheme defined in Components.
type Reference struct {
	Ref string `json:"$ref"`
}

func Ref(to string) Reference {
	return Reference{
		Ref: to,
	}
}
func (Reference) sealed() {}

var _ Scheme = Ref("")

// The specification of a security scheme.
type Spec struct {
	// The type of security, e.g. oauth2.
	Type Type `json:"type"`

	// The name of the header, query or cookie parameter used.
	Name string `json:"name"`

	// A human readable description. May contain markdown.
	Description *string `json:"description,omitempty"`

	// Poor man's sum type. Provided iff Type is TypeAPIKey.
	*ApiKey

	// Poor man's sum type. Provided iff Type is TypeHttp.
	*Http

	// Poor man's sum type. Provided iff Type is TypeOAuth2.
	*Oauth2

	// Poor man's sum type. Provided iff Type is TypeOpenIdConnect.
	*OpenIdConnect
}

// A type of security.
type Type string

const (
	TypeAPIKey        = Type("apiKey")
	TypeHttp          = Type("http")
	TypeMutalTLS      = Type("mutualTLS")
	TypeOAuth2        = Type("oauth2")
	TypeOpenIdConnect = Type("openIDConnect")
)

// An emplacement where the authentication may be stored.
type ApiKeyIn string

const (
	ApiKeyInQuery  = ApiKeyIn("query")
	ApiKeyInHeader = ApiKeyIn("header")
	ApiKeyInCookie = ApiKeyIn("cookie")
)

type ApiKey struct {
	// The name of the field containing the API key.
	Name string `json:"string"`

	// Where the API key is stored.
	In ApiKeyIn `json:"in"`
}

type Http struct {
	Scheme       string  `json:"scheme"`
	BearerFormat *string `json:"bearerFormat"`
}

type Oauth2 struct {
	Flows OAuthFlows `json:"flows"`
}

type OpenIdConnect struct {
	OpenIdConnectUrl string `json:"openIdConnectUrl"`
}

type OAuthFlows struct {
	// Configuration for the OAuth implicit flow.
	Implicit *OAuthFlow `json:"implicit,omitempty"`

	// Configuration for the OAuth password flow.
	Password *OAuthFlow `json:"password,omitempty"`

	// Configuration for the OAuth client credentials (aka "application") flow.
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`

	// Configuration for the OAuth authorization code (aka "accesCode") flow.
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

type OAuthFlow struct {
	// The authorization URL to be used for this flow. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	AuthorizationUrl string `json:"authorizationUrl"`

	// The token URL to be used for this flow. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	TokenUrl string `json:"tokenUrl"`

	// The URL to be used for obtaining refresh tokens. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	RefreshUrl *string `json:"refreshUrl,omitempty"`

	// The available scopes for the OAuth2 security scheme. A map between the scope name and a short description for it. The map MAY be empty.
	Scopes map[string]string `json:"scopes"`
}
