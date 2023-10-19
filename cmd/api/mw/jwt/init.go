package jwt

import (
	"context"
	"github.com/SIN5t/tiktok_v2/cmd/api/biz/model/ApiGateway"
	"github.com/SIN5t/tiktok_v2/cmd/api/rpc"
	config "github.com/SIN5t/tiktok_v2/config/const"
	"github.com/SIN5t/tiktok_v2/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
	"net/http"
	"time"
)

var (
	JwtMiddleware *jwt.HertzJWTMiddleware
	IdentityKey   = config.JwtIdentityKey
)

func InitJwt() {
	var err error
	JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:         "test zone",
		Key:           []byte("secret key"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, utils.H{
				"code":    code,
				"token":   token,
				"expire":  expire.Format(time.RFC3339),
				"message": "success",
			})
		},
		// LoginHandler 将调用该方法，返回值中含有userId,后面用于保存到jwt的MapClaims中
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var req = &ApiGateway.DouyinUserLoginRequest{}
			if err := c.BindAndValidate(req); err != nil {
				return nil, err
			}
			if len(req.Username) == 0 || len(req.Password) == 0 {
				return "", jwt.ErrMissingLoginValues
			}
			// 这里应该调用rpc的方法，因为user不在当前服务下，不能直接调它的service
			loginResponse, err := rpc.Login(ctx, &user.UserLoginRequest{Username: req.Username, Password: req.Password})
			if err != nil {
				return nil, err
			}

			return loginResponse.UserId, nil
		},
		// PayloadFunc：用于在登录时向 JWT 添加额外负载数据的回调函数 例如添加了userId，这个userId来自上面的Authenticator
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					IdentityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityKey: IdentityKey,
		// 这个是在外部调用auth的时候调用的,用于设置获取身份信息的函数，此处提取 token 的负载，并配合 IdentityKey 将用户id存入上下文信息。
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &ApiGateway.User{
				Id: claims[IdentityKey].(int64),
			}
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			hlog.CtxErrorf(ctx, "jwt biz err = %+v", e.Error())
			return e.Error()
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(http.StatusOK, utils.H{
				"code":    code,
				"message": message,
			})
		},
	})
	if err != nil {
		panic(err)
	}
}
