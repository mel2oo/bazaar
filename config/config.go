package config

import (
	"bazaar/internal/domain/db"
	"bazaar/internal/domain/storage"
	"bazaar/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Server     Server         `mapstructure:"server"`
	Logger     logger.Config  `mapstructure:"logger"`
	Storage    storage.Config `mapstructure:"storage"`
	Counchbase db.Config      `mapstructure:"database"`
	Yara       struct {
		Address string `mapstructure:"address"`
	} `mapstructure:"yara"`
}

type Server struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Address string `mapstructure:"address"`
}

// 读取配置文件，反序列化到指定结构体
func Load(path string) (*Config, error) {
	c := new(Config)

	// 添加配置文件目录
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}

	return c, nil
}
