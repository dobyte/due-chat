package route

const (
	ValidateToken int32 = 1000 // 校验Token
	CreateRoom    int32 = 1001 // 创建房间
	DismissRoom   int32 = 1002 // 解散房间
	EnterRoom     int32 = 1003 // 进入房间
	LeaveRoom     int32 = 1004 // 离开房间
	SendMessage   int32 = 1005 // 发送消息
	FetchMembers  int32 = 1006 // 拉取成员
	MessageNotify int32 = 1007 // 消息通知
)
