package minio

import (
	"fmt"
	"github.com/SIN5t/tiktok_v2/pkg/viper"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/context"
)

var MinioClient *minio.Client

func InitMinIo() {
	// 初始化minioClient
	minioViper := viper.Init("minio")
	endpoint := fmt.Sprintf("%s:%d", minioViper.GetString("server.host"), minioViper.GetInt("server.port"))
	useSSL := false // 使用http连接要配置为false
	accessKeyID := minioViper.GetString("MINIO_ROOT_USER")
	secretAccessKey := minioViper.GetString("MINIO_ROOT_PASSWORD")
	videoBucket := minioViper.GetString("video_bucket")
	picBucket := minioViper.GetString("pic_bucket")
	location := minioViper.GetString("location")

	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		klog.Fatal(err.Error())
	}
	klog.Infof("%v\n", MinioClient)
	klog.Info("初始化成功")

	// Make new buckets
	err = MinioClient.MakeBucket(context.Background(), videoBucket, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := MinioClient.BucketExists(context.Background(), videoBucket)
		if errBucketExists == nil && exists {
			klog.Infof("We already own bucket： %s\n", videoBucket)
		} else {
			klog.Fatal(err)
		}
	}
	err = MinioClient.MakeBucket(context.Background(), picBucket, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := MinioClient.BucketExists(context.Background(), picBucket)
		if errBucketExists == nil && exists {
			klog.Infof("We already own bucket： %s\n", picBucket)
		} else {
			klog.Fatal(err)
		}
	}

}

func UploadFile(bucketName string, objectName string, filePath string, filetype string) error {

	// 上传的文件不显示特定的内容类型（例如二进制文件），则可以将 ContentType 设置为空字符串。
	info, err := MinioClient.FPutObject(context.Background(), bucketName, objectName, filePath+objectName, minio.PutObjectOptions{ContentType: filetype})
	if err != nil {
		klog.Debugf("error occurred %s\n", err.Error())
		return err
	}
	klog.Info(info)
	return nil
}

func DownloadFile(bucketName string, objectName string, fileDownloadPath string) error {
	err := MinioClient.FGetObject(context.Background(), bucketName, objectName, fileDownloadPath+objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
