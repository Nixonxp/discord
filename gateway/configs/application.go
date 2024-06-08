package config

type ApplicationConfig struct {
	Port               string `envconfig:"APP_PORT" default:":8180"`
	SwaggerPort        string `envconfig:"SWAGGER_PORT" default:":8081"`
	MetricsPort        string `envconfig:"METRICS_PORT" default:":8082"`
	AuthServiceHost    string `envconfig:"AUTH_SERVICE_PORT" default:"localhost:8280"`
	UserServiceHost    string `envconfig:"USER_SERVICE_PORT" default:"localhost:8380"`
	ServerServiceHost  string `envconfig:"SERVER_SERVICE_PORT" default:"localhost:8480"`
	ChannelServiceHost string `envconfig:"CHANNEL_SERVICE_PORT" default:"localhost:8580"`
	ChatServiceHost    string `envconfig:"CHAT_SERVICE_PORT" default:"localhost:8680"`
	AuthSecretKey      string `envconfig:"AUTH_SECRET_KEY" default:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}
