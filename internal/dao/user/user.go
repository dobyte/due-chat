package user

import (
	"github.com/dobyte/due-chat/internal/dao/user/internal"
	"gorm.io/gorm"
)

type (
	Columns = internal.Columns
	OrderBy = internal.OrderBy
	FilterFunc = internal.FilterFunc
	UpdateFunc = internal.UpdateFunc
	ColumnFunc = internal.ColumnFunc
	OrderFunc = internal.OrderFunc
)

type User struct {
	*internal.User
}

func NewUser(db *gorm.DB) *User {
	return &User{User: internal.NewUser(db)}
}
