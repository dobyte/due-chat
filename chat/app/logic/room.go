package logic

import (
	"context"
	"github.com/dobyte/due-chat/chat/app/route"
	"github.com/dobyte/due-chat/internal/code"
	"github.com/dobyte/due/v2/cluster"
	"github.com/dobyte/due/v2/errors"
	"github.com/dobyte/due/v2/session"
	"github.com/dobyte/due/v2/utils/xtime"
	"sync"
)

type Room struct {
	id      uint64            // 房间ID
	name    string            // 房间名称
	creator *Member           // 创建者
	manager *Manager          // 管理器
	rw      sync.RWMutex      // 锁
	members map[int64]*Member // 房间成员列表
}

func newRoom(id uint64, name string, creator *Member) *Room {
	r := &Room{}
	r.id = id
	r.name = name
	r.creator = creator
	r.members = make(map[int64]*Member, 1)
	r.members[creator.uid] = creator

	return r
}

// SendMessage 发送消息
func (r *Room) SendMessage(ctx context.Context, sender *Member, content string) error {
	r.rw.RLock()
	targets := make([]int64, 0, len(r.members))
	for uid := range r.members {
		if uid != sender.uid {
			targets = append(targets, uid)
		}
	}
	r.rw.RUnlock()

	if len(targets) == 0 {
		return nil
	}

	err := r.manager.proxy.Broadcast(ctx, &cluster.BroadcastArgs{
		Kind: session.User,
		Message: &cluster.Message{
			Route: route.MessageNotify,
			Data: &MessageNotify{
				Sender: &MemberInfo{
					UID:      sender.uid,
					Nickname: sender.nickname,
				},
				Content:   content,
				Timestamp: xtime.Now().Unix(),
			},
		},
	})
	if err != nil {
		return errors.NewError(err, code.InternalError)
	}

	return nil
}
