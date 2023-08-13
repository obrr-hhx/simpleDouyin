package redis

import "strconv"

const (
	followerSuffix = ":follower"
	followSuffix   = ":follow"
)

type Follows struct{}

// AddFollow add a follow
func (f Follows) AddFollow(userId, followerId int64) {
	add(rdbFollows, strconv.FormatInt(followerId, 10)+followSuffix, userId)
}

// AddFollower add a follower
func (f Follows) AddFollower(userId, followerId int64) {
	add(rdbFollows, strconv.FormatInt(userId, 10)+followerSuffix, followerId)
}

// DelFollow delete a follow
func (f Follows) DelFollow(userId, followerId int64) {
	del(rdbFollows, strconv.FormatInt(followerId, 10)+followSuffix, userId)
}

// DelFollower delete a follower
func (f Follows) DelFollower(userId, followerId int64) {
	del(rdbFollows, strconv.FormatInt(userId, 10)+followerSuffix, followerId)
}

// CheckFollow check if a follow exists
func (f Follows) CheckFollow(followerId int64) bool {
	return check(rdbFollows, strconv.FormatInt(followerId, 10)+followSuffix)
}

// CheckFollower check if a follower exists
func (f Follows) CheckFollower(userId int64) bool {
	return check(rdbFollows, strconv.FormatInt(userId, 10)+followerSuffix)
}

// ExistFollow check if a follow exists
func (f Follows) ExistFollow(userId, followerId int64) bool {
	return exist(rdbFollows, strconv.FormatInt(followerId, 10)+followSuffix, userId)
}

// ExistFollower check if a follower exists
func (f Follows) ExistFollower(userId, followerId int64) bool {
	return exist(rdbFollows, strconv.FormatInt(userId, 10)+followerSuffix, followerId)
}

// CountFollow count the number of follows
func (f Follows) CountFollow(followerId int64) int64 {
	sum, _ := count(rdbFollows, strconv.FormatInt(followerId, 10)+followSuffix)
	return sum
}

// CountFollower count the number of followers
func (f Follows) CountFollower(userId int64) int64 {
	sum, _ := count(rdbFollows, strconv.FormatInt(userId, 10)+followerSuffix)
	return sum
}

// GetFollow get the follows of a user
func (f Follows) GetFollow(followerId int64) []int64 {
	return get(rdbFollows, strconv.FormatInt(followerId, 10)+followSuffix)
}

// GetFollower get the followers of a user
func (f Follows) GetFollower(userId int64) []int64 {
	return get(rdbFollows, strconv.FormatInt(userId, 10)+followerSuffix)
}

// GetFriend get the friend of the id via intersection
func (f Follows) GetFriend(userId int64) (friends []int64) {
	ks1 := strconv.FormatInt(userId, 10) + followerSuffix
	ks2 := strconv.FormatInt(userId, 10) + followSuffix
	v, _ := rdbFollows.SInter(ks1, ks2).Result()
	for _, vs := range v {
		v_i64, _ := strconv.ParseInt(vs, 10, 64)
		friends = append(friends, v_i64)
	}
	return friends
}
