package db

import (
	"github.com/SIN5t/tiktok_v2/config"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"time"
)

var DB *gorm.DB

// BaseModel model ID and other info
type BaseModel struct {
	ID         int64     `gorm:"primarykey"`
	CreateTime time.Time `gorm:""`
	UpdateAt   time.Time `gorm:"index"`
	DeleteAt   gorm.DeletedAt
	IsDeleted  bool
}

// Video model video info
type Video struct {
	BaseModel
	AuthorID int64 `gorm:"index:idx_video_authorid;not null"`

	PlayUrl  string `gorm:"type:varchar(200);not null"`
	CoverUrl string `gorm:"type:varchar(200);not null"`

	FavCount int64 `gorm:"type:int;default:0;not null"`
	ComCount int64 `gorm:"type:int;default:0;not null"`

	IsFavorite bool `gorm:"type:bool;default:false;not null"`

	Data  []byte `gorm:"column:video_data"`
	Title string `gorm:"type:varchar(50);not null"`
}

func InitDB() {

	mysqlDSN := config.GetDsn("mysql.Source")
	var err error

	DB, err = gorm.Open(
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
	/*if err := DB.Use(tracing.NewPlugin()); err != nil {
		klog.Fatalf("use tracing plugin failed: %s", err.Error())
	}*/
	err = DB.AutoMigrate(&Video{})
	if err != nil {
		klog.Fatal("err create video table: %s", err.Error())
	}
}
