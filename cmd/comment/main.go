package main

import (
	"github.com/SIN5t/tiktok_v2/cmd/comment/service"
	comment "github.com/SIN5t/tiktok_v2/kitex_gen/comment/commentservice"
	"log"
)

func main() {
	svr := comment.NewServer(new(service.CommentServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
