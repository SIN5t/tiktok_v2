package service

import (
	"bytes"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
	"strings"
)

func ReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}

func GetJpegFromFfmpeg(filePath string, destFilePath string, fileName string, frameNum int) {

	//如果保存的路径不存在，则创建一个
	destPathWithoutFileName := strings.TrimSuffix(destFilePath, fileName)
	if _, err := os.Stat(destPathWithoutFileName); os.IsNotExist(err) {
		os.MkdirAll(destPathWithoutFileName, os.ModePerm)
	}

	reader := ReadFrameAsJpeg(filePath, frameNum)
	img, err := imaging.Decode(reader)
	if err != nil {
		klog.Fatal(err.Error())
	}
	err = imaging.Save(img, destFilePath)
	if err != nil {
		klog.Fatal(err.Error())
	}
}
