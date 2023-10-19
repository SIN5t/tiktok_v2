package service

import (
	"errors"
	db "github.com/SIN5t/tiktok_v2/cmd/user/dal"
	"github.com/SIN5t/tiktok_v2/pkg/jwt/utils"
)

// CheckUser 数据库验证用户名密码是否正确
func CheckUser(username string, password string) (*db.UserReg, error) {
	//user := &db.UserReg{}
	var user db.UserReg
	err := db.UserMysqlDB.
		Select("user_id").
		Where("username = ? and password = ?", username, utils.MD5(password)).
		First(&user).
		Error
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	return &user, nil
}
