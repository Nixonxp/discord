package config

type ApplicationConfig struct {
	Port              string `envconfig:"APP_PORT" default:":8280"`
	MetricsPort       string `envconfig:"METRICS_PORT" default:":8282"`
	UserServiceHost   string `envconfig:"USER_SERVICE_PORT" default:"localhost:8380"`
	AuthSecretKey     string `envconfig:"AUTH_SECRET_KEY" default:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	OAuthClientID     string `envconfig:"OAUTH_CLIENT_ID" default:"{key}"`
	OAuthClientSecret string `envconfig:"OAUTH_CLIENT_SECRET" default:"{secret}"`
	OAuthRedirectUrl  string `envconfig:"OAUTH_REDIRECT_URL" default:"{url}"`
}
