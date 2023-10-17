package initialize

import "github.com/SIN5t/tiktok_v2/cmd/api/rpc"

// Init :rpc\middleware(jwt)\logger
// 在Go语言中，编译器会自动调用init()函数来进行一些初始化操作。因此，“init”是一个特殊的函数名，不能被用作包名或变量名。
func Init() {
	rpc.InitVideoClient()
	rpc.InitUserClient()
}
