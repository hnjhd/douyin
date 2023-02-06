package controller

import (
	"TikTok/dao"
	"TikTok/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoid := c.Query("video_id")
	video_id, _ := strconv.ParseInt(videoid, 10, 64)
	actiontype := c.Query("action_type")
	action_type, _ := strconv.ParseInt(actiontype, 10, 64)
	usi := service.UserServiceImpl{}
	userid, err := usi.GetparseTokens(token)
	user_id := int64(userid)
	if err != nil {
		newLikeList := dao.LikeList{
			VideoId: video_id,
			UserId:  user_id,
		}

		if action_type == 1 {
			newLikeList.InsertLikeList()
			c.JSON(http.StatusOK, gin.H{
				"status_code": 0,
				"status_msg":  nil,
			})
		} else {
			newLikeList.DeleteLikeList()
			c.JSON(http.StatusOK, gin.H{
				"status_code": 0,
				"status_msg":  nil,
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		})
	}
}
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	usi := service.UserServiceImpl{}
	userid, err := usi.GetparseTokens(token)
	user_id := int64(userid)
	video_list := dao.FindUserLike(user_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{

			"status_code": 0,
			"status_msg":  nil,
			"video_list":  video_list,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		})
	}
}
