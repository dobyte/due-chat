package middleware

import (
	"github.com/dobyte/due-chat/internal/code"
	"github.com/dobyte/due-chat/internal/define"
	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/log"
)

func Auth(middleware *node.Middleware, ctx node.Context) {
	if ctx.UID() == 0 {
		if err := ctx.Response(&define.Res{Code: code.Unauthorized.Code()}); err != nil {
			log.Errorf("response message failed, err: %v", err)
		}
	} else {
		middleware.Next(ctx)
	}
}
