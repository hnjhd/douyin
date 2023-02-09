package main

import (
	"TikTok/dao"
	"TikTok/middleware/ffmpeg"
	"TikTok/middleware/ftp"

	"github.com/gin-gonic/gin"
)

func main() {
	InitDependence()
	
	//配置路由
	r := gin.Default()
	//初始化路由
	initRouter(r)
	//启动一个服务
	r.Run(":8080")
}

// 初始化各种配置
func InitDependence() {
	// 初始化数据库
	dao.Init()
	// 初始化表格
	dao.StructInit()
	// 初始化FTP服务器
	ftp.InitFtp()
	// 初始化SSH
	ffmpeg.InitSSH()
}
