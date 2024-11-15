package logic

import (
	"github.com/dobyte/due-chat/chat/app/route"
	"github.com/dobyte/due-chat/internal/code"
	"github.com/dobyte/due-chat/internal/middleware"
	usersvc "github.com/dobyte/due-chat/internal/service/user/client"
	userpb "github.com/dobyte/due-chat/internal/service/user/pb"
	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/codes"
	"github.com/dobyte/due/v2/log"
)

type Core struct {
	proxy   *node.Proxy
	manager *Manager
}

func NewCore(proxy *node.Proxy) *Core {
	return &Core{proxy: proxy, manager: newManager(proxy)}
}

func (c *Core) Init() {
	c.proxy.Router().Group(func(group *node.RouterGroup) {
		// 验证TOKEN
		group.AddRouteHandler(route.ValidateToken, false, c.validateToken)
		// 注入授权检测中间件
		group.Middleware(middleware.Auth)
		// 拉取玩家信息
		//group.AddRouteHandler(route.FetchProfile, false, c.fetchProfile)
		// 创建聊天室
		group.AddRouteHandler(route.CreateRoom, false, c.createRoom)
		// 进入聊天室
		group.AddRouteHandler(route.EnterRoom, false, c.enterRoom)
		// 离开聊天室
		group.AddRouteHandler(route.LeaveRoom, true, c.leaveRoom)
		// 发送消息
		group.AddRouteHandler(route.SendMessage, true, c.sendMessage)
	})
}

// 验证TOKEN
func (c *Core) validateToken(ctx node.Context) {
	ctx.Task(func(ctx node.Context) {
		req := &ValidateTokenReq{}
		res := &ValidateTokenRes{}
		ctx.Defer(func() {
			if err := ctx.Response(res); err != nil {
				log.Errorf("response message failed: %v", err)
			}
		})

		if err := ctx.Parse(req); err != nil {
			log.Errorf("parse request message failed: %v", err)
			res.Code = code.InternalError.Code()
			return
		}

		if req.Token == "" {
			res.Code = code.InvalidArgument.Code()
			return
		}

		client, err := usersvc.NewClient(c.proxy.NewMeshClient)
		if err != nil {
			log.Errorf("create client failed: %v", err)
			res.Code = code.InternalError.Code()
			return
		}

		reply, err := client.ValidateToken(ctx.Context(), &userpb.ValidateTokenArgs{
			Token: req.Token,
		})
		if err != nil {
			res.Code = codes.Convert(err).Code()
			return
		}

		if err = ctx.BindGate(reply.UID); err != nil {
			log.Errorf("bind gate failed, uid = %v err = %v", reply.UID, err)
			res.Code = code.InternalError.Code()
			return
		}

		res.Code = code.OK.Code()
	})
}

// 创建聊天室
func (c *Core) createRoom(ctx node.Context) {
	req := &CreateRoomReq{}
	res := &CreateRoomRes{}
	ctx.Defer(func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response message failed: %v", err)
		}
	})

	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
		res.Code = code.InternalError.Code()
		return
	}

	var member *Member

	room, err := c.manager.CreateRoom(req.Name, member)
	if err != nil {
		res.Code = codes.Convert(err).Code()
		return
	}

	res.Code = code.OK.Code()
	res.Data = &CreateRoomResData{RoomID: room.id}
}

// 进入聊天室
func (c *Core) enterRoom(ctx node.Context) {
	req := &EnterRoomReq{}
	res := &EnterRoomRes{}
	ctx.Defer(func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response message failed: %v", err)
		}
	})

	if err := ctx.Parse(req); err != nil {
		log.Errorf("parse request message failed: %v", err)
		res.Code = code.InternalError.Code()
		return
	}

	member, err := c.manager.LoadMember(ctx.UID())
	if err != nil {
		res.Code = codes.Convert(err).Code()
		return
	}

	room, err := c.manager.EnterRoom(req.RoomID, member)
	if err != nil {
		res.Code = codes.Convert(err).Code()
		return
	}

	res.Code = code.OK.Code()
	res.Data = &EnterRoomResData{ID: room.id, Name: room.name}
}

// 离开聊天室
func (c *Core) leaveRoom(ctx node.Context) {
	res := &LeaveRoomRes{}
	ctx.Defer(func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response message failed: %v", err)
		}
	})

	if err := ctx.Next(); err != nil {
		log.Errorf("request next failed: %v", err)
		res.Code = code.IllegalRequest.Code()
	}
}

// 发送消息
func (c *Core) sendMessage(ctx node.Context) {
	res := &SendMessageRes{}
	ctx.Defer(func() {
		if err := ctx.Response(res); err != nil {
			log.Errorf("response message failed: %v", err)
		}
	})

	if err := ctx.Next(); err != nil {
		log.Errorf("request next failed: %v", err)
		res.Code = code.IllegalRequest.Code()
	}
}
