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
	Server     Server     `mapstructure:"server"`
	Logger     Logger     `mapstructure:"logger"`
	Postgres   Postgres   `mapstructure:"postgres"`
	Redis      Redis      `mapstructure:"redis"`
	S3         S3         `mapstructure:"s3"`
	VecService VecService `mapstructure:"vec_service"`
	Metrics    Metrics    `mapstructure:"metrics"`
	OAuth      OAuth      `mapstructure:"oauth"`
}

type Server struct {
	Addr         string        `mapstructure:"addr"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	CtxTimeout   time.Duration `mapstructure:"ctx_timeout"`
	Debug        bool          `mapstructure:"debug"`
	Cookie       *Cookie       `mapstructure:"cookie"`
	Session      *Session      `mapstructure:"session"`
	CORS         *CORS         `mapstructure:"cors"`
}

type Cookie struct {
	Name     string        `mapstructure:"name"`
	Expire   time.Duration `mapstructure:"expire"`
	Secure   bool          `mapstructure:"secure"`
	HttpOnly bool          `mapstructure:"http_only"`
}

type Session struct {
	Expire time.Duration `mapstructure:"expire"`
	Secret string        `mapstructure:"secret"`
}

type CORS struct {
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowOrigins     []string `mapstructure:"allow_origins"`
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

type Redis struct {
	RedisAddr      string        `mapstructure:"redis_addr"`
	RedisPass      string        `mapstructure:"redis_pass"`
	RedisDB        int           `mapstructure:"redis_db"`
	RedisDefaultDB int           `mapstructure:"redis_default_db"`
	MinIdleConns   int           `mapstructure:"min_idle_conns"`
	PoolSize       int           `mapstructure:"pool_size"`
	PoolTimeout    time.Duration `mapstructure:"pool_timeout"`
}

type S3 struct {
	Endpoint       string `mapstructure:"endpoint"`
	Bucket         string `mapstructure:"bucket"`
	Region         string `mapstructure:"region"`
	AccessKey      string `mapstructure:"access_key"`
	SecretAccess   string `mapstructure:"secret_access_key"`
	ForcePathStyle bool   `mapstructure:"force_path_style"`
	// The buffer size for file uploads, including multipart uploads, in megabytes.
	UploadBufferSizeMB   int   `mapstructure:"upload_buffer_size_mb"`
	MultipartChunkSizeMB int64 `mapstructure:"multipart_chunk_size_mb"`
}

type OAuth struct {
	Google *OAuthGoogle `mapstructure:"google"`
}

type OAuthGoogle struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type VecService struct {
	URL string `mapstructure:"url"`
}

type Metrics struct {
	URL  string `mapstructure:"url"`
	Name string `mapstructure:"name"`
}

func FetchAndLoadConfig() (*Config, error) {
	path := fetchConfigPath()
	if path == "" {
		return nil, fmt.Errorf("unable to fetch config path")
	}

	configParts := strings.Split(path, "/")
	configFilename := configParts[len(configParts)-1]
	configPath := strings.Join(configParts[:len(configParts)-1], "/")
	return LoadViper(configFilename, configPath)
}

func LoadViper(filename string, path string) (*Config, error) {
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

	return ParseConfig(v)
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
