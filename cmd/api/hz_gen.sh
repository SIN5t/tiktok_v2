#!/bin/sh

if [ "$1" = "init" ]; then
    hz new -mod github.com/SIN5t/tiktok_v2/cmd/api
fi

hz update -idl ../../idl/relation.proto
hz update -idl ../../idl/user.proto
hz update -idl ../../idl/video.proto
hz update -idl ../../idl/favorite.proto
hz update -idl ../../idl/comment.proto
hz update -idl ../../idl/message.proto