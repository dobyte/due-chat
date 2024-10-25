package api

type LoginReq struct {
	Account  string `json:"account" validate:"required"`  // 账号
	Password string `json:"password" validate:"required"` // 密码
}

type LoginResData struct {
	Gate  string `json:"gate"`  // 网关
	Token string `json:"token"` // Token
}

type RegisterReq struct {
	Account  string `json:"account" validate:"required"`  // 账号
	Password string `json:"password" validate:"required"` // 密码
}
