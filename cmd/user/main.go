package main

import (
	"fmt"
	db "github.com/SIN5t/tiktok_v2/cmd/user/dal"
	"github.com/SIN5t/tiktok_v2/cmd/user/handler"
	"github.com/SIN5t/tiktok_v2/kitex_gen/user/userservice"
	"github.com/SIN5t/tiktok_v2/pkg/etcd"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"net"
)

var (
	config          = viper.Init("user")
	serverName      = config.GetString("server.name")
	etcdEndpoin     = fmt.Sprintf("%s:%d", config.GetString("etcd.host"), config.GetInt("etcd.port"))
	otlpRpcReceiver = fmt.Sprintf("%s:%d", config.GetString("otel.host"), config.GetInt("otel.port"))
	serviceAddr     = fmt.Sprintf("%s:%d", config.GetString("server.host"), config.GetInt("server.port"))
)

func main() {
	//初始化数据库
	db.InitMysqlDB()
	db.InitRdb()

	// 服务注册
	registry, err := etcd.NewEtcdRegistry([]string{etcdEndpoin})
	if err != nil {
		klog.Fatal(err.Error())
	}

	//使用net.ResolveTCPAddr函数将服务地址（serviceAddr）解析为TCP地址（*net.TCPAddr）。
	addr, err := net.ResolveTCPAddr("tcp", serviceAddr)
	if err != nil {
		klog.Fatal(err.Error())
	}

	// openTelemetry 链路追踪
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serverName),
		provider.WithExportEndpoint(otlpRpcReceiver),
		provider.WithInsecure(),
	)

	userServer := userservice.NewServer(
		//&handler.UserServiceImpl{},
		new(handler.UserServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(registry), //注册中心
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serverName}),
		server.WithSuite(tracing.NewServerSuite()), //链路追踪
	)
	if err = userServer.Run(); err != nil {
		klog.Fatalf("%v stopped with error: %v", serverName, err.Error())
	}
}
