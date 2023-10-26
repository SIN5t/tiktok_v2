package db

import (
	"fmt"
	commConfig "github.com/SIN5t/tiktok_v2/config"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
	"time"
)

var (
	VideoMysqlDB *gorm.DB
	redisConfig  = viper.Init("db")
	VideoRdb     *redis.Client
)

// BaseModel model ID and other info
type BaseModel struct {
	ID         int64     `gorm:"primarykey"`
	CreateTime time.Time `gorm:"index"`
	DeleteAt   gorm.DeletedAt
	IsDeleted  bool
}

// Video model video info
type Video struct {
	BaseModel
	AuthorID  int64  `gorm:"index:idx_video_authorid;not null"`
	VideoName string `gorm:"type:varchar(200);not null"`
	VideoPath string `gorm:"type:varchar(200);not null"`
	Title     string `gorm:"type:varchar(50);not null"`
}

func InitMysqlDB() {
	mysqlDSN := commConfig.GetDsn("mysql.Source")
	var err error

	VideoMysqlDB, err = gorm.Open(
		mysql.Open(mysqlDSN),
		&gorm.Config{
			Logger: logger.New(
				logrus.NewWriter(), // 自定义日志输出,来自openTelemetry
				logger.Config{
					SlowThreshold: time.Second, // 慢 SQL 阈值
					LogLevel:      logger.Info, // Log level
					Colorful:      true,        // 彩色打印
				})},
	)
	if err != nil {
		klog.Fatal("fail to initializeLog db: ", err.Error())
	}
	if err := VideoMysqlDB.Use(tracing.NewPlugin()); err != nil {
		klog.Fatalf("use tracing plugin failed: %s", err.Error())
	}
	err = VideoMysqlDB.AutoMigrate(&Video{})
	if err != nil {
		klog.Fatal("err create video table: %s", err.Error())
	}
}

func InitRdb() {
	VideoRdb = redis.NewClient(
		&redis.Options{
			Network:  "",
			Addr:     fmt.Sprintf("%s:%d", redisConfig.GetString("redis.addr"), redisConfig.GetInt("redis.port")),
			Password: redisConfig.GetString("redis.password"),
			DB:       redisConfig.GetInt("redis.video_db"),
			PoolSize: redisConfig.GetInt("redis.pool_size"),
		},
	)
}
