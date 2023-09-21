package main

import (
	"fmt"
	db "github.com/SIN5t/tiktok_v2/cmd/video/dal"
	"github.com/SIN5t/tiktok_v2/cmd/video/handler"
	"github.com/SIN5t/tiktok_v2/kitex_gen/video/videoservice"
	"github.com/SIN5t/tiktok_v2/pkg/etcd"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

var (
	config      = viper.Init("video")
	serviceName = config.GetString("server.name") //server.WithServerBasicInfo用到
	serviceAddr = fmt.Sprintf("%s:%d", config.GetString("server.host"), config.GetInt("server.port"))
	etcdAddr    = fmt.Sprintf("%s:%d", config.GetString("etcd.host"), config.GetInt("etcd.port"))
	//signingKey  = config.Viper.GetString("JWT.signingKey")
)

func main() {
	db.InitDB()
	//服务注册
	registry, err := etcd.NewEtcdRegistry([]string{etcdAddr})
	if err != nil {
		log.Fatal(err.Error())
	}
	//使用net.ResolveTCPAddr函数将服务地址（serviceAddr）解析为TCP地址（*net.TCPAddr）。
	addr, err := net.ResolveTCPAddr("tcp", serviceAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	// TODO obs 链路追踪
	/*provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serviceName),
		provider.WithExportEndpoint(etcdAddr),
		provider.WithInsecure(),
	)*/

	videoServer := videoservice.NewServer(
		new(handler.VideoServiceImpl), //这个service就是mvc中的service
		server.WithServiceAddr(addr),  //tcp 的地址
		server.WithRegistry(registry), //注册中心
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)

	if err := videoServer.Run(); err != nil {
		klog.Fatalf("%v stopped with error: %v", serviceName, err.Error())
	}
}
