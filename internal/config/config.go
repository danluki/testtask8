package config

import (
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv   string `json:"env" yaml:"env" default:"development"`
	Database struct {
		Url string `json:"url" yaml:"url" default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	}
	Http struct {
		Host               string        `json:"host" yaml:"host" default:"0.0.0.0"`
		Port               string        `json:"port" yaml:"port" default:"8080"`
		ReadTimeout        time.Duration `json:"readTimeout" yaml:"readTimeout" default:"5s"`
		WriteTimeout       time.Duration `json:"writeTimeout" yaml:"writeTimeout" default:"5s"`
		MaxHeaderMegabytes int           `json:"maxHeaderMegabytes" yaml:"maxHeaderMegabytes" default:"1"`
	}
}

func DefaultConfig() *Config {
	cfg := &Config{}
	defaults.SetDefaults(cfg)

	return cfg
}

func Load(path string) (*Config, error) {
	if path != "" {
		viper.SetConfigFile(path)
	} else {
		path = GetConfigDefaultPath()
		viper.AddConfigPath(path)
		viper.AddConfigPath("/etc/testtask")
		viper.AddConfigPath("$HOME/.testtask")
	}

	viper.SetEnvPrefix("TestTask")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	conf := &Config{}
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func GetConfigDefaultPath() string {
	dir := getSourcePath()

	return dir + "/../../configs/"
}

func getSourcePath() string {
	_, filename, _, _ := runtime.Caller(1)

	return path.Dir(filename)
}
