package app

import (
	"github.com/dobyte/due-chat/login/app/api"
	"github.com/dobyte/due/component/http/v2"
)

func Init(proxy *http.Proxy) {
	// API
	a := api.NewAPI(proxy)
	// 路由器
	router := proxy.Router()
	// 登录
	router.Post("/login", a.Login)
	// 注册
	router.Post("/register", a.Register)
}
