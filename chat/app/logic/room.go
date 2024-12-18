package logic

import (
	"context"
	"github.com/dobyte/due-chat/chat/app/route"
	"github.com/dobyte/due-chat/internal/code"
	"github.com/dobyte/due/v2/cluster"
	"github.com/dobyte/due/v2/errors"
	"github.com/dobyte/due/v2/log"
	"github.com/dobyte/due/v2/session"
	"github.com/dobyte/due/v2/utils/xconv"
	"github.com/dobyte/due/v2/utils/xtime"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
)

const (
	stateNormal  = 0 // 正常
	stateDestroy = 1 // 销毁中
)

type Room struct {
	id      int64             // 聊天室ID
	name    string            // 聊天室名称
	creator *Member           // 创建者
	manager *Manager          // 管理器
	state   atomic.Int32      // 房间状态
	rw      sync.RWMutex      // 锁
	members map[int64]*Member // 聊天室成员列表
}

func newRoom(manager *Manager, user *User, name string) *Room {
	r := &Room{}
	r.id = manager.id.Add(1)
	r.name = name
	r.creator = &Member{user: user, time: xtime.Now()}
	r.manager = manager
	r.members = make(map[int64]*Member, 1)
	r.members[user.uid] = r.creator

	return r
}

// 发送消息
func (r *Room) doSendMessage(sender int64, content string) error {
	user, ok := r.manager.doGetUser(sender)
	if !ok {
		return errors.NewError(code.IllegalOperation)
	}

	r.rw.RLock()

	if r.state.Load() != stateNormal {
		r.rw.RLock()
		return errors.NewError(code.IllegalOperation)
	}

	targets := make([]int64, 0, len(r.members))
	for uid := range r.members {
		targets = append(targets, uid)
	}

	r.rw.RUnlock()

	if len(targets) == 0 {
		return nil
	}

	if err := r.manager.proxy.Multicast(context.Background(), &cluster.MulticastArgs{
		Kind:    session.User,
		Targets: targets,
		Message: &cluster.Message{
			Route: route.MessageNotify,
			Data: &MessageNotify{
				Sender:    user.doMakeUserInfo(),
				Content:   content,
				Timestamp: xtime.Now().Unix(),
			},
		},
	}); err != nil && !errors.Is(err, errors.ErrNotFoundUserLocation) {
		return errors.NewError(err, code.InternalError)
	}

	return nil
}

// 解散房间
func (r *Room) doDismissRoom(uid int64) error {
	r.rw.Lock()

	if member, ok := r.members[uid]; !ok && member != r.creator {
		r.rw.Unlock()
		return errors.NewError(code.IllegalOperation)
	}

	if !r.state.CompareAndSwap(stateNormal, stateDestroy) {
		r.rw.Unlock()
		return errors.NewError(code.IllegalOperation)
	}

	r.rw.Unlock()

	eg := &errgroup.Group{}

	for u := range r.members {
		member := r.members[u]

		eg.Go(func() error {
			member.user.doClearRoom()

			r.manager.proxy.UnbindActor(member.user.uid, roomActor)

			if err := r.manager.proxy.UnbindNode(context.Background(), uid); err != nil {
				log.Errorf("unbind node failed: %v", err)
			}

			return nil
		})
	}

	_ = eg.Wait()

	clear(r.members)

	r.manager.rooms.Delete(r.id)

	if actor, ok := r.manager.proxy.Actor(roomActor, xconv.String(r.id)); ok {
		actor.Destroy()
	}

	return nil
}

// 进入房间
func (r *Room) doEnterRoom(uid int64) (err error) {
	user, err := r.manager.doLoadUser(uid)
	if err != nil {
		return err
	}

	r.rw.Lock()

	if r.state.Load() != stateNormal {
		r.rw.Unlock()
		return errors.NewError(code.IllegalOperation)
	}

	if err = user.doSaveRoom(r); err != nil {
		r.rw.Unlock()
		return err
	}

	r.members[uid] = &Member{user: user, time: xtime.Now()}

	r.rw.Unlock()

	defer func() {
		if err != nil {
			r.rw.Lock()
			defer r.rw.Unlock()

			if r.state.Load() != stateNormal {
				return
			}

			user.doClearRoom()

			delete(r.members, uid)
		}
	}()

	if err = r.manager.proxy.BindActor(uid, roomActor, xconv.String(r.id)); err != nil {
		log.Errorf("bind actor failed: %v", err)
		return errors.NewError(err, code.InternalError)
	}

	defer func() {
		if err != nil {
			r.manager.proxy.UnbindActor(uid, roomActor)
		}
	}()

	if err = r.manager.proxy.BindNode(context.Background(), uid); err != nil {
		log.Errorf("bind node failed: %v", err)
		return errors.NewError(err, code.InternalError)
	}

	return
}

// 离开房间
func (r *Room) doLeaveRoom(uid int64) error {
	r.rw.Lock()

	if r.state.Load() != stateNormal {
		r.rw.Unlock()
		return errors.NewError(code.IllegalOperation)
	}

	member, ok := r.members[uid]
	if !ok || member == r.creator {
		r.rw.Unlock()
		return errors.NewError(code.IllegalOperation)
	}

	delete(r.members, uid)

	member.user.doClearRoom()

	r.manager.proxy.UnbindActor(uid, roomActor)

	r.rw.Unlock()

	if err := r.manager.proxy.UnbindNode(context.Background(), uid); err != nil {
		log.Errorf("unbind node failed: %v", err)
	}

	return nil
}

// 拉取成员列表
func (r *Room) fetchMembers() ([]*MemberInfo, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	if r.state.Load() != stateNormal {
		return nil, errors.NewError(code.IllegalOperation)
	}

	list := make([]*MemberInfo, 0, len(r.members))
	for _, member := range r.members {
		list = append(list, member.doMakeMemberInfo())
	}

	return list, nil
}

// 构建聊天室信息
func (r *Room) doMakeRoomInfo() *RoomInfo {
	return &RoomInfo{
		ID:      r.id,
		Name:    r.name,
		Creator: r.creator.doMakeMemberInfo(),
	}
}
