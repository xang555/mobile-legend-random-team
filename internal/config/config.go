package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Server holds HTTP server configuration values.
type Server struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// Logging stores logger configuration options.
type Logging struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

// Team encapsulates team generation configuration.
type Team struct {
	Composition     []string            `mapstructure:"composition"`
	AllowDuplicates bool                `mapstructure:"allow_duplicates"`
	Heroes          map[string][]string `mapstructure:"heroes"`
}

// Config is the top-level application configuration object.
type Config struct {
	Server  Server  `mapstructure:"server"`
	Logging Logging `mapstructure:"logging"`
	Team    Team    `mapstructure:"team"`
}

// Load reads configuration from file and environment variables.
func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	v.SetEnvPrefix("RMT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config: failed to read config %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: failed to unmarshal config %w", err)
	}

	return &cfg, nil
}

// Addr returns the string host:port combination for the server.
func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
