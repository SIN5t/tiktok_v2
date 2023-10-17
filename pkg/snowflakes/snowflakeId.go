package snowflakes

import (
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
)

// GenerateSnowFlakeId 生成一个雪花ID
func GenerateSnowFlakeId(nodeNum int64) int64 {

	node, err := snowflake.NewNode(nodeNum)
	if err != nil {
		klog.Fatal(err.Error())
	}
	return node.Generate().Int64()
}
