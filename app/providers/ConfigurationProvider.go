package providers

import "github.com/mikeychowy/fiber-crayplate/app/configuration"

var appConfig *configuration.Configuration

func SetConfiguration(config *configuration.Configuration) {
	appConfig = config
}

func GetConfiguration() (config *configuration.Configuration) {
	return appConfig
}
