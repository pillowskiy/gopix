package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   Server
	Logger   Logger
	Postgres Postgres
}

type Server struct {
	Addr         string        `mapstructure:"addr"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	CtxTimeout   time.Duration `mapstructure:"ctx_timeout"`
	Debug        bool          `mapstructure:"debug"`
}

type Logger struct {
	Mode              string `mapstructure:"mode"`
	DisableCaller     bool   `mapstructure:"disable_caller"`
	DisableStacktrace bool   `mapstructure:"disable_stacktrace"`
	Encoding          string `mapstructure:"encoding"`
	Level             string `mapstructure:"level"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSL      bool   `mapstructure:"ssl"`
	Driver   string `mapstructure:"driver"`
}

func FetchAndLoadConfig() (*viper.Viper, error) {
	path := fetchConfigPath()
	if path == "" {
		return nil, fmt.Errorf("unable to fetch config path")
	}

	configParts := strings.Split(path, "/")
	configFilename := configParts[len(configParts)-1]
	configPath := strings.Join(configParts[:len(configParts)-1], "/")
	return LoadConfig(configFilename, configPath)
}

func LoadConfig(filename string, path string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(path)
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file with path '%s/%s' not found", path, filename)
		}
		return nil, err
	}

	return v, nil
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return &c, nil
}
