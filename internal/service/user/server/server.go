package server

import (
	"context"
	"github.com/dobyte/due-chat/internal/code"
	jwtcomp "github.com/dobyte/due-chat/internal/component/jwt"
	mysqlcomp "github.com/dobyte/due-chat/internal/component/mysql"
	userdao "github.com/dobyte/due-chat/internal/dao/user"
	"github.com/dobyte/due-chat/internal/model"
	"github.com/dobyte/due-chat/internal/service/user/pb"
	"github.com/dobyte/due/v2/cluster/mesh"
	"github.com/dobyte/due/v2/errors"
	"github.com/dobyte/due/v2/log"
	"github.com/dobyte/due/v2/utils/xconv"
	"github.com/dobyte/due/v2/utils/xtime"
	"github.com/dobyte/jwt"
	"golang.org/x/crypto/bcrypt"
)

var _ pb.UserServer = &Server{}

type Server struct {
	pb.UnimplementedUserServer
	proxy   *mesh.Proxy
	jwt     *jwtcomp.JWT
	userDao *userdao.User
}

const defaultGate = "ws://127.0.0.1:3533"

func NewServer(proxy *mesh.Proxy) *Server {
	return &Server{
		proxy:   proxy,
		jwt:     jwtcomp.Instance(),
		userDao: userdao.NewUser(mysqlcomp.Instance()),
	}
}

func (s *Server) Init() {
	s.proxy.AddServiceProvider("user", &pb.User_ServiceDesc, s)
}

// Register 注册
func (s *Server) Register(ctx context.Context, args *pb.RegisterArgs) (*pb.RegisterReply, error) {
	user, err := s.doQueryUserByAccount(ctx, args.Account)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, errors.NewError(code.AccountExists)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("generate password failed:%v", err)
		return nil, errors.NewError(err, code.InternalError)
	}

	_, err = s.userDao.Insert(ctx, &model.User{
		Account:     args.Account,
		Password:    xconv.String(password),
		Nickname:    args.Nickname,
		RegisterAt:  xtime.Now(),
		RegisterIP:  args.ClientIP,
		LastLoginAt: xtime.Now(),
		LastLoginIP: args.ClientIP,
	})
	if err != nil {
		log.Errorf("insert user failed: %v", err)
		return nil, errors.NewError(err, code.InternalError)
	}

	return &pb.RegisterReply{}, nil
}

// Login 登录
func (s *Server) Login(ctx context.Context, args *pb.LoginArgs) (*pb.LoginReply, error) {
	user, err := s.doQueryUserByAccount(ctx, args.Account)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.NewError(code.NotFoundUser)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(args.Password)); err != nil {
		return nil, errors.NewError(err, code.WrongAccountOrPassword)
	}

	s.doUpdateUserLastLoginInfo(ctx, user.ID, args.ClientIP)

	token, err := s.jwt.GenerateToken(jwt.Payload{
		s.jwt.IdentityKey(): user.ID,
	})
	if err != nil {
		log.Errorf("generate token failed, id = %d err = %v", user.ID, err)
		return nil, errors.NewError(err, code.InternalError)
	}

	return &pb.LoginReply{
		Gate:  defaultGate,
		Token: token.Token,
	}, nil
}

// ValidateToken 验证Token
func (s *Server) ValidateToken(ctx context.Context, args *pb.ValidateTokenArgs) (*pb.ValidateTokenReply, error) {
	identity, err := s.jwt.ExtractIdentity(args.Token)
	if err != nil {
		return nil, errors.NewError(err, code.Unauthorized)
	}

	uid := xconv.Int64(identity)
	if uid <= 0 {
		return nil, errors.NewError(err, code.Unauthorized)
	}

	return &pb.ValidateTokenReply{UID: uid}, nil
}

// FetchUser 拉取用户
func (s *Server) FetchUser(ctx context.Context, args *pb.FetchUserArgs) (*pb.FetchUserReply, error) {
	user, err := s.doQueryUserByUID(ctx, args.UID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.NewError(code.NotFoundUser)
	}

	return &pb.FetchUserReply{User: &pb.UserInfo{
		UID:         user.ID,
		Account:     user.Account,
		Nickname:    user.Nickname,
		RegisterAt:  user.RegisterAt.Format(xtime.DateTime),
		RegisterIP:  user.RegisterIP,
		LastLoginAt: user.LastLoginAt.Format(xtime.DateTime),
		LastLoginIP: user.LastLoginIP,
	}}, nil
}

// 根据用户ID查询用户信息
func (s *Server) doQueryUserByUID(ctx context.Context, uid int64) (*model.User, error) {
	user, err := s.userDao.FindOne(ctx, func(cols *userdao.Columns) interface{} {
		return map[string]interface{}{
			cols.ID: uid,
		}
	})
	if err != nil {
		log.Errorf("find user failed, uid = %s err = %v", uid, err)
		return nil, errors.NewError(err, code.InternalError)
	}

	return user, nil
}

// 根据账号查询用户信息
func (s *Server) doQueryUserByAccount(ctx context.Context, account string) (*model.User, error) {
	user, err := s.userDao.FindOne(ctx, func(cols *userdao.Columns) interface{} {
		return map[string]interface{}{
			cols.Account: account,
		}
	})
	if err != nil {
		log.Errorf("find user failed, account = %s err = %v", account, err)
		return nil, errors.NewError(err, code.InternalError)
	}

	return user, nil
}

// 更新用户最新登录信息
func (s *Server) doUpdateUserLastLoginInfo(ctx context.Context, uid int64, clientIP string) {
	_, err := s.userDao.Update(ctx, func(cols *userdao.Columns) interface{} {
		return map[string]interface{}{
			cols.ID: uid,
		}
	}, func(cols *userdao.Columns) interface{} {
		return map[string]interface{}{
			cols.LastLoginAt: xtime.Now(),
			cols.LastLoginIP: clientIP,
		}
	})
	if err != nil {
		log.Errorf("update login info failed: %v", err)
	}
}
