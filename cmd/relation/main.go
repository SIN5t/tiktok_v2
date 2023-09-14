package main

import (
	relation2 "github.com/SIN5t/tiktok_v2/cmd/relation/service"
	relation "github.com/SIN5t/tiktok_v2/kitex_gen/relation/relationservice"
	"log"
)

func main() {
	svr := relation.NewServer(new(relation2.RelationServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
