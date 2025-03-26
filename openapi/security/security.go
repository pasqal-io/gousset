package security

type Requirement = map[string][]string

type Scheme interface {
	sealed()
}
type Reference string

func Ref(to string) Reference {
	return Reference(to)
}
func (Reference) sealed() {}

var _ Scheme = Reference("")

type Spec struct {
	Type        Type   `json:"type"`
	Description string `json:"description"`
	Name        string `json:"name"`
	*ApiKey
	*Http
	*Oauth2
	*OpenIdConnect
}

type Type string

const (
	TypeAPIKey        = Reference("apiKey")
	TypeHttp          = Reference("http")
	TypeMutalTLS      = Reference("mutualTLS")
	TypeOAuth2        = Reference("oauth2")
	TypeOpenIdConnect = Reference("openIDConnect")
)

type In string

const (
	InQuery  = In("query")
	InHeader = In("header")
	InCookie = In("cookie")
)

type ApiKey struct {
	Name string `json:"string"`
	In   In     `json:"in"`
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
	Implicit          OAuthFlow `json:"implicit"`
	Password          OAuthFlow `json:"password"`
	ClientCredentials OAuthFlow `json:"clientCredentials"`
	AuthorizationCode OAuthFlow `json:"authorizationCode"`
}

type OAuthFlow struct {
	AuthorizationUrl string            `json:"authorizationUrl"`
	TokenUrl         string            `json:"tokenUrl"`
	RefreshUrl       *string           `json:"refreshUrl"`
	Scopes           map[string]string `json:"scopes"`
}
