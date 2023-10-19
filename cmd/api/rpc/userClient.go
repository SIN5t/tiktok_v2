package rpc

import (
	"fmt"
	"github.com/SIN5t/tiktok_v2/kitex_gen/user"
	"github.com/SIN5t/tiktok_v2/kitex_gen/user/userservice"
	"github.com/SIN5t/tiktok_v2/pkg/etcd"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"golang.org/x/net/context"
	"time"
)

var userClient userservice.Client

func InitUserClient() {
	userViperConfig := viper.Init("user")
	serverName := userViperConfig.GetString("server.name")
	etcdEndPoint := fmt.Sprintf("%s:%d", userViperConfig.GetString("etcd.host"), userViperConfig.GetInt("etcd.port"))
	otel := fmt.Sprintf("%s:%d", userViperConfig.GetString("otel.host"), userViperConfig.GetInt("otel.port"))

	etcdResolver, err := etcd.NewEtcdResolver([]string{etcdEndPoint})
	if err != nil {
		klog.Fatal(err.Error())
	}
	provider.NewOpenTelemetryProvider(
		provider.WithInsecure(),
		provider.WithExportEndpoint(otel),
		provider.WithServiceName(serverName),
	)

	newUserClient, err := userservice.NewClient(
		serverName,
		client.WithMuxConnection(1),
		client.WithRPCTimeout(time.Second*15),
		client.WithConnectTimeout(time.Second*15),
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(tracing.NewClientSuite()),        // tracer 添加了一个追踪器，用于跟踪客户端的请求和响应
		client.WithResolver(etcdResolver),                 // resolver
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serverName}),
	)
	if err != nil {
		klog.Fatal(err.Error())
	}
	userClient = newUserClient
}

func Login(ctx context.Context, request *user.UserLoginRequest) (*user.UserLoginResponse, error) {
	resp, err := userClient.Login(ctx, request)
	return resp, err
}

func Register(ctx context.Context, req *user.UserRegisterRequest) (*user.UserRegisterResponse, error) {
	resp, err := userClient.Register(ctx, req)
	return resp, err

}
