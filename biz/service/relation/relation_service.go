package relation

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/obrr-hhx/simpleDouyin/biz/model/common"
	"github.com/obrr-hhx/simpleDouyin/biz/model/social/relation"
	user_service "github.com/obrr-hhx/simpleDouyin/biz/service/user"
	"github.com/obrr-hhx/simpleDouyin/dal/db"
	"github.com/obrr-hhx/simpleDouyin/pkg/errno"
)

const (
	FOLLOW   = 1
	UNFOLLOW = 2
)

type RelationService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewRealtionService(ctx context.Context, c *app.RequestContext) *RelationService {
	return &RelationService{ctx: ctx, c: c}
}

// FollowAction follow or unfollow by request
func (r *RelationService) FollowAction(req *relation.DouyinRelationActionRequest) (flag bool, err error) {
	_, err = db.CheckUserExistById(req.ToUserId)
	if err != nil {
		return false, err
	}
	if req.ActionType != FOLLOW && req.ActionType != UNFOLLOW {
		return false, errno.ParamErr
	}
	current_user_id, _ := r.c.Get("current_user_id")
	if req.ToUserId == current_user_id.(int64) {
		return false, errno.ParamErr
	}
	new_follow_relation := &db.Follows{
		UserId:     req.ToUserId,
		FollowerId: current_user_id.(int64),
	}
	follow_exist, _ := db.QueryFollowExist(new_follow_relation.UserId, new_follow_relation.FollowerId)
	if req.ActionType == FOLLOW {
		if follow_exist {
			return false, errno.FollowRelationAlreadyExistErr
		}
		flag, err = db.AddNewFollow(new_follow_relation)
	} else {
		if !follow_exist {
			return false, errno.FollowRelationNotExistErr
		}
		flag, err = db.DeleteFollow(new_follow_relation)
	}
	return flag, err
}

// GetFollowList get follow list by request
func (r *RelationService) GetFollowList(req *relation.DouyinRelationFollowListRequest) (list []*common.User, err error) {
	_, err = db.CheckUserExistById(req.UserId)
	if err != nil {
		return nil, err
	}

	var follow_list []*common.User
	current_user_id, exists := r.c.Get("current_user_id")
	if !exists {
		current_user_id = int64(0)
	}
	followIdList, err := db.GetFollowList(req.UserId)
	if err != nil {
		return follow_list, err
	}

	// get the follow user info by their id
	for _, followId := range followIdList {
		userInfo, err := user_service.NewUserService(r.ctx, r.c).GetUserInfo(followId, current_user_id.(int64))
		if err != nil {
			continue
		}
		user := common.User{
			Id:              userInfo.Id,
			Name:            userInfo.Name,
			FollowCount:     userInfo.FollowCount,
			FollowerCount:   userInfo.FollowerCount,
			IsFollow:        userInfo.IsFollow,
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			Signature:       userInfo.Signature,
			TotalFavorited:  userInfo.TotalFavorited,
			WorkCount:       userInfo.WorkCount,
			FavoriteCount:   userInfo.FavoriteCount,
		}
		follow_list = append(follow_list, &user)
	}
	return follow_list, nil
}

// GetFollowerList get follower list by request
func (r *RelationService) GetFollowerList(req *relation.DouyinRelationFollowerListRequest) (list []*common.User, err error) {
	user_id := req.UserId
	var follower_list []*common.User
	current_user_id, exists := r.c.Get("current_user_id")
	if !exists {
		current_user_id = int64(0)
	}

	followerIdList, err := db.GetFollowerList(user_id)
	if err != nil {
		return follower_list, err
	}

	// get the follower user info by their id
	for _, followerId := range followerIdList {
		userInfo, err := user_service.NewUserService(r.ctx, r.c).GetUserInfo(followerId, current_user_id.(int64))
		if err != nil {
			continue
		}
		user := common.User{
			Id:              userInfo.Id,
			Name:            userInfo.Name,
			FollowCount:     userInfo.FollowCount,
			FollowerCount:   userInfo.FollowerCount,
			IsFollow:        userInfo.IsFollow,
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			Signature:       userInfo.Signature,
			TotalFavorited:  userInfo.TotalFavorited,
			WorkCount:       userInfo.WorkCount,
			FavoriteCount:   userInfo.FavoriteCount,
		}
		follower_list = append(follower_list, &user)
	}
	return follower_list, nil
}

// GetFriendList get friend list by request
func (r *RelationService) GetFriendList(req *relation.DouyinRelationFriendListRequest) (list []*common.User, err error) {
	user_id := req.UserId
	var friend_list []*common.User
	current_user_id, exists := r.c.Get("current_user_id")
	if !exists {
		current_user_id = int64(0)
	}

	FriendIdList, err := db.GetFriendList(user_id)
	if err != nil {
		return friend_list, err
	}

	// get the follow user info by their id
	for _, friendId := range FriendIdList {
		userInfo, err := user_service.NewUserService(r.ctx, r.c).GetUserInfo(friendId, current_user_id.(int64))
		if err != nil {
			continue
		}
		user := common.User{
			Id:              userInfo.Id,
			Name:            userInfo.Name,
			FollowCount:     userInfo.FollowCount,
			FollowerCount:   userInfo.FollowerCount,
			IsFollow:        userInfo.IsFollow,
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			Signature:       userInfo.Signature,
			TotalFavorited:  userInfo.TotalFavorited,
			WorkCount:       userInfo.WorkCount,
			FavoriteCount:   userInfo.FavoriteCount,
		}
		friend_list = append(friend_list, &user)
	}
	return friend_list, nil
}
