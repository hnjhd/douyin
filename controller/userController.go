package controller

import (
	"TikTok/dao"
	"errors"
	"fmt"
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
	u, err := dao.GetUserByUserName(username)
	//fmt.Println(u)
	if err != nil {
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
		//fmt.Println(c)
		fmt.Println(-2)
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
		//fmt.Println("注册返回的id: ", newUser.Model.ID)
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
	//fmt.Println("-orz2\n\n\n\n\n\n\n\n")
	u, err := dao.GetUserByUserName(username)
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
		return
	}
	//fmt.Println("-orz3\n\n\n\n\n\n\n\n")
	password = dao.GetSha256(password)
	//fmt.Println("-orz4\n\n\n\n\n\n\n\n")
	if u.Password == password {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  nil,
			"user_id":     u.Model.ID,
			"token":       u.NewToken(),
		})
		//fmt.Println(-3)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "密码错误",
			"user_id":     nil,
			"token":       nil,
		})
		//fmt.Println(-4)
	}
}

// UserInfo GET douyin/user/ 用户信息
func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	id, err := strconv.ParseInt(user_id, 10, 64)
	fmt.Println(id, err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
			"user":        nil,
		})
		return
	}
	//fmt.Println(-1)
	u, err := dao.GetUserByUserID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
			"user":        nil,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  nil,
			"user":        getuser(u),
		})
	}
}
