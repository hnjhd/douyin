package service

import (
	"TikTok/dao"
	"mime/multipart"
	"time"
)

type Video struct {
	dao.TableVideo
	Author UserDTO `json:"author"`
	FavoriteCount int64 `json:"favorite_count"`
	CommentCount int64 `json:"comment_count"`
	IsFavorite bool `json:"is_favorite"`
}

type VideoService interface {
	// Feed
	// 传入当前时间和用户id返回视频
	Feed(lastTime time.Time, userId int64) ([]Video, time.Time, error)

	// GetVideo
	// 根据视频id和用户id查询视频
	GetVideo(videoId int64, userId int64) (Video, error)

	// Publish
	// 上传视频
	Publish(data *multipart.FileHeader, userId int64, title string) error

	// List
	// 通过用户id查询对应用户id视频
	List(userId int64, curId int64) ([]Video, error)

	// GetVideoIdList
	// 通过用户id查询用户发布的视频
	GetVideoIdList(userId int64) ([]int64, error)
}
