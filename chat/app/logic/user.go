package logic

import (
	"context"
	"github.com/dobyte/due-chat/internal/code"
	usersvc "github.com/dobyte/due-chat/internal/service/user/client"
	userpb "github.com/dobyte/due-chat/internal/service/user/pb"
	"github.com/dobyte/due/v2/errors"
	"sync"
)

type User struct {
	uid      int64        // 用户ID
	account  string       // 用户账号
	nickname string       // 用户昵称
	rw       sync.RWMutex // 锁
	room     *Room        // 用户所在房间
}

func newUser(mgr *Manager, uid int64) (*User, error) {
	client, err := usersvc.NewClient(mgr.proxy.NewMeshClient)
	if err != nil {
		return nil, errors.NewError(err, code.InternalError)
	}

	reply, err := client.FetchUser(context.Background(), &userpb.FetchUserArgs{UID: uid})
	if err != nil {
		return nil, err
	}

	u := &User{}
	u.uid = reply.User.UID
	u.account = reply.User.Account
	u.nickname = reply.User.Nickname

	return u, nil
}

// 获取用户所在房间
func (u *User) doGetRoom() *Room {
	u.rw.RLock()
	defer u.rw.RUnlock()

	return u.room
}

// 保存用户房间
func (u *User) doSaveRoom(room *Room) error {
	u.rw.Lock()
	defer u.rw.Unlock()

	if u.room != nil {
		return errors.NewError(code.IllegalOperation)
	}

	u.room = room

	return nil
}

// 清理用户房间
func (u *User) doClearRoom() (room *Room) {
	u.rw.Lock()
	room = u.room
	u.room = nil
	u.rw.Unlock()

	return
}

// 检测用户是否有已进入的聊天室
func (u *User) doHasRoom() bool {
	u.rw.RLock()
	defer u.rw.RUnlock()

	return u.room != nil
}

// 构建用户信息
func (u *User) doMakeUserInfo() *UserInfo {
	return &UserInfo{
		UID:      u.uid,
		Nickname: u.nickname,
	}
}
