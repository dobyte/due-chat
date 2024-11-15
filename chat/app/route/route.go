package route

const (
	ValidateToken int32 = 1000 // 校验Token
	FetchProfile  int32 = 1001 // 拉取玩家信息
	CreateRoom    int32 = 1002 // 创建房间
	DismissRoom   int32 = 1003 // 解散房间
	EnterRoom     int32 = 1004 // 进入房间
	LeaveRoom     int32 = 1005 // 离开房间
	SendMessage   int32 = 1006 // 发送消息
	FetchMembers  int32 = 1007 // 拉取成员
	MessageNotify int32 = 1008 // 消息通知
)
