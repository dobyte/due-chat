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
	"github.com/dobyte/due/v2/utils/xrand"
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
		Token: token.Token,
	}, nil
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

// 分配网关
func (s *Server) doAssignGate() string {
	if opts := server.GetOpts(); len(opts.Gates) > 0 {
		return opts.Gates[xrand.Int(0, len(opts.Gates)-1)]
	} else {
		return ""
	}
}
