package api

type LoginReq struct {
	Account  string `json:"account" validate:"required"`  // 账号
	Password string `json:"password" validate:"required"` // 密码
}

type LoginRes struct {
	Gate  string `json:"gate"`  // 网关
	Token string `json:"token"` // Token
}

type RegisterReq struct {
	Account  string `json:"account" validate:"required"`  // 账号
	Nickname string `json:"nickname" validate:"required"` // 昵称
	Password string `json:"password" validate:"required"` // 密码
}
