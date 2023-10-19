package service

import (
	db "github.com/SIN5t/tiktok_v2/cmd/user/dal"
)

func SaveUserToMysql(user *db.UserReg) error {

	//将该用户存入mysql数据库
	return db.UserMysqlDB.Create(user).Error
}

// NameExist 如果用户名存在，返回true，否则返回false
func NameExist(userName string) bool {
	///查询是否有重复
	if db.UserMysqlDB.Where("username = ?", userName).First(&db.UserReg{}).RowsAffected > 0 {
		return true
	}
	return false
}
