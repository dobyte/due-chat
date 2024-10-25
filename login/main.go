package main

import (
	"github.com/dobyte/due-chat/login/app"
	"github.com/dobyte/due/component/http/v2"
	"github.com/dobyte/due/registry/nacos/v2"
	"github.com/dobyte/due/transport/grpc/v2"
	"github.com/dobyte/due/v2"
)

// @title 登录服API文档
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	// 创建容器
	container := due.NewContainer()
	// 创建服务注册发现
	registry := nacos.NewRegistry()
	// 创建RPC传输器
	transporter := grpc.NewTransporter()
	// 创建HTTP组件
	component := http.NewHttp(
		http.WithRegistry(registry),
		http.WithTransporter(transporter),
	)
	// 初始化应用
	app.Init(component.Proxy())
	// 添加HTTP组件
	container.Add(component)
	// 启动容器
	container.Serve()
}
