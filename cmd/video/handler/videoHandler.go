package handler

import (
	"context"
	"github.com/SIN5t/tiktok_v2/cmd/video/service"
	config "github.com/SIN5t/tiktok_v2/config/const"
	"github.com/SIN5t/tiktok_v2/kitex_gen/video"
	KafkaVideo "github.com/SIN5t/tiktok_v2/pkg/kafka/video"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/google/uuid"
	"github.com/prometheus/common/log"
	"os"
	"strings"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// Feed implements the VideoServiceImpl interface.
// 要求按照投稿时间倒叙输出
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	latestTime := req.LatestTime
	token := req.Token

	// 获取视频
	err, videoListRes, nextTime := service.FeedService(ctx, latestTime, token)
	if err != nil {
		log.Error(err)
		response := &video.FeedResponse{
			StatusCode: 1,
			StatusMsg:  "视频查询失败",
		}
		return response, nil
	}

	return &video.FeedResponse{
		StatusCode: 0,
		StatusMsg:  "刷新成功",
		VideoList:  videoListRes,
		NextTime:   nextTime - 1, //注意下次的时间是这次的最后
	}, nil
}

// PublishAction implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {

	videoViper := viper.Init("video")
	tempVideoLoc := videoViper.GetString("location.video")

	//将请求中的视频（[]byte）拿出来,保存到temp目录下，之后交给消息队列处理，上传oss等
	fileName := strings.Replace(uuid.New().String(), "-", "", -1) + ".mp4" //为视频生成唯一的视频名称
	filePath := tempVideoLoc + fileName                                    // 存在 ./temp/video/临时目录下
	err = os.WriteFile(filePath, req.GetData(), config.FileAuth)           //路径不存在自动创建
	if err != nil {
		resp = &video.PublishActionResponse{
			StatusCode: config.FailResponse,
			StatusMsg:  "视频保存出现错误",
		}
		return resp, err
	}

	// 封装消息，调用消息队列，发布消息上传的任务,注意也要发送当前作者的id
	videoMsg := KafkaVideo.VideoMsg{
		VideoPath: filePath, //文件路径包含文件名
		VideoName: fileName,
		AuthorId:  req.UserId,
		Title:     req.Title,
	}

	// 生成之前需要开启消费者,消费者会阻塞，需要开一个goroutine,在main中开启
	KafkaVideo.ProduceMsg(KafkaVideo.NewProducer(), videoMsg)

	resp = &video.PublishActionResponse{
		StatusCode: config.Success,
		StatusMsg:  "视频上传中...",
	}
	return
}

// PublishList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	// TODO: Your code here...
	return
}
