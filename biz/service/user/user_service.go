package service

import (
	"context"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/obrr-hhx/simpleDouyin/biz/model/basic/user"
	"github.com/obrr-hhx/simpleDouyin/biz/model/common"
	"github.com/obrr-hhx/simpleDouyin/dal/db"
	"github.com/obrr-hhx/simpleDouyin/pkg/constants"
	"github.com/obrr-hhx/simpleDouyin/pkg/errno"
	"github.com/obrr-hhx/simpleDouyin/pkg/utils"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

// NewUserService create user service
func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{
		ctx: ctx,
		c:   c,
	}
}

// UserRegister register user and return user id
func (s *UserService) UserRegister(req *user.DouyinUserRegisterRequest) (user_id int64, err error) {
	user, err := db.QueryUser(req.Username)
	if err != nil {
		return 0, err
	}
	if user != nil {
		return 0, errno.UserAlreadyExistErr
	}
	passWord, err := utils.Crypt(req.Password)
	user_id, err = db.CreateUser(&db.User{
		UserName:        req.Username,
		Password:        passWord,
		Avatar:          constants.TestAva,
		BackgroundImage: constants.TestBackGround,
		Signature:       constants.TestSign,
	})
	return user_id, err
}

// UserInfo User api
func (s *UserService) UserInfo(req *user.DouyinUserRequest) (*common.User, error) {
	query_user_id := req.UserId
	current_user_id, exists := s.c.Get("current_user_id")
	if !exists {
		current_user_id = 0
	}

	return s.GetUserInfo(query_user_id, current_user_id.(int64))
}

// GetUserInfo
//
//	@Description: Query the information of query_user_id according to the current user user_id
//	@receiver *UserService
//	@param query_user_id int64
//	@param user_id int64  "Currently logged-in user id, may be 0"
//	@return *user.User
//	@return error
func (s *UserService) GetUserInfo(query_user_id int64, user_id int64) (*common.User, error) {
	u := &common.User{}
	errChan := make(chan error)
	defer close(errChan)
	var wg sync.WaitGroup
	wg.Add(7)

	// parallel do query work
	// Query user information
	go func() {
		defer wg.Done()
		user, err := db.QueryUserById(query_user_id)
		if err != nil {
			errChan <- err
			return
		}
		u.Name = user.UserName
		u.Avatar = utils.URLconvert(s.ctx, s.c, user.Avatar)
		u.BackgroundImage = utils.URLconvert(s.ctx, s.c, user.BackgroundImage)
		u.Signature = user.Signature
	}()

	// Get the number of the user published videos
	go func() {
		defer wg.Done()
		WorkCount, err := db.GetWorkOut(query_user_id)
		if err != nil {
			errChan <- err
			return
		}
		u.WorkCount = WorkCount
	}()

	// Get the number of user is following
	go func() {
		defer wg.Done()
		FollowCount, err := db.GetFollowCount(query_user_id)
		if err != nil {
			errChan <- err
			return
		}
		u.FollowCount = FollowCount
	}()

	// Get the number of user's fans
	go func() {
		defer wg.Done()
		FollowerCount, err := db.GetFollowerCount(query_user_id)
		if err != nil {
			errChan <- err
			return
		}
		u.FollowerCount = FollowerCount
	}()

	// Check if the user is followed
	go func() {
		defer wg.Done()
		if user_id != 0 {
			IsFollow, err := db.QueryFollowExist(user_id, query_user_id)
			if err != nil {
				errChan <- err
				return
			}
			u.IsFollow = IsFollow
		} else {
			u.IsFollow = false
		}
	}()

	// Get the total number of the user's likes
	go func() {
		defer wg.Done()
		LikeCount, err := db.GetFavoriteCountByUserID(query_user_id)
		if err != nil {
			errChan <- err
			return
		}
		u.FavoriteCount = LikeCount
	}()

	// Get the number of the user is favorited
	go func() {
		defer wg.Done()
		TotalFavorited, err := db.QueryTotalFavoritedByAuthorID(query_user_id)
		if err != nil {
			errChan <- err
			return
		}
		u.TotalFavorited = TotalFavorited
	}()

	wg.Wait()
	select {
	case result := <-errChan:
		return &common.User{}, result
	default:
	}
	u.Id = query_user_id
	return u, nil
}
