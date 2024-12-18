package app

import (
	"github.com/dobyte/due-chat/chat/app/logic"
	"github.com/dobyte/due/v2/cluster/node"
)

func Init(proxy *node.Proxy) {
	logic.NewCore(proxy).Init()
}
