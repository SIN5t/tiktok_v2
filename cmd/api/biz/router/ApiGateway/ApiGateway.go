// Code generated by hertz generator. DO NOT EDIT.

package ApiGateway

import (
	ApiGateway "github.com/SIN5t/tiktok_v2/cmd/api/biz/handler/ApiGateway"
	"github.com/cloudwego/hertz/pkg/app/server"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_douyin := root.Group("/douyin", _douyinMw()...)
		{
			_feed := _douyin.Group("/feed", _feedMw()...)
			_feed.GET("/", append(_feed0Mw(), ApiGateway.Feed)...)
		}
		{
			_user := _douyin.Group("/user", _userMw()...)
			{
				_register := _user.Group("/register", _registerMw()...)
				_register.POST("/", append(_userregisterMw(), ApiGateway.UserRegister)...)
			}
		}
	}
}
