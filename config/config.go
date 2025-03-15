package config

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Database Database `yaml:"database"`
		Secret   Secret
		Serving  Serving `yaml:"serving"`
	}
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Username string `yaml:"username" env:"DB_USERNAME"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
		Name     string `yaml:"name" env:"DB_NAME"`
	}
	Secret struct {
		BotToken           string `env:"BOT_TOKEN" envDefault:""`
		GitHubToken        string `env:"GITHUB_TOKEN" envDefault:""`
		StackOverflowToken string `env:"STACKOVERFLOW_TOKEN" envDefault:""`
	}
	Serving struct {
		Host         string        `yaml:"host" env:"HOST" envDefault:"localhost"`
		ScrapperPort int           `yaml:"scrapper_port" env:"SCRAPPER_PORT" envDefault:"8080"`
		BotPort      int           `yaml:"bot_port" env:"BOT_PORT" envDefault:"8081"`
		Debug        bool          `yaml:"debug" env:"DEBUG"`
		Interval     time.Duration `yaml:"interval" env:"INTERVAL" envDefault:"5"`
	}
)

func (d *Database) ToDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d%s?target_session_attrs=read-write&ssl-mode=disable",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}

func NewConfig(name string) (*Config, error) {
	cfg, err := NewConfigFromFile(name)
	if err != nil {
		return nil, err
	}

	if err := configFromEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

//go:embed default-config.yaml
var configBytes []byte

func NewConfigFromFile(name string) (*Config, error) {
	conf := &Config{}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(name)

	if err := v.ReadConfig(bytes.NewReader(configBytes)); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if err := v.MergeInConfig(); err != nil {
		if errors.Is(err, &viper.ConfigParseError{}) {
			return nil, fmt.Errorf("merge config: %w", err)
		}
	}

	if err := v.Unmarshal(conf); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return conf, nil
}

func configFromEnv(conf *Config) error {
	if err := env.Parse(conf); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	return nil
}
