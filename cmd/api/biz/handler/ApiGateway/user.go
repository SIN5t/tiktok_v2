// Code generated by hertz generator.

package ApiGateway

import (
	"context"
	"github.com/SIN5t/tiktok_v2/cmd/api/biz/model/ApiGateway"
	"github.com/SIN5t/tiktok_v2/cmd/api/rpc"
	config "github.com/SIN5t/tiktok_v2/config/const"
	"github.com/SIN5t/tiktok_v2/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// UserRegister .
// @router /douyin/user/register/ [POST]
func UserRegister(ctx context.Context, c *app.RequestContext) {

	req := &user.UserRegisterRequest{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	response, err := rpc.Register(ctx, req)
	if err != nil {
		c.JSON(consts.StatusOK, ApiGateway.DouyinUserRegisterResponse{
			StatusCode: config.FailResponse,
			StatusMsg:  "注册失败！",
		})
		hlog.Fatalf("rpc.Register failed, err: %v", err.Error())
	}
	//response 这里是kitex的，转为ApiGateway的比较好
	c.JSON(consts.StatusOK, response)

}
