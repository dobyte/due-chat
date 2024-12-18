package logic

import (
	"context"
	"github.com/dobyte/due-chat/internal/code"
	"github.com/dobyte/due/v2/cluster/node"
	"github.com/dobyte/due/v2/codes"
	"github.com/dobyte/due/v2/errors"
	"github.com/dobyte/due/v2/log"
	"github.com/dobyte/due/v2/utils/xconv"
	"golang.org/x/sync/singleflight"
	"sync"
	"sync/atomic"
)

type Manager struct {
	proxy *node.Proxy        // 代理API
	id    atomic.Int64       // 聊天室自增ID
	sfg   singleflight.Group // single flight
	users sync.Map           // 用户列表
	rooms sync.Map           // 聊天室列表
}

func newManager(proxy *node.Proxy) *Manager {
	return &Manager{proxy: proxy}
}

// 获取成员信息
func (m *Manager) doGetUser(uid int64) (*User, bool) {
	v, ok := m.users.Load(uid)
	if !ok {
		return nil, false
	}

	return v.(*User), true
}

// LoadMember 加载成员信息
func (m *Manager) doLoadUser(uid int64) (*User, error) {
	if v, ok := m.users.Load(uid); ok {
		return v.(*User), nil
	}

	v, err, _ := m.sfg.Do(xconv.String(uid), func() (interface{}, error) {
		user, err := newUser(m, uid)
		if err != nil {
			return nil, err
		}

		m.users.Store(uid, user)

		return user, nil
	})
	if err != nil {
		return nil, err
	}

	return v.(*User), nil
}

// 创建聊天室
func (m *Manager) doCreateRoom(uid int64, name string) (room *Room, err error) {
	var (
		user  *User
		actor *node.Actor
		ctx   = context.Background()
	)

	if user, err = m.doLoadUser(uid); err != nil {
		return
	}

	if user.doHasRoom() {
		return nil, errors.NewError(code.IllegalOperation)
	}

	room = newRoom(m, user, name)

	if err = user.doSaveRoom(room); err != nil {
		return
	}

	defer func() {
		if err != nil {
			user.doClearRoom()
		}
	}()

	if actor, err = m.proxy.Spawn(newRoomProcessor, node.WithActorID(xconv.String(room.id)), node.WithActorArgs(room)); err != nil {
		log.Errorf("spawn actor faile: %v", err)
		return nil, errors.NewError(err, code.InternalError)
	}

	defer func() {
		if err != nil {
			actor.Destroy()
		}
	}()

	if err = m.proxy.BindActor(uid, actor.Kind(), actor.ID()); err != nil {
		log.Errorf("bind actor failed: %v", err)
		return nil, errors.NewError(err, code.InternalError)
	}

	defer func() {
		if err != nil {
			m.proxy.UnbindActor(uid, roomActor)
		}
	}()

	if err = m.proxy.BindNode(ctx, uid); err != nil {
		log.Errorf("bind node failed: %v", err)
		return nil, errors.NewError(err, codes.InternalError)
	}

	m.rooms.Store(room.id, room)

	return
}
