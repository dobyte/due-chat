package app

import (
	usersvc "github.com/dobyte/due-chat/internal/service/user/server"
	"github.com/dobyte/due/v2/cluster/mesh"
)

func Init(proxy *mesh.Proxy) {
	// 初始化用户服务
	usersvc.NewServer(proxy).Init()
}
