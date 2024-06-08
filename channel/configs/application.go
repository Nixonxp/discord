package config

type ApplicationConfig struct {
	Port                       string `envconfig:"APP_PORT" default:":8580"`
	MetricsPort                string `envconfig:"METRICS_PORT" default:":8582"`
	MongoHost                  string `envconfig:"MONGO_HOST" default:"localhost"`
	MongoDb                    string `envconfig:"MONGO_DB" default:"discord"`
	MongoPort                  string `envconfig:"MONGO_PORT" default:"27117"`
	MongoUser                  string `envconfig:"MONGO_USER" default:"discord"`
	MongoPassword              string `envconfig:"MONGO_PASSWORD" default:"example"`
	ServiceCollection          string `envconfig:"MONGO_SERVICE_COLLECTION" default:"channels"`
	ChannelSubscribeCollection string `envconfig:"MONGO_CHANNEL_SUBSCRIBE_COLLECTION" default:"channel_subscribe"`
}
