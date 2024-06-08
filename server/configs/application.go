package config

type ApplicationConfig struct {
	Port                      string `envconfig:"APP_PORT" default:":8480"`
	MetricsPort               string `envconfig:"METRICS_PORT" default:":8482"`
	ChatServiceHost           string `envconfig:"METRICS_PORT" default:":8682"`
	MongoHost                 string `envconfig:"MONGO_HOST" default:"localhost"`
	MongoDb                   string `envconfig:"MONGO_DB" default:"discord"`
	MongoPort                 string `envconfig:"MONGO_PORT" default:"27117"`
	MongoUser                 string `envconfig:"MONGO_USER" default:"discord"`
	MongoPassword             string `envconfig:"MONGO_PASSWORD" default:"example"`
	ServiceCollection         string `envconfig:"MONGO_SERVICE_COLLECTION" default:"servers"`
	ServerSubscribeCollection string `envconfig:"MONGO_SERVICE_COLLECTION" default:"server_subscribe"`
}
