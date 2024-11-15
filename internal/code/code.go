package code

import "github.com/dobyte/due/v2/codes"

var (
	OK                     = codes.OK
	Canceled               = codes.Canceled
	Unknown                = codes.Unknown
	InvalidArgument        = codes.InvalidArgument
	DeadlineExceeded       = codes.DeadlineExceeded
	NotFound               = codes.NotFound
	InternalError          = codes.InternalError
	Unauthorized           = codes.Unauthorized
	IllegalInvoke          = codes.IllegalInvoke
	IllegalRequest         = codes.IllegalRequest
	NotFoundUser           = codes.NewCode(100, "not found user")
	WrongAccountOrPassword = codes.NewCode(101, "wrong account or password")
	AccountExists          = codes.NewCode(102, "account exists")
	IllegalOperation       = codes.NewCode(103, "illegal operation")
)
