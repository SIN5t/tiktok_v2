package db

import (
	"context"
	"encoding/json"
	config "github.com/SIN5t/tiktok_v2/config/const"
	"github.com/redis/go-redis/v9"
	"strconv"
)

// RedisAddVideo sortedSet(ZSet),key:videos,val{{Score:time,Member:序列化对象},{...},...}
func RedisAddVideo(video *Video) {

	videoMarsh, _ := json.Marshal(video)
	// 存入redis
	VideoRdb.ZAdd(
		context.Background(),
		config.RedisZsetVideoKey,
		redis.Z{
			Score:  float64(video.CreateTime.Unix()),
			Member: videoMarsh,
		},
	)

}

// RedisGetVideos 存的值是序列化对象
func RedisGetVideos(latestTime int64) ([]*Video, error) {
	latestTimeStr := strconv.FormatInt(latestTime, 10)
	// 需要按照给定的分数查询分数的前30个视频，返回序列化的视频切片
	videosMarshList, err := VideoRdb.ZRangeArgs(
		context.Background(),
		redis.ZRangeArgs{
			Key:     config.RedisZsetVideoKey,
			Start:   0,
			Stop:    "(" + latestTimeStr, // 括号代表开区间
			ByScore: true,
			ByLex:   false,
			Rev:     true,
			Offset:  0,
			Count:   10,
		},
	).Result()
	if err != nil {
		return nil, err
	}

	videoList := make([]*Video, 0)
	//反序列化返回
	for _, videoJsonStr := range videosMarshList {
		var videoUnmarshal Video
		err := json.Unmarshal([]byte(videoJsonStr), &videoUnmarshal)
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, &videoUnmarshal)
	}
	return videoList, nil
}
