package db

import (
	"time"

	"github.com/obrr-hhx/simpleDouyin/pkg/constants"
)

type Video struct {
	ID          int64     `json:"id"`
	AuthorID    int64     `json:"author_id"`
	VideoUrl    string    `json:"video_url"`
	CoverUrl    string    `json:"cover_url"`
	PublishTime time.Time `json:"publish_time"`
	Title       string    `json:"title"`
}

func (Video) TableName() string {
	return constants.VideosTableName
}

// CreateVideo create video
func CreateVideo(video *Video) (int64, error) {
	err := DB.Create(video).Error
	if err != nil {
		return 0, err
	}
	return video.ID, err
}

// GetVideoByLastTime get video by last time
func GetVideoByLastTime(lastTime time.Time) ([]*Video, error) {
	videos := make([]*Video, constants.VideoFeedCount)
	if err := DB.Where("publish_time < ?", lastTime).Order("publish_time desc").Limit(constants.VideoFeedCount).Find(&videos).Error; err != nil {
		return videos, err
	}
	return videos, nil
}

// GetVideoByAuthorID get video by author id
func GetVideoByAuthorID(authorID int64) ([]*Video, error) {
	var videos []*Video
	if err := DB.Where("author_id = ?", authorID).Find(&videos).Error; err != nil {
		return videos, err
	}
	return videos, nil
}

// GetVideoListByVideoIDList get videos by video id list
func GetVideoListByVideoIDList(videoIDList []int64) ([]*Video, error) {
	var videos []*Video
	if err := DB.Where("id in ?", videoIDList).Find(&videos).Error; err != nil {
		return videos, err
	}
	return videos, nil
}

// GetWorkOut get the num of video published by the user
func GetWorkOut(userID int64) (int64, error) {
	var count int64
	if err := DB.Model(&Video{}).Where("author_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CheckVideoExist check video exist
func CheckVideoExistById(videoID int64) (bool, error) {
	var video Video
	if err := DB.Where("id = ?", videoID).Find(&video).Error; err != nil {
		return false, err
	}
	if video == (Video{}) {
		return false, nil
	}
	return true, nil
}

// DeleteVideo delete video
func DeleteVideo(videoID int64) (bool, error) {
	if err := DB.Where("id = ?", videoID).Delete(&Video{}).Error; err != nil {
		return false, err
	}
	return true, nil
}
