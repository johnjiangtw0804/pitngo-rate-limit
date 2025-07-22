package env

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	// App config
	AppName     string `mapstructure:"APP_NAME"`
	AppPort     string `mapstructure:"APP_PORT"`
	AppEnv      string `mapstructure:"APP_Env"`
	AppTimeZone string `mapstructure:"APP_TIMEZONE"`

	// Redis config
	RedisHost       string `mapstructure:"REDIS_HOST"`
	RedisPort       string `mapstructure:"REDIS_PORT"`
	RedisDB         string `mapstructure:"REDIS_DB"`
	RedisPassword   string `mapstructure:"REDIS_PASSWORD"`
	RedisDBProtocol string `mapstructure:"REDIS_DB_PROTOCOL"`
}

func LoadConfig() (*Configuration, error) {
	var config Configuration
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("error to decode, %v", err)
		return nil, err
	}

	return &config, nil
}
