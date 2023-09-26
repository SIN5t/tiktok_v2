package db

import "context"

// RedisAddVideos sortedSet(ZSet),key:videos,val{{Score:time,Member:序列化对象},{...},...}
func RedisAddVideos(videos []*Video) {

}

// RedisGetVideos 存的值是序列化对象
func RedisGetVideos(time string) ([]string, error) {
	result, err := VideoRdb.ZRevRange(context.Background(), "videos", 0, 29).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
