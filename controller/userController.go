package controller

import (
	"TikTok/dao"
	"TikTok/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type user struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func getuser(u dao.User) user {
	newuser := user{
		Id:            u.Model.ID,
		Name:          u.UserName,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	return newuser
}

// Register POST douyin/user/register/ 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	usi := service.UserServiceImpl{}
	u, err := usi.GetUserByUsername(username)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
			"user_id":     nil,
			"token":       nil,
		})
		return
	}
	if username == u.UserName {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "用户名已存在",
			"user_id":     nil,
			"token":       nil,
		})
	} else {
		newUser := dao.User{
			UserName: username,
			Password: password,
		}
		if fg, err := newUser.InsertUser(); fg != true {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 1,
				"status_msg":  err.Error(),
				"user_id":     nil,
				"token":       nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  nil,
			"user_id":     newUser.Model.ID,
			"token":       newUser.NewToken(),
		})
	}
}

// Login POST douyin/user/login/ 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	usi := service.UserServiceImpl{}
	u, err := usi.GetUserByUsername(username)
	if err != nil || u.Model.ID == 0 {
		if u.Model.ID == 0 {
			err = errors.New("用户不存在")
		}
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
			"user_id":     nil,
			"token":       nil,
		})
		//log.Println("Login ", err)
		return
	}
	password = usi.GetSha256(password)
	if u.Password == password {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  nil,
			"user_id":     u.Model.ID,
			"token":       u.NewToken(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "密码错误",
			"user_id":     nil,
			"token":       nil,
		})
	}
}

// UserInfo GET douyin/user/ 用户信息
func UserInfo(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")
	id, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
			"user":        nil,
		})
		return
	}
	usi := service.UserServiceImpl{}
	userid, err := usi.GetparseTokens(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
			"user":        nil,
		})
		return
	}
	if uint(id) != userid {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "token错误",
			"user":        nil,
		})
		return
	}
	u, err := usi.GetUserById(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
			"user":        nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  nil,
		"user":        getuser(u),
	})
}
