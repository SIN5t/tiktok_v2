package jwt

import (
	"context"
	db2 "github.com/SIN5t/tiktok_v2/cmd/user/dal"
	"github.com/SIN5t/tiktok_v2/cmd/user/service"
	utils2 "github.com/SIN5t/tiktok_v2/pkg/jwt/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
	"net/http"
	"time"
)

var (
	JwtMiddleware *jwt.HertzJWTMiddleware
	IdentityKey   = "identity"
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
		// login会调用这里的方法！所以这里要自己写
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginStruct struct {
				Username string `form:"account" json:"account" query:"username" vd:"(len($) > 0 && len($) < 30); msg:'Illegal format'"`
				Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 30); msg:'Illegal format'"`
			}
			if err := c.BindAndValidate(&loginStruct); err != nil {
				return nil, err
			}
			// TODO 这里应该调用rpc的方法，因为将来rpc是在当前的模块下的，而user不在，不能直接调它的service

			users, err := service.CheckUser(loginStruct.Username, utils2.MD5(loginStruct.Password))
			if err != nil {
				return nil, err
			}

			return users, nil
		},
		IdentityKey: IdentityKey,
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &db2.UserReg{
				UserId: claims[IdentityKey].(int64),
			}
		},
		// PayloadFunc：用于在登录时向 JWT 添加额外负载数据的回调函数
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*db2.UserReg); ok {
				return jwt.MapClaims{
					IdentityKey: v.UserId,
				}
			}
			return jwt.MapClaims{}
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
