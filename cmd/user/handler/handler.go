package handler

import (
	"context"
	"github.com/SIN5t/tiktok_v2/cmd/user/dal"
	"github.com/SIN5t/tiktok_v2/cmd/user/service"
	config "github.com/SIN5t/tiktok_v2/config/const"
	user "github.com/SIN5t/tiktok_v2/kitex_gen/user"
	"github.com/SIN5t/tiktok_v2/pkg/jwt/utils"
	"github.com/SIN5t/tiktok_v2/pkg/snowflakes"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {

	//雪花算法创建id
	userViper := viper.Init("user")

	userId := snowflakes.GenerateSnowFlakeId(int64(userViper.GetInt("snowflake.node")))
	userName := req.GetUsername()
	password := req.GetPassword()

	if exist := service.NameExist(userName); exist {
		// 返回用户名存在
		return &user.UserRegisterResponse{
			StatusCode: config.FailResponse,
			StatusMsg:  "用户名重复，请重新输入",
		}, nil
	}
	// 加密password
	encryptPwd := utils.MD5(password)

	createUser := &db.UserReg{
		UserId:   userId,
		Username: userName,
		Password: encryptPwd,
	}
	err = service.SaveUserToMysql(createUser)
	if err != nil {
		return &user.UserRegisterResponse{
			StatusCode: config.FailResponse,
			StatusMsg:  "创建用户失败",
		}, err
	}

	return &user.UserRegisterResponse{
		StatusCode: config.Success,
		StatusMsg:  "成功创建账户",
		UserId:     userId,
	}, err
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	userBasicInfo, err := service.CheckUser(req.GetUsername(), req.GetPassword())
	if err == nil {
		return &user.UserLoginResponse{
			UserId: userBasicInfo.UserId,
		}, nil
	}
	return nil, err
}

// UserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	// TODO: Your code here...
	return
}
