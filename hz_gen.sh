#!/bin/sh

if [ "$1" = "init" ]; then
    hz new -mod github.com/obrr-hhx/simpleDouyin -force
    hz update -idl idl/publish.proto
    hz update -idl idl/relation.proto
fi

hz update -idl idl/user.proto
hz update -idl idl/feed.proto
hz update -idl idl/favorite.proto
hz update -idl idl/comment.proto
hz update -idl idl/message.proto