package core

import (
	"log"
)

func StartProxy() {
	err := FetchNodesFromSubscription()
	if err != nil {
		log.Printf("获取订阅节点失败: %v", err)
		return
	}

	log.Println("启动 Shadowsocks 代理")
	if err := StartShadowsocksProxy(); err != nil {
		log.Fatalf("Shadowsocks 代理启动失败: %v", err)
	}
}
