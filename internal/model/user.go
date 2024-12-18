package model

import "time"

//go:generate gorm-dao-generator -model-dir=. -model-names=User -dao-dir=../dao/ -sub-pkg-enable=true
type User struct {
	ID          int64     `gorm:"column:id"`            // ID
	Account     string    `gorm:"column:account"`       // 账号
	Password    string    `gorm:"column:password"`      // 密码
	Nickname    string    `gorm:"column:nickname"`      // 昵称
	RegisterAt  time.Time `gorm:"column:register_at"`   // 注册时间
	RegisterIP  string    `gorm:"column:register_ip"`   // 注册IP
	LastLoginAt time.Time `gorm:"column:last_login_at"` // 最新登录时间
	LastLoginIP string    `gorm:"column:last_login_ip"` // 最新登录IP
}
