package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config - Configuration for app (both client and server for convenience)
type Config struct {
	Server   *ServerConf     `mapstructure:"server"`
	Hashcash *HashcashConf   `mapstructure:"hashcash"`
	Client   *ClientConnConf `mapstructure:"client"`
}

type ServerConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ClientConnConf struct {
	ChallengeLen  int           `mapstructure:"challenge_len"`
	QueueLenLimit int           `mapstructure:"queue_len_limit"`
	WriteTimeout  time.Duration `mapstructure:"write_timeout"`

	DebugRead  bool `mapstructure:"debug_read"`
	DebugWrite bool `mapstructure:"debug_write"`
}

type HashcashConf struct {
	ZerosCount    int `mapstructure:"zeros_count"`
	MaxIterations int `mapstructure:"max_iterations"`
}

func ReadConfig() (*Config, error) {
	viper.AddConfigPath("./configs")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	config := &Config{}
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
