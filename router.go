package main

import (
	"TikTok/controller"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)

	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)

	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// videoController
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)
}
