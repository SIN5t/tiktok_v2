package response

import (
	"github.com/SIN5t/tiktok_v2/kitex_gen/relation"
	"github.com/SIN5t/tiktok_v2/kitex_gen/user"
)

type RelationAction struct {
	Base
}

type FollowerList struct {
	Base
	UserList []*user.User `json:"user_list"`
}

type FollowList struct {
	Base
	UserList []*user.User `json:"user_list"`
}

type FriendList struct {
	Base
	UserList []*relation.FriendUser `json:"user_list"`
}
