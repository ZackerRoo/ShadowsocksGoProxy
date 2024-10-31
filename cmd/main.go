package main

import (
	"proxy_tool/config"
	"proxy_tool/core"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 启动代理
	core.StartProxy()
}
