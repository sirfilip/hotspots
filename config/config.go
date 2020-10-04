package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvPrefix("necsam")
	viper.AutomaticEnv()
	viper.SetDefault("environment", "development")
	viper.SetDefault("port", ":3000")
	viper.SetDefault("app_name", "Hotspots")
	viper.SetDefault("api_url", "http://localhost:3000")
	viper.SetDefault("email_host", "127.0.0.1")
	viper.SetDefault("email_host_user", "tester@example.com")
	viper.SetDefault("email_host_password", "tester")
	viper.SetDefault("email_port", "1025")
	viper.SetDefault("email_use_tls", "true")
	viper.SetDefault("email_tls_skip_verify", "false")
	viper.SetDefault("app_secret", "hs-230dfce7354c18cc61ec63e0352de3aa")
	viper.SetDefault("dbname", fmt.Sprintf("necsam_%s", Get("environment")))
	viper.SetDefault("mongourl", "mongodb://localhost:27017")
}

func Get(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int64 {
	return viper.GetInt64(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}
