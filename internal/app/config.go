package app

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Server struct {
	Host   string `validate:"required"`
	Port   int    `validate:"required"`
	Enable bool   `validate:"required"`
}

type FTPClient struct {
	Addr            Server `validate:"required"`
	User            string `validate:"required"`
	Pass            string `validate:"required"`
	RemoteDirectory string `validate:"required"`
}

type Stage struct {
	Positions         int    `validate:"required"`
	AttemptBeforeNext int    `validate:"required"`
	Threads           int    `validate:"required"`
	IsCookie          bool   `validate:"required"`
	IsEmbededCharters bool   `validate:"required"`
	IsMarkWatched     bool   `validate:"required"`
	Extractors        string `validate:"required"`
}
type DownloadConfig struct {
	Threads       int     `validate:"required"`
	PercentToNext int     `validate:"required"`
	WorkPath      string  `validate:"required"`
	Stages        []Stage `validate:"required"`
}

type Config struct {
	GRPC Server    `validate:"required"`
	HTTP Server    `validate:"required"`
	FTP  FTPClient `validate:"required"`

	PathToDB string `validate:"required"`
	NameDB   string `validate:"required"`

	Cache Server `validate:"required"`

	Downloader DownloadConfig `validate:"required"`
}

func ReadConfig(path string, fileName string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("ReadConfig: %w", err)
	}

	b, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		return nil, fmt.Errorf("ReadConfig (Marshal): %w", err)
	}

	res := &Config{}
	err = yaml.Unmarshal(b, res)
	if err != nil {
		return nil, fmt.Errorf("ReadConfig (Unmarshal): %w", err)
	}

	if err := valid(res); err != nil {
		return nil, fmt.Errorf("ReadConfig (Validate): %w", err)
	}

	return res, nil
}

func valid(conf *Config) error {
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(conf); err != nil {
		return fmt.Errorf("validation faild: %w", err)
	}

	return nil
}
