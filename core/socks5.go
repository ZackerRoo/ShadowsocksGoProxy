package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"proxy_tool/config"
	"strconv"

	ss "github.com/shadowsocks/go-shadowsocks2/core"
)

func StartShadowsocksProxy() error {
	port := config.AppConfig.Proxy.Port
	addr := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Shadowsocks 代理正在监听 %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("连接接受失败:", err)
			continue
		}
		go handleShadowsocksConnection(conn)
	}
}

func handleShadowsocksConnection(clientConn net.Conn) {
	defer clientConn.Close()

	targetConn, err := connectToTarget(clientConn.RemoteAddr().String())
	if err != nil {
		log.Println("Shadowsocks 请求失败:", err)
		return
	}
	defer targetConn.Close()

	// 中继数据
	relay(clientConn, targetConn)
}

func connectToTarget(address string) (net.Conn, error) {
	nodes := config.AppConfig.Nodes
	if len(nodes) == 0 {
		return nil, errors.New("无可用代理节点")
	}

	node := nodes[0]
	fullAddress := net.JoinHostPort(node.Address, node.Port)
	log.Printf("通过 Shadowsocks 节点 %s\n 连接 %s", node.Name, fullAddress)

	cipher, err := ss.PickCipher(node.Method, nil, node.Password)
	if err != nil {
		log.Printf("创建 Shadowsocks 加密器失败: %v", err)
		return nil, fmt.Errorf("创建 Shadowsocks 加密器失败: %w", err)
	}

	// 使用 Shadowsocks 连接目标
	proxyConn, err := ss.Dial("tcp", fullAddress, cipher)
	if err != nil {
		log.Printf("连接 Shadowsocks 服务器失败: %v", err)
		return nil, fmt.Errorf("连接 Shadowsocks 服务器失败: %w", err)
	}

	return proxyConn, nil
}

func encodeShadowsocksTarget(address string) []byte {
	host, portStr, _ := net.SplitHostPort(address)
	port, _ := strconv.Atoi(portStr)

	var destAddr bytes.Buffer

	ip := net.ParseIP(host)
	if ip4 := ip.To4(); ip4 != nil {
		destAddr.WriteByte(0x01)
		destAddr.Write(ip4)
	} else if ip6 := ip.To16(); ip6 != nil {
		destAddr.WriteByte(0x04)
		destAddr.Write(ip6)
	} else {
		destAddr.WriteByte(0x03)
		destAddr.WriteByte(byte(len(host)))
		destAddr.WriteString(host)
	}

	binary.Write(&destAddr, binary.BigEndian, uint16(port))

	return destAddr.Bytes()
}

func relay(conn1, conn2 net.Conn) {
	done := make(chan struct{}, 2)
	go func() {
		io.Copy(conn1, conn2)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(conn2, conn1)
		done <- struct{}{}
	}()
	<-done
}
