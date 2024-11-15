package logic

import (
	"github.com/dobyte/due-chat/internal/code"
	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/errors"
	"github.com/dobyte/due/v2/utils/xconv"
	"sync"
	"sync/atomic"
)

type Manager struct {
	id    atomic.Uint64 // 房间自增ID
	rooms sync.Map      // 房间列表
	proxy *node.Proxy   // 代理API
}

func newManager(proxy *node.Proxy) *Manager {
	return &Manager{proxy: proxy}
}

// GetMember 获取成员信息
func (m *Manager) GetMember(uid int64) (*Member, bool) {

}

// LoadMember 加载成员信息
func (m *Manager) LoadMember(uid int64) (*Member, error) {

}

// CreateRoom 创建聊天室
func (m *Manager) CreateRoom(name string, creator *Member) (*Room, error) {
	room := newRoom(m.id.Add(1), name, creator)

	actor, err := m.proxy.Spawn(newRoomProcessor, node.WithActorID(xconv.String(room.id)), node.WithActorArgs(room))
	if err != nil {
		return nil, errors.NewError(err, code.InternalError)
	}

	if err = m.proxy.BindActor(creator.uid, actor.Kind(), actor.ID()); err != nil {
		actor.Destroy()
		return nil, errors.NewError(err, code.InternalError)
	}

	m.rooms.Store(room.id, room)

	return room, nil
}

// EnterRoom 进入聊天室
func (m *Manager) EnterRoom(id uint64, member *Member) (*Room, error) {

}
