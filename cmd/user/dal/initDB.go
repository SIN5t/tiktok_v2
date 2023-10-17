package db

import (
	"fmt"
	commConfig "github.com/SIN5t/tiktok_v2/config"
	"github.com/SIN5t/tiktok_v2/kitex_gen/user"
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
	UserMysqlDB *gorm.DB
	redisConfig = viper.Init("db")
	UserRdb     *redis.Client
)

type UserReg struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	//主键
	UserId int64 `json:"user_id,omitempty" gorm:"primaryKey"`
}

func InitMysqlDB() {
	mysqlDSN := commConfig.GetDsn("mysql.Source")
	var err error

	UserMysqlDB, err = gorm.Open(
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
	if err := UserMysqlDB.Use(tracing.NewPlugin()); err != nil {
		klog.Fatalf("use tracing plugin failed: %s", err.Error())
	}
	err = UserMysqlDB.AutoMigrate(&UserReg{}, &user.User{}) // 注册和登入分表
	if err != nil {
		klog.Fatal("err create video table: %s", err.Error())
	}
}

func InitRdb() {
	UserRdb = redis.NewClient(
		&redis.Options{
			Network:  "",
			Addr:     fmt.Sprintf("%s:%d", redisConfig.GetString("redis.addr"), redisConfig.GetInt("redis.port")),
			Password: redisConfig.GetString("redis.password"),
			DB:       redisConfig.GetInt("redis.video_db"),
			PoolSize: redisConfig.GetInt("redis.pool_size"),
		},
	)
}
