package api

import (
	"context"
	"github.com/dobyte/due-chat/internal/code"
	usersvc "github.com/dobyte/due-chat/internal/service/user/client"
	userpb "github.com/dobyte/due-chat/internal/service/user/pb"
	"github.com/dobyte/due/component/http/v2"
	"github.com/dobyte/due/v2/log"
	"github.com/go-playground/validator/v10"
)

type API struct {
	proxy    *http.Proxy
	validate *validator.Validate
}

func NewAPI(proxy *http.Proxy) *API {
	return &API{
		proxy:    proxy,
		validate: validator.New(),
	}
}

// Login 登录
// @Summary 登录
// @Tags 登录
// @Schemes
// @Accept json
// @Produce json
// @Param request body LoginReq true "请求参数"
// @Response 200 {object} http.Resp{Data=LoginRes} "响应参数"
// @Router /login [post]
func (a *API) Login(ctx http.Context) error {
	req := &LoginReq{}

	if err := ctx.Bind().JSON(req); err != nil {
		return ctx.Failure(code.InvalidArgument)
	}

	if err := a.validate.Struct(req); err != nil {
		return ctx.Failure(code.InvalidArgument)
	}

	client, err := usersvc.NewClient(a.proxy.NewMeshClient)
	if err != nil {
		log.Errorf("create client failed: %v", err)
		return ctx.Failure(code.InternalError)
	}

	reply, err := client.Login(context.Background(), &userpb.LoginArgs{
		Account:  req.Account,
		Password: req.Password,
		ClientIP: ctx.IP(),
	})
	if err != nil {
		return ctx.Failure(err)
	}

	return ctx.Success(&LoginRes{
		Gate:  reply.Gate,
		Token: reply.Token,
	})
}

// Register 注册
// @Summary 注册
// @Tags 注册
// @Schemes
// @Accept json
// @Produce json
// @Param request body RegisterReq true "请求参数"
// @Response 200 {object} http.Resp{} "响应参数"
// @Router /register [post]
func (a *API) Register(ctx http.Context) error {
	req := &RegisterReq{}

	if err := ctx.Bind().JSON(req); err != nil {
		return ctx.Failure(code.InvalidArgument)
	}

	if err := a.validate.Struct(req); err != nil {
		return ctx.Failure(code.InvalidArgument)
	}

	client, err := usersvc.NewClient(a.proxy.NewMeshClient)
	if err != nil {
		log.Errorf("create client failed: %v", err)
		return ctx.Failure(code.InternalError)
	}

	_, err = client.Register(context.Background(), &userpb.RegisterArgs{
		Account:  req.Account,
		Password: req.Password,
		Nickname: req.Nickname,
		ClientIP: ctx.IP(),
	})
	if err != nil {
		return ctx.Failure(err)
	}

	return ctx.Success()
}
