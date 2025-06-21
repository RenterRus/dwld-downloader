package app

import (
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Server struct {
	Host   string
	Port   int
	Enable bool
}

type FTPClient struct {
	Addr            Server
	User            string
	Pass            string
	RemoteDirectory string
}

type Config struct {
	GRPC Server
	HTTP Server
	FTP  FTPClient

	PathToDB string
	Cache    Server

	EagerMode bool
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

	return res, nil
}
