package config

import (
	"log"

	"github.com/spf13/viper"
)

type ProxyConfig struct {
	Protocol string
	Port     int
}

type Node struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Port     string `json:"port"`     // 端口需要存储为字符串
	Protocol string `json:"protocol"` // 如 ss 或 vmess
	Method   string `json:"method"`   // 加密方式，例如 aes-128-gcm
	Password string `json:"password"` // 密码
}

type Config struct {
	Proxy           ProxyConfig
	SubscriptionURL string `mapstructure:"subscription_url"`
	Nodes           []Node
}

var AppConfig *Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	AppConfig = &config
}
