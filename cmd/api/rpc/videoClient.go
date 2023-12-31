package rpc

import (
	"context"
	"fmt"
	"github.com/SIN5t/tiktok_v2/kitex_gen/video"
	"github.com/SIN5t/tiktok_v2/kitex_gen/video/videoservice"
	"github.com/SIN5t/tiktok_v2/pkg/etcd"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"log"
	"time"
)

//返回handler构建出的一个client
//生成的client包含了用于与服务端通信的代码，可以作为客户端调用远程方法。
//而生成的service包含了实现具体服务逻辑的代码，可以作为服务端接收并处理请求。

var videoClient videoservice.Client //需要配置，需要初始化

// InitVideoClient
// 在构建client和server 通信之前，先关注连通性
func InitVideoClient() {
	//加载配置文件
	videoConfig := viper.Init("video")
	etcdAddr := fmt.Sprintf("%s:%d", videoConfig.GetString("etcd.host"), videoConfig.GetInt("etcd.port"))
	otlpRpcReceiver := fmt.Sprintf("%s:%d", videoConfig.GetString("otel.host"), videoConfig.GetInt("otel.port"))
	serverName := videoConfig.GetString("server.name") //指定客户端所连接的服务的名称

	resolver, err := etcd.NewEtcdResolver([]string{etcdAddr})
	if err != nil {
		hlog.Fatal(err)
	}

	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serverName),
		provider.WithExportEndpoint(otlpRpcReceiver),
		provider.WithInsecure(),
	)

	newClient, err := videoservice.NewClient(
		serverName,
		//client.WithMiddleware(middleware.CommonMiddleware),
		//client.WithInstanceMW(middleware.ClientMiddleware), //实例级别的中间件。这个中间件可能会对客户端的每个请求进行预处理或后处理
		client.WithMuxConnection(1),                        // mux
		client.WithRPCTimeout(300*time.Second),             // rpc timeout
		client.WithConnectTimeout(300000*time.Millisecond), // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()),  // retry
		client.WithSuite(tracing.NewClientSuite()),         // tracer 添加了一个追踪器，用于跟踪客户端的请求和响应
		client.WithResolver(resolver),                      // resolver
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serverName}),
	)
	if err != nil {
		log.Fatal(err)
	}
	videoClient = newClient
}

// FeedClient 给client前端数据，让client与server通信，拿到server返回的结果
func FeedClient(ctx context.Context, req *video.FeedRequest) (*video.FeedResponse, error) {
	return videoClient.Feed(ctx, req)
}

// PublishAction 将client调用server的方法进行rpc通信
func PublishAction(ctx context.Context, req *video.PublishActionRequest) (*video.PublishActionResponse, error) {
	return videoClient.PublishAction(ctx, req)
}
