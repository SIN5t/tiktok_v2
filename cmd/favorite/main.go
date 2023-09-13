package main

import (
	favorite2 "github.com/SIN5t/tiktok_v2/cmd/favorite/service"
	favorite "github.com/SIN5t/tiktok_v2/kitex_gen/favorite/favoriteservice"
	"log"
)

func main() {
	svr := favorite.NewServer(new(favorite2.FavoriteServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
