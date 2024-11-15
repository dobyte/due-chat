package logic

import (
	"github.com/dobyte/due-chat/chat/app/route"
	"github.com/dobyte/due-chat/internal/code"
	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/codes"
	"github.com/dobyte/due/v2/log"
)

const roomActor = "Room"

type roomProcessor struct {
	node.BaseProcessor
	actor *node.Actor
	room  *Room
}

func newRoomProcessor(actor *node.Actor, args ...any) node.Processor {
	return &roomProcessor{
		actor: actor,
		room:  args[0].(*Room),
	}
}

// Kind 设置处理器类型
func (p *roomProcessor) Kind() string {
	return roomActor
}

// Init 初始化处理器
func (p *roomProcessor) Init() {
	// 进入聊天室
	p.actor.AddRouteHandler(route.EnterRoom, p.enterRoom)
	// 离开聊天室
	p.actor.AddRouteHandler(route.LeaveRoom, p.leaveRoom)
	// 发送消息
	p.actor.AddRouteHandler(route.SendMessage, p.sendMessage)
	// 拉取成员
	p.actor.AddRouteHandler(route.FetchMembers, p.fetchMembers)
}

// 进入聊天室
func (p *roomProcessor) enterRoom(ctx node.Context) {
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
}

// 进入聊天室
func (p *roomProcessor) leaveRoom(ctx node.Context) {

}

// 发送消息
func (p *roomProcessor) sendMessage(ctx node.Context) {
	req := &SendMessageReq{}
	res := &SendMessageRes{}
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

	member, ok := p.room.manager.GetMember(ctx.UID())
	if !ok {
		res.Code = code.IllegalOperation.Code()
		return
	}

	if err := p.room.SendMessage(ctx.Context(), member, req.Content); err != nil {
		res.Code = codes.Convert(err).Code()
		return
	}

	res.Code = code.OK.Code()
}
