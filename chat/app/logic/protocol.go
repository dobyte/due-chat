package logic

type ValidateTokenReq struct {
	Token string `json:"token"` // 登录Token
}

type ValidateTokenRes struct {
	Code int `json:"code"` // 响应码
}

type CreateRoomReq struct {
	Name string `json:"name"` // 聊天室名称
}

type CreateRoomRes struct {
	Code int                `json:"code"` // 响应码
	Data *CreateRoomResData `json:"data"` // 响应数据
}

type CreateRoomResData struct {
	Room *RoomInfo `json:"room"` // 聊天室
}

type DismissRoomRes struct {
	Code int `json:"code"` // 响应码
}

type EnterRoomReq struct {
	RoomID uint64 `json:"roomID"` // 聊天室ID
}

type EnterRoomRes struct {
	Code int               `json:"code"` // 响应码
	Data *EnterRoomResData `json:"data"` // 响应数据
}

type EnterRoomResData struct {
	Room *RoomInfo `json:"room"` // 聊天室
}

type LeaveRoomRes struct {
	Code int `json:"code"` // 响应码
}

type FetchMembersRes struct {
	Code int                  `json:"code"` // 响应码
	Data *FetchMembersResData `json:"data,omitempty"`
}

type FetchMembersResData struct {
	List []*MemberInfo `json:"list"` // 成员列表
}

type SendMessageReq struct {
	Content string `json:"content"` // 发送内容
}

type SendMessageRes struct {
	Code int `json:"code"` // 响应码
}

type MessageNotify struct {
	Sender    *UserInfo `json:"sender"`    // 消息发送者
	Content   string    `json:"content"`   // 消息内容
	Timestamp int64     `json:"timestamp"` // 发送时间
}

type MemberInfo struct {
	User      *UserInfo `json:"user"`      // 用户信息
	Timestamp int64     `json:"timestamp"` // 加入时间
}

type UserInfo struct {
	UID      int64  `json:"uid"`      // 用户ID
	Nickname string `json:"nickname"` // 用户昵称
}

type RoomInfo struct {
	ID      int64       `json:"id"`      // 聊天室ID
	Name    string      `json:"name"`    // 聊天室名称
	Creator *MemberInfo `json:"creator"` // 聊天室创建者
}
