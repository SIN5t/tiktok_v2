// Code generated by hertz generator.

package ApiGateway

import (
	"github.com/SIN5t/tiktok_v2/cmd/api/mw/jwt"
	"github.com/cloudwego/hertz/pkg/app"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _douyinMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _feedMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _feed0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _registerMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userregisterMw() []app.HandlerFunc {
	// your code...
	return nil
}

// group : /publish/...
func _publishMw() []app.HandlerFunc {

	return []app.HandlerFunc{
		jwt.JwtMiddleware.MiddlewareFunc(),
	}
}

func _publishactionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _publishlistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userinfoMw() []app.HandlerFunc {
	return []app.HandlerFunc{
		jwt.JwtMiddleware.MiddlewareFunc(),
	}
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userloginMw() []app.HandlerFunc {
	// your code...
	return nil
}
