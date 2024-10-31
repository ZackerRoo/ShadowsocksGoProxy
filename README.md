# ShadowsocksGoProxy

ShadowsocksGoProxy 是一个基于 Go 的代理工具，支持 Shadowsocks 协议并集成了 SOCKS5 代理功能。该工具旨在提供高效、安全的网络代理服务，适用于各种网络环境和需求。

## 特性

- **支持 Shadowsocks**：使用 Shadowsocks 协议加密数据，提供安全的代理连接。
- **支持 SOCKS5**：实现 SOCKS5 代理，灵活处理网络请求。
- **节点订阅**：从订阅链接获取代理节点信息，支持多节点选择。
- **高并发**：利用 Go 的 goroutines 处理高并发连接。

## 安装与使用

1. 克隆仓库：

   ```sh
   git clone https://github.com/ZackerRoo/ShadowsocksGoProxy.git
   ```

2. 构建并运行项目：

   ```sh
   cd proxy_tool
   go build -o proxy_tool
   ./proxy_tool
   ```

3. 使用代理（例如，使用 curl 进行测试）：

   ```sh
   curl -x socks5h://127.0.0.1:1080 https://www.google.com
   ```

## 配置

配置文件 `config.json` 支持以下字段：

- `SubscriptionURL`：用于订阅节点列表的 URL。
- `Proxy`：
  - `Port`：代理服务器监听的端口。
  - `Protocol`：代理协议类型（目前支持 `shadowsocks` 和 `socks5`）。

## 贡献

欢迎贡献代码！如果您有改进建议或遇到任何问题，请随时创建 issue 或提交 pull request。

## 许可证

本项目使用 MIT 许可证。
