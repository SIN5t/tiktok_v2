package response

import "github.com/SIN5t/tiktok_v2/kitex_gen/message"

type MessageChat struct {
	Base
	MessageList []*message.Message `json:"message_list"`
}

type MessageAction struct {
	Base
}
