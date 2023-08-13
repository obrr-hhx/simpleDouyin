package db

import (
	"time"

	"github.com/obrr-hhx/simpleDouyin/pkg/constants"
	"github.com/obrr-hhx/simpleDouyin/pkg/errno"
	"gorm.io/gorm"
)

type Comment struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	VideoID     int64          `json:"video_id"`
	CommentText string         `json:"comment_text"`
	CreatedAt   time.Time      `json:"create_at"`
	DeleteAt    gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

func (Comment) TableName() string {
	return constants.CommentsTableName
}

// AddNewComment add a comment
func AddNewComment(comment *Comment) error {
	if ok, _ := CheckUserExistById(comment.UserID); !ok {
		return errno.UserIsNotExistErr
	}
	if ok, _ := CheckVideoExistById(comment.VideoID); !ok {
		return errno.VideoIsNotExistErr
	}

	if err := DB.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func CheckCommentExistById(id int64) (bool, error) {
	var comment Comment
	if err := DB.Where("id = ?", id).Find(&comment).Error; err != nil {
		return false, err
	}
	if comment.ID == 0 {
		return false, nil
	}
	return true, nil
}

// DeleteCommentById delete a comment by id
func DeleteCommentById(id int64) error {
	if ok, _ := CheckCommentExistById(id); !ok {
		return errno.CommentIsNotExistErr
	}
	if err := DB.Where("id = ?", id).Delete(&Comment{}).Error; err != nil {
		return err
	}
	return nil
}

// GetCommentListByVideoId get comment list by video id
func GetCommentListByVideoId(videoId int64) ([]*Comment, error) {
	var comments []*Comment
	if ok, _ := CheckVideoExistById(videoId); !ok {
		return comments, errno.VideoIsNotExistErr
	}
	if err := DB.Table(constants.CommentsTableName).Where("video_id = ?", videoId).Find(&comments).Error; err != nil {
		return comments, err
	}
	return comments, nil
}

// GetCommentCountByVideoId get comment count by video id
func GetCommentCountByVideoId(videoId int64) (int64, error) {
	var count int64
	if ok, _ := CheckVideoExistById(videoId); !ok {
		return count, errno.VideoIsNotExistErr
	}
	if err := DB.Model(&Comment{}).Where("video_id = ?", videoId).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}
