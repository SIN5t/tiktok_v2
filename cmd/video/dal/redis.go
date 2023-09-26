package db

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

// RedisAddVideos sortedSet(ZSet),key:videos,val{{Score:time,Member:序列化对象},{...},...}
func RedisAddVideos(videos []*Video) {

}

// RedisGetVideos 存的值是序列化对象
func RedisGetVideos(time string) ([]*Video, error) {
	// 需要按照分数查询分数的前30个视频，返回序列化的视频切片
	videosMarshList, err := VideoRdb.ZRangeArgs(
		context.Background(),
		redis.ZRangeArgs{
			Key:     "RedisVideo",
			Start:   0,
			Stop:    "(" + time, // 括号代表开区间
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
		var videoUnmarshal *Video
		err := json.Unmarshal([]byte(videoJsonStr), &videoUnmarshal)
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, videoUnmarshal)
	}
	return videoList, nil
}
