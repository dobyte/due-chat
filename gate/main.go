package main

import (
	"github.com/dobyte/due-chat/internal/component/jwt"
	"github.com/dobyte/due/locate/redis/v2"
	"github.com/dobyte/due/network/ws/v2"
	"github.com/dobyte/due/registry/nacos/v2"
	"github.com/dobyte/due/v2"
	"github.com/dobyte/due/v2/cluster/gate"
	"net/http"
)

func main() {
	// 创建容器
	container := due.NewContainer()
	// 创建网络服务器
	server := ws.NewServer()
	// 监听HTTP连接升级WS协议
	server.OnUpgrade(func(w http.ResponseWriter, r *http.Request) bool {
		if _, err := jwt.Instance().Http().ExtractPayload(r); err != nil {
			return false
		}
		return true
	})
	// 创建用户定位器
	locator := redis.NewLocator()
	// 创建服务注册发现
	registry := nacos.NewRegistry()
	// 创建网关组件
	component := gate.NewGate(
		gate.WithServer(server),
		gate.WithLocator(locator),
		gate.WithRegistry(registry),
	)
	// 添加网关组件
	container.Add(component)
	// 启动容器
	container.Serve()
}
