// Code generated by hertz generator.

package main

import (
	"fmt"
	"github.com/SIN5t/tiktok_v2/cmd/api/initialize"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/hertz/pkg/app/server"
)

var config = viper.Init("api")
var hostPorts = fmt.Sprintf("%s:%d", config.GetString("server.host"), config.GetInt("server.port"))

func main() {
	//初始化client、中间件、jwt、logger...
	initialize.Init()

	// TODO 链路追踪

	h := server.New(
		server.WithHostPorts(hostPorts),
	)

	register(h)
	h.Spin()
}
