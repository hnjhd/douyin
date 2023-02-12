package dao

import (
	"TikTok/config"
	"TikTok/middleware/ftp"
	"io"
	"log"
	"time"
)

type TableVideo struct {
	Id          int64 `json:"id"`
	AuthorId    int64
	PlayUrl     string `json:"play_url"`
	CoverUrl    string `json:"cover_url"`
	PublishTime time.Time
	Title       string `json:"title"`
}

// 根据作者id查询视频
func GetVideosByAuthorId(authorId int64) ([]TableVideo, error) {
	var data []TableVideo
	result := DB.Where(&TableVideo{AuthorId: authorId}).Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// 根据作者id查询视频集合id
func GetVideoIdsByAuthorId(authorId int64) ([]int64, error) {
	var id []int64
	result := DB.Model(&TableVideo{}).Where("author_id", authorId).Pluck("id", &id)
	if result.Error != nil {
		return nil, result.Error
	}
	return id, nil
}

// 根据视频id查询视频
func GetVideosByVideoId(videoId int64) (TableVideo, error) {
	var tableVideo TableVideo
	tableVideo.Id = videoId
	result := DB.First(&tableVideo)
	if result.Error != nil {
		return tableVideo, result.Error
	}
	return tableVideo, nil
}

// 通过FTP上传视频
func VideoFTP(file io.Reader, videoName string) error {
	ftp.TIKTOK_FTP.Cwd("~")
	err := ftp.TIKTOK_FTP.Cwd("videos")
	if err != nil {
		log.Println("切换到videos路径失败")
		return err
	}
	err = ftp.TIKTOK_FTP.Stor(videoName+".mp4", file)
	if err != nil {
		log.Println("上传失败")
		return err
	}
	return nil
}

// 通过FTP上传图片
func ImageFTP(file io.Reader, imageName string) error {
	err := ftp.TIKTOK_FTP.Cwd("images")
	if err != nil {
		log.Println("切换到images路径失败")
		return err
	}
	err = ftp.TIKTOK_FTP.Stor(imageName, file)
	if err != nil {
		log.Println("上传失败")
		return err
	}
	return nil
}

// 保存视频
func Save(videoName string, imageName string, authorId int64, title string) error {
	var tableVideo TableVideo
	tableVideo.PublishTime = time.Now()
	tableVideo.PlayUrl = config.PLAY_URL_PREFIX + videoName + ".mp4"
	tableVideo.CoverUrl = config.COVER_URL_PREFIX + imageName + ".jpg"
	tableVideo.AuthorId = authorId
	tableVideo.Title = title
	result := DB.Save(&tableVideo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 获取传入时间之前的视频
func GetVideoByLastTime(lastTime time.Time) ([]TableVideo, error) {
	videos := make([]TableVideo, config.VideoCount)
	result := DB.Where("publish_time < ?", lastTime).Order("publish_time desc").Limit(config.VideoCount).Find(&videos)
	if result.Error != nil {
		return videos, result.Error
	}
	return videos, nil
}
