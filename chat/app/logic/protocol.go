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
	RoomID uint64 `json:"roomID"` // 房间ID
}

type EnterRoomReq struct {
	RoomID uint64 `json:"roomID"` // 房间ID
}

type EnterRoomRes struct {
	Code int               `json:"code"` // 响应码
	Data *EnterRoomResData `json:"data"` // 响应数据
}

type EnterRoomResData struct {
	ID   uint64 `json:"id"`   // 聊天室ID
	Name string `json:"name"` // 聊天室名称
}

type LeaveRoomRes struct {
	Code int `json:"code"` // 响应码
}

type SendMessageReq struct {
	Content string `json:"content"` // 发送内容
}

type SendMessageRes struct {
	Code int `json:"code"` // 响应码
}

type MessageNotify struct {
	Sender    *MemberInfo `json:"sender"`    // 消息发送者
	Content   string      `json:"content"`   // 消息内容
	Timestamp int64       `json:"timestamp"` // 发送时间
}

type MemberInfo struct {
	UID      int64  `json:"uid"`      // 用户ID
	Nickname string `json:"nickname"` // 用户昵称
}
