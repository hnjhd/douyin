package service

import (
	"TikTok/config"
	"TikTok/dao"
	"TikTok/middleware/ffmpeg"
	"log"
	"mime/multipart"
	"time"
	"github.com/satori/go.uuid"
)

type VideoServiceImpl struct {
	UserService
}

// Feed
func (videoService VideoServiceImpl) Feed(lastTime time.Time, userId int64) ([]Video, time.Time, error) {
	videos := make([]Video, 0, config.VideoCount)
	tableVideos, err := dao.GetVideoByLastTime(lastTime)
	if err != nil {
		log.Println("dao.GetVideoByLastTime(lastTime) 失败", err)
		return nil, time.Time{}, err
	}
	err = videoService.copyVideos(&videos, &tableVideos, userId)
	if err != nil {
		log.Println("videoService.copyVideos(&videos, &tableVideos, userId) 失败", err)
		return nil, time.Time{}, err
	}
	return videos, tableVideos[len(tableVideos) - 1].PublishTime, nil
}

// GetVideo
func (videoService *VideoServiceImpl) GetVideo(videoId int64, userId int64) (Video, error) {
	var video Video
	data, err := dao.GetVideosByVideoId(videoId)
	if err != nil {
		log.Println("dao.GetVideosByVideoId(videoId) 失败", err)
		return video, err
	}
	videoService.createVideo(&video, &data, userId)
	return video, nil
}

// Publish
func (videoService *VideoServiceImpl) Publish(data *multipart.FileHeader, userId int64, title string) error {
	file, err := data.Open()
	if err != nil {
		log.Println("data.Open() 失败", err)
		return err
	}
	videoName := uuid.NewV4().String()
	log.Println("视频名字" + videoName)
	err = dao.VideoFTP(file, videoName)
	if err != nil {
		log.Println("dao.VideoFTP(file, videoName) 失败", err)
		return err
	}
	defer file.Close()
	imageName := uuid.NewV4().String()
	ffmpeg.Ffchan <- ffmpeg.Ffmsg{
		VideoName: videoName,
		ImageName: imageName,
	}
	err = dao.Save(videoName, imageName, userId, title)
	if err != nil {
		log.Println("dao.Save(videoName, imageName, userId, title) 失败", err)
		return err
	}
	return nil
}

// List
func (videoService *VideoServiceImpl) List(userId int64, curId int64) ([]Video, error) {
	data, err := dao.GetVideosByAuthorId(userId)
	if err != nil {
		log.Println("dao.GetVideosByAuthorId(userId) 失败", err)
		return nil, err
	}
	result := make([]Video, 0, len(data))
	err = videoService.copyVideos(&result, &data, curId)
	if err != nil {
		log.Println("videoService.copyVideos(&result, &data, curId) 失败", err)
		return nil, err
	}
	return result, nil
}

// GetVideoIdList
func (videoService *VideoServiceImpl) GetVideoIdList(authorId int64) ([]int64, error) {
	ids, err := dao.GetVideoIdsByAuthorId(authorId)
	if err != nil {
		log.Println("dao.GetVideoIdsByAuthorId(authorId) 失败", err)
		return nil, err
	}
	return ids, nil
}

// 拷贝视频
func (videoService *VideoServiceImpl) copyVideos(result *[]Video, data *[]dao.TableVideo, userId int64) error {
	for _, temp := range *data {
		var video Video
		videoService.createVideo(&video, &temp, userId)
		*result = append(*result, video)
	}
	return nil
}

// 插入视频
func (videoService *VideoServiceImpl) createVideo(video *Video, data *dao.TableVideo, userId int64) {
	var err error
	video.TableVideo = *data
	tableUser, err := videoService.GetUserById(uint(userId))
	if err != nil {
		log.Println("usi.GetUseById(uint(id)) 失败", err)
	}
	// TODO 后期拓展
	video.Author = UserDTO{
		Id:             userId,
		Name:           tableUser.UserName,
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
		TotalFavorited: 0,
		FavoriteCount:  0,
	}
}