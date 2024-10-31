package core

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"proxy_tool/config"
	"strings"
)

func FetchNodesFromSubscription() error {
	url := config.AppConfig.SubscriptionURL
	if url == "" {
		return errors.New("未设置订阅地址")
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("订阅地址返回错误状态码")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	decodedData, err := base64.StdEncoding.DecodeString(string(body))

	if err != nil {
		return errors.New("订阅内容Base64解码失败: " + err.Error())
	}
	lines := strings.Split(string(decodedData), "\n")
	// fmt.Printf("lines: %v\n", lines)

	var nodes []config.Node
	for _, line := range lines {
		if line == "" {
			continue
		}
		node, err := parseNode(line)
		if err != nil {
			log.Printf("解析节点失败: %v", err)
			continue
		}
		nodes = append(nodes, node)
	}

	// for _, node := range nodes {
	// 	log.Printf("解析到节点: %s", node)
	// }

	config.AppConfig.Nodes = nodes
	log.Printf("成功从订阅地址获取到 %d 个节点", len(nodes))
	return nil
}

func parseNode(line string) (config.Node, error) {
	if strings.HasPrefix(line, "ss://") {
		return parseSSNode(line)
	} else if strings.HasPrefix(line, "vmess://") {
		log.Println("不支持的节点类型: vmess")
	}
	return config.Node{}, errors.New("不支持的协议类型")
}

func parseSSNode(line string) (config.Node, error) {
	line = strings.TrimPrefix(line, "ss://")

	parts := strings.SplitN(line, "@", 2)
	if len(parts) < 2 {
		return config.Node{}, errors.New("无效的 ss:// 节点格式")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return config.Node{}, errors.New("Base64 解码失败: " + err.Error())
	}
	encryptionParts := strings.SplitN(string(decoded), ":", 2)
	if len(encryptionParts) != 2 {
		return config.Node{}, errors.New("加密方式或密码解析失败")
	}
	method := encryptionParts[0]
	password := encryptionParts[1]

	addressParts := strings.Split(parts[1], ":")
	if len(addressParts) != 2 {
		return config.Node{}, errors.New("服务器地址或端口解析失败")
	}
	address := addressParts[0]
	port := strings.Split(addressParts[1], "#")[0]

	name := ""
	if strings.Contains(line, "#") {
		nameParts := strings.Split(addressParts[1], "#")
		if len(nameParts) == 2 {
			name = nameParts[1]
		}
	}

	return config.Node{
		Name:     name,
		Address:  address,
		Port:     port,
		Protocol: "ss",
		Method:   method,
		Password: password,
	}, nil
}
