package configurator

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

func ReadEnvironment(config Configurable) {
	err := envconfig.Process("", config)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
