package main

import (
	message2 "github.com/SIN5t/tiktok_v2/cmd/message/service"
	message "github.com/SIN5t/tiktok_v2/kitex_gen/message/messageservice"
	"log"
)

func main() {
	svr := message.NewServer(new(message2.MessageServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
