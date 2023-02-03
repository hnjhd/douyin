package main

import (
	"NewWork/dao"
	"github.com/gin-gonic/gin"
)

func main() {
	//初始化数据库
	dao.Init()
	//初始化表格
	dao.StructInit()
	//配置路由
	r := gin.Default()
	//初始化路由
	initRouter(r)
	//启动一个服务
	r.Run(":8080")
}
