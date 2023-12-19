package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Nodes      []string `yaml:"nodes"`       // 全部节点
	Hostname   string   `yaml:"hostname"`    // 当前节点的主机名
	Env        string   `yaml:"env"`         // 当前环境
	HttpServer string   `yaml:"http_server"` // 当前节点的地址
	Protect    bool     `yaml:"protect"`     // 是否开启保护模式
}

func LoadConfig(configFile string) (*Config, error) {
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		return nil, err
	}

	return config, err
}
