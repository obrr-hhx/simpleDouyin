package db

import (
	"github.com/obrr-hhx/simpleDouyin/pkg/constants"
	"github.com/obrr-hhx/simpleDouyin/pkg/errno"
)

type User struct {
	ID              int64  `json:"id"`
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
}

func (User) TableName() string {
	return constants.UserTableName
}

// CreateUser create user info
func CreateUser(user *User) (int64, error) {
	err := DB.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

// QueryUser query user by user_name
func QueryUser(userName string) (*User, error) {
	user := &User{}
	if err := DB.Where("user_name = ?", userName).Find(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// QueryUserById query user by id
func QueryUserById(id int64) (*User, error) {
	user := &User{}
	if err := DB.Where("id = ?", id).Find(user).Error; err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errno.UserIsNotExistErr
	}

	return user, nil
}

// VarifyUser varify username and password in the database
func VarifyUser(userName, password string) (int64, error) {
	var user User
	if err := DB.Where("user_name = ? AND password = ?", userName, password).Find(&user).Error; err != nil {
		return 0, err
	}
	if user.ID == 0 {
		return 0, errno.PasswordIsNotVerified
	}
	return user.ID, nil
}

// CheckUserExistById check user if exist by id
func CheckUserExistById(id int64) (bool, error) {
	var user User
	if err := DB.Where("id = ?", id).Find(&user).Error; err != nil {
		return false, err
	}
	if user == (User{}) {
		return false, errno.UserIsNotExistErr
	}
	return true, nil
}
