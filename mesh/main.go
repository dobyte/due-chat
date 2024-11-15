package main

import (
	"github.com/dobyte/due-chat/mesh/app"
	"github.com/dobyte/due/locate/redis/v2"
	"github.com/dobyte/due/registry/nacos/v2"
	"github.com/dobyte/due/transport/grpc/v2"
	"github.com/dobyte/due/v2"
	"github.com/dobyte/due/v2/cluster/mesh"
	"github.com/dobyte/due/v2/component/pprof"
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
	// 创建网格组件
	component := mesh.NewMesh(
		mesh.WithLocator(locator),
		mesh.WithRegistry(registry),
		mesh.WithTransporter(transporter),
	)
	// 初始化应用
	app.Init(component.Proxy())
	// 添加网格组件
	container.Add(component, pprof.NewPProf())
	// 启动容器
	container.Serve()
}
