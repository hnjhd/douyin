package controller

import (
	"TikTok/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
	NextTime  int64           `json:"next_time"`
}

type VideoListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}

// Feed /feed/
func Feed(c *gin.Context) {
	createTime := c.Query("latest_time")
	var lastTime time.Time
	if createTime != "0" {
		temp, _ := strconv.ParseInt(createTime, 10, 64)
		if len(createTime) > 10 {
			temp /= 1000
		}
		lastTime = time.Unix(temp, 0)
	} else {
		lastTime = time.Now()
	}
	videoService := GetVideoService()
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	feed, nextTime, err := videoService.Feed(lastTime, userId)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "获取视频失败",
			},
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		VideoList: feed,
		NextTime:  nextTime.Unix(),
	})
}

// Publish /publish/action/
func Publish(c *gin.Context) {
	token := c.Query("token")
	videoService := GetVideoService()
	userid, err := videoService.GetparseTokens(token)
	data, err := c.FormFile("data")
	log.Println("转化前token: ", token)
	log.Println("转化前id: ", userid)
	userId := int64(userid)
	log.Println("转化后id: ", userid)
	title := c.PostForm("title")
	if err != nil {
		log.Println("获取数据失败", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg: err.Error(),
		})
		return
	}
	// videoService := GetVideoService()
	err = videoService.Publish(data, userId, title)
	if err != nil {
		log.Println("上传视频失败", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg: "上传成功",
	})
}

// PublishList /publish/list/
func PublishList(c *gin.Context) {
	Id, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoService := GetVideoService()
	list, err := videoService.List(userId, curId)
	if err != nil {
		log.Println("获取视频列表失败")
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "获取列表失败",
			},
		})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg: "",
		},
		VideoList: list,
	})
}

func GetVideoService() service.VideoServiceImpl {
	var userService service.UserServiceImpl
	var videoService service.VideoServiceImpl
	videoService.UserService = &userService
	return videoService
}
