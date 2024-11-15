package main

import (
	"github.com/dobyte/due-chat/chat/app"
	"github.com/dobyte/due/locate/redis/v2"
	"github.com/dobyte/due/registry/nacos/v2"
	"github.com/dobyte/due/transport/grpc/v2"
	"github.com/dobyte/due/v2"
	"github.com/dobyte/due/v2/cluster/node"
)

func main() {
	// 创建容器
	container := due.NewContainer()
	// 创建用户定位器
	locator := redis.NewLocator()
	// 创建服务注册发现
	registry := nacos.NewRegistry()
	// 创建RPC传输器
	transporter := grpc.NewTransporter()
	// 创建节点组件
	component := node.NewNode(
		node.WithLocator(locator),
		node.WithRegistry(registry),
		node.WithTransporter(transporter),
	)
	// 初始化应用
	app.Init(component.Proxy())
	// 添加节点组件
	container.Add(component)
	// 启动容器
	container.Serve()
}
