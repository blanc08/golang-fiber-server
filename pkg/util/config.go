package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Stores all configuration of the application
// The value are read by viper from a config file for environment variables.
type Config struct {
	DBDriver                 string        `mapstructure:"DB_DRIVER"`
	DBSource                 string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress        string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress        string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	V4AsymmetricPublicKeyHex string        `mapstructure:"V4_ASYMMETRIC_PUBLIC_KEY_HEX"`
	V4AsymmetricSecretKeyHex string        `mapstructure:"V4_ASYMMETRIC_SECRET_KEY_HEX"`
	AccessTokenDuration      time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration     time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig read configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("local")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&config)
	return
}
