// Code generated by hertz generator.

package relation

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	relation "github.com/obrr-hhx/simpleDouyin/biz/model/social/relation"
)

// RelationAction .
// @router douyin/relation/action/ [POST]
func RelationAction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.DouyinRelationActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(relation.DouyinRelationActionResponse)

	c.JSON(consts.StatusOK, resp)
}

// RelationFollowList .
// @router douyin/relation/follow/list/ [POST]
func RelationFollowList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.DouyinRelationFollowListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(relation.DouyinRelationFollowListResponse)

	c.JSON(consts.StatusOK, resp)
}

// RelationFollowerList .
// @router douyin/relation/follower/list/ [POST]
func RelationFollowerList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.DouyinRelationFollowerListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(relation.DouyinRelationFollowerListResponse)

	c.JSON(consts.StatusOK, resp)
}

// RelationFriendList .
// @router douyin/relation/friend/list/ [POST]
func RelationFriendList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req relation.DouyinRelationFriendListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(relation.DouyinRelationFriendListResponse)

	c.JSON(consts.StatusOK, resp)
}
