package service

import (
	db "github.com/SIN5t/tiktok_v2/cmd/video/dal"
	"github.com/SIN5t/tiktok_v2/kitex_gen/video"
	"golang.org/x/net/context"
	"time"
)

func FeedService(ctx context.Context, latestTime int64, token string) (err error, videoListRes []*video.Video, nextTime int64) {

	videoListRaw, err := db.GetVideosByLastTime(ctx, latestTime, 30)

	flag := false
	if token != "" {
		//TODO 校验 jwt
		flag = true //已登入
	}

	videoListRes = make([]*video.Video, 0) //存的元素类型是指针，而不是对象本身
	for _, value := range videoListRaw {

		//如果登入
		if flag {
			//TODO 当前视频是否点赞

			// TODO 当前用户是否关注视频作者

		}

		// TODO 作者信息完善，包括点赞数、关注数，名字，是否关注等
		VideoAuthor := &video.User{Id: value.AuthorID}

		curVideo := video.Video{
			Id:            value.ID,
			Author:        VideoAuthor,
			PlayUrl:       value.PlayUrl,
			CoverUrl:      value.CoverUrl,
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			Title:         value.Title,
		}
		videoListRes = append(videoListRes, &curVideo) //存入地址
	}
	if len(videoListRes) <= 0 {
		return nil, videoListRes, time.Now().UnixMilli()
	}
	/*time1 := videoListRaw[len(videoListRaw)].CreateTime.UnixMilli()
	fmt.Println(time1)*/
	return nil, videoListRes, videoListRaw[len(videoListRaw)-1].CreateTime.UnixMilli()
}
