package test

import (
	db "github.com/SIN5t/tiktok_v2/cmd/user/dal"
	"github.com/SIN5t/tiktok_v2/cmd/user/service"
	"testing"
)

func TestGorm(t *testing.T) {
	service.CheckUser("test01", "123456")
}

func TestSaveUser(t *testing.T) {
	user := db.UserReg{
		Username: "test03",
		Password: "123456",
		UserId:   132456798456,
	}
	service.SaveUserToMysql(&user)
}
