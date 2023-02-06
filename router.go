package main

import (
	"TikTok/controller"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")
	go apiRouter.GET("/user/", controller.UserInfo)
	go apiRouter.POST("/user/register/", controller.Register)
	go apiRouter.POST("/user/login/", controller.Login)
	go apiRouter.POST("/favorite/action", controller.FavoriteAction)
	go apiRouter.GET("favorite/list", controller.FavoriteList)
}
