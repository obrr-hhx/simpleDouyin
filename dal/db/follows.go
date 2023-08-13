package db

import (
	"time"

	"gorm.io/gorm"

	redis "github.com/obrr-hhx/simpleDouyin/mw/redis"
	"github.com/obrr-hhx/simpleDouyin/pkg/constants"
)

// Follows follower is fan of user
type Follows struct {
	ID         int64          `json:"id"`
	UserId     int64          `json:"user_id"`
	FollowerId int64          `json:"follower_id"`
	CreatedAt  time.Time      `json:"create_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

var rdFollows redis.Follows

func (Follows) TableName() string {
	return constants.FollowsTableName
}

// AddNewFollow add a follow and check if the id exists
func AddNewFollow(follow *Follows) (bool, error) {
	exist, err := QueryUserById(follow.UserId)
	if exist == nil || err != nil {
		return false, err
	}
	exist, err = QueryUserById(follow.FollowerId)
	if exist == nil || err != nil {
		return false, err
	}
	err = DB.Create(follow).Error
	if err != nil {
		return false, err
	}

	// add to redis
	if rdFollows.CheckFollow(follow.FollowerId) {
		rdFollows.AddFollow(follow.UserId, follow.FollowerId)
	}
	if rdFollows.CheckFollower(follow.UserId) {
		rdFollows.AddFollower(follow.UserId, follow.FollowerId)
	}

	return true, nil
}

// DeleteFollow delete a follow and check if the id exists and delete in redis
func DeleteFollow(follow *Follows) (bool, error) {
	exist, err := QueryUserById(follow.UserId)
	if exist == nil || err != nil {
		return false, err
	}
	exist, err = QueryUserById(follow.FollowerId)
	if exist == nil || err != nil {
		return false, err
	}
	err = DB.Where("user_id = ? AND follower_id = ?", follow.UserId, follow.FollowerId).Delete(follow).Error
	if err != nil {
		return false, err
	}

	// delete in redis
	if rdFollows.CheckFollow(follow.FollowerId) {
		rdFollows.DelFollow(follow.UserId, follow.FollowerId)
	}
	if rdFollows.CheckFollower(follow.UserId) {
		rdFollows.DelFollower(follow.UserId, follow.FollowerId)
	}

	return true, nil
}

// QueryFollowExist check the relation between user and follower if exist
func QueryFollowExist(userId, followerId int64) (bool, error) {
	// check in redis first
	if rdFollows.CheckFollow(followerId) {
		return rdFollows.ExistFollow(userId, followerId), nil
	}
	if rdFollows.CheckFollower(userId) {
		return rdFollows.ExistFollower(userId, followerId), nil
	}

	follow := &Follows{
		UserId:     userId,
		FollowerId: followerId,
	}

	if err := DB.Where("user_id = ? AND follower_id = ?", userId, followerId).Find(follow).Error; err != nil {
		return false, err
	}
	if follow.ID == 0 {
		return false, nil
	}
	return true, nil
}

// GetFollowCount query the number of the user is following
func GetFollowCount(follower_id int64) (int64, error) {
	// check in redis first
	if rdFollows.CheckFollow(follower_id) {
		return rdFollows.CountFollow(follower_id), nil
	}

	// query in database and add to redis
	followings, err := getFollowIdList(follower_id)
	if err != nil {
		return 0, err
	}
	go addFollowRelationToRedis(follower_id, followings)
	return int64(len(followings)), nil
}

// GetFollowerCount query the number of the user's follower
func GetFollowerCount(user_id int64) (int64, error) {
	// check in redis first
	if rdFollows.CheckFollower(user_id) {
		return rdFollows.CountFollower(user_id), nil
	}

	// query in database and add to redis
	followers, err := getFollowerIdList(user_id)
	if err != nil {
		return 0, err
	}
	go addFollowerRelationToRedis(user_id, followers)
	return int64(len(followers)), nil
}

// GetFollowList query the list of the user is following
func GetFollowList(follower_id int64) ([]int64, error) {
	if rdFollows.CheckFollow(follower_id) {
		return rdFollows.GetFollow(follower_id), nil
	}
	return getFollowIdList(follower_id)
}

// GetFollowerList query the list of the user's follower
func GetFollowerList(user_id int64) ([]int64, error) {
	if rdFollows.CheckFollower(user_id) {
		return rdFollows.GetFollower(user_id), nil
	}
	return getFollowerIdList(user_id)
}

// GetFriendList query the list of the user's friend
func GetFriendList(user_id int64) ([]int64, error) {
	if !rdFollows.CheckFollow(user_id) {
		followings, err := getFollowIdList(user_id)
		if err != nil {
			return *new([]int64), err
		}
		addFollowRelationToRedis(user_id, followings)
	}
	if !rdFollows.CheckFollower(user_id) {
		followers, err := getFollowerIdList(user_id)
		if err != nil {
			return *new([]int64), err
		}
		addFollowerRelationToRedis(user_id, followers)
	}
	return rdFollows.GetFriend(user_id), nil
}

func addFollowRelationToRedis(follower_id int64, followings []int64) {
	for _, following := range followings {
		rdFollows.AddFollow(follower_id, following)
	}
}

func addFollowerRelationToRedis(user_id int64, followers []int64) {
	for _, follower := range followers {
		rdFollows.AddFollower(user_id, follower)
	}
}

func getFollowIdList(follower_id int64) ([]int64, error) {
	var follows []Follows
	if err := DB.Where("follower_id = ?", follower_id).Find(&follows).Error; err != nil {
		return nil, err
	}
	var followings []int64
	for _, follow := range follows {
		followings = append(followings, follow.UserId)
	}
	return followings, nil
}

func getFollowerIdList(user_id int64) ([]int64, error) {
	var followers []Follows
	if err := DB.Where("user_id = ?", user_id).Find(&followers).Error; err != nil {
		return nil, err
	}
	var followerIds []int64
	for _, follower := range followers {
		followerIds = append(followerIds, follower.FollowerId)
	}
	return followerIds, nil
}
