package main

import (
	user2 "github.com/SIN5t/tiktok_v2/cmd/user/service"
	user "github.com/SIN5t/tiktok_v2/kitex_gen/user/userservice"
	"log"
)

func main() {
	svr := user.NewServer(new(user2.UserServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
