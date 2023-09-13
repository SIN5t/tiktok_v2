package handler

import (
	"context"
	"github.com/SIN5t/tiktok_v2/cmd/video/service"
	video "github.com/SIN5t/tiktok_v2/kitex_gen/video"
	"github.com/prometheus/common/log"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// Feed implements the VideoServiceImpl interface.
// 要求按照投稿时间倒叙输出
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	latestTime := req.LatestTime
	token := req.Token
	//数据库获取视频
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
	// TODO: Your code here...
	return
}

// PublishList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	// TODO: Your code here...
	return
}
