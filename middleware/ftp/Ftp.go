package ftp

import (
	"TikTok/config"
	"log"
	"time"

	"github.com/dutchcoders/goftp"
)

var TIKTOK_FTP *goftp.FTP

func InitFtp() {
	var err error
	TIKTOK_FTP, err = goftp.Connect(config.FTP_IP)
	if err != nil {
		log.Println("获取FTP失败")
	}
	err = TIKTOK_FTP.Login(config.FTP_USER, config.FTP_PASSWORD)
	if err != nil {
		log.Println("登录FTP失败")
	}
	go keepAlive()
}

func keepAlive() {
	time.Sleep(time.Duration(config.FTP_HEARTBEAT_TIME) * time.Second)
	TIKTOK_FTP.Noop()
}