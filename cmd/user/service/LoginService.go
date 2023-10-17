package service

import (
	db "github.com/SIN5t/tiktok_v2/cmd/user/dal"
)

// CheckUser 数据库验证用户名密码是否正确
func CheckUser(username string, password string) (*db.UserReg, error) {
	user := &db.UserReg{}
	if err := db.UserMysqlDB.Model("user_reg").Where("Username = ? and Password = ?").First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
