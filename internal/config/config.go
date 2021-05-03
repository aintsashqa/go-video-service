package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HttpConfig  HttpConfig  `mapstructure:"http"`
	MinioConfig MinioConfig `mapstructure:"minio"`
}

type HttpConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	MaxHeaderMBytes int           `mapstructure:"max_header_mbytes"`
}

type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	EnableSSL       bool   `mapstructure:"enable_ssl"`
}

func Init(filename string) (Config, error) {
	var cfg Config

	viper.AddConfigPath("configs")
	viper.SetConfigName(filename)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	err := viper.Unmarshal(&cfg)
	return cfg, err
}
