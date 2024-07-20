package configs

import (
	"github.com/spf13/viper"
)

type conf struct {
	WebServerPort          string `mapstructure:"WEB_SERVER_PORT"`
	OtelServiceName        string `mapstructure:"OTEL_SERVICE_NAME"`
	OtelCollectorURL       string `mapstructure:"OTEL_COLLECTOR_URL"`
	ViaCEPClientURL        string `mapstructure:"VIA_CEP_CLIENT_URL"`
	WeatherAPIClientURL    string `mapstructure:"WEATHER_API_CLIENT_URL"`
	WeatherAPIClientAPIKey string `mapstructure:"WEATHER_API_CLIENT_API_KEY"`
}

type confValidator struct {
	WebServerPort      string `mapstructure:"WEB_SERVER_PORT"`
	OtelServiceName    string `mapstructure:"OTEL_SERVICE_NAME"`
	OtelCollectorURL   string `mapstructure:"OTEL_COLLECTOR_URL"`
	TemppcAPICleintURL string `mapstructure:"TEMPPC_API_CLIENT_URL"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("temppc_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, nil
}

func LoadConfigValidator(path string) (*confValidator, error) {
	var cfg *confValidator
	viper.SetConfigName("validator_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, nil
}
