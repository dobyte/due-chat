package logic

import "github.com/dobyte/due/v2/utils/xtime"

type Member struct {
	user *User      // 用户信息
	time xtime.Time // 加入时间
}

// 构建成员信息
func (m *Member) doMakeMemberInfo() *MemberInfo {
	return &MemberInfo{
		User:      m.user.doMakeUserInfo(),
		Timestamp: m.time.Unix(),
	}
}
