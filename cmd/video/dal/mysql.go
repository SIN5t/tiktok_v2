package db

import (
	"golang.org/x/net/context"
	"gorm.io/plugin/dbresolver"
	"time"
)

func GetVideosByLastTime(ctx context.Context, lastTime int64, limit int) (videoList []*Video, err error) {

	//Clauses为数据库查询添加额外的选项或条件,WithContext可以做：超时和取消、事务、日志等
	conn := VideoMysqlDB.Clauses(dbresolver.Read).WithContext(ctx)

	lastTimeUnixMilli := time.UnixMilli(lastTime)
	err = conn.Limit(limit).Order("create_time desc").Where("create_time < ?", lastTimeUnixMilli).Find(&videoList).Error
	if err != nil {
		return nil, err
	}

	return videoList, nil
}

func SaveVideoToMysql(video *Video) error {
	if err := VideoMysqlDB.Create(video).Error; err != nil {
		return err
	}
	return nil
}
