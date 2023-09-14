package init

import "github.com/SIN5t/tiktok_v2/cmd/api/rpc"

// Init :rpc\middleware(jwt)\logger
func Init() {
	rpc.InitVideoClient()
}
