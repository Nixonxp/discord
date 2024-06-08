package config

type ApplicationConfig struct {
	Port        string `envconfig:"APP_PORT" default:":8380"`
	MetricsPort string `envconfig:"METRICS_PORT" default:":8382"`
	PostgresDSN string `envconfig:"POSTGRES_DSN" default:"user=admin password=password123 host=localhost port=5432 dbname=discord sslmode=require pool_max_conns=10"`
}
