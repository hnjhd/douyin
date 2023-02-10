package ffmpeg

import (
	"TikTok/config"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

type Ffmsg struct {
	VideoName string
	ImageName string
}

var ClientSSH *ssh.Client
var Ffchan chan Ffmsg

func InitSSH() {
	var err error
	SSHConfig := &ssh.ClientConfig{
		Timeout: 5 * time.Second,
		User: config.SSH_USER,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if config.SSH_TYPE == "password" {
		SSHConfig.Auth = []ssh.AuthMethod{ssh.Password(config.SSH_PASSOWRD)}
	}
	addr := fmt.Sprintf("%s:%d", config.SSH_HOST, config.SSH_PORT)
	ClientSSH, err = ssh.Dial("tcp", addr, SSHConfig)
	if err != nil {
		log.Println("ssh client失败", err)
	}
	Ffchan = make(chan Ffmsg, config.SSH_MAX_MESSAGE_COUNT)
	go dispatcher()
	go keepAlive()
}

func dispatcher() {
	for ffmsg := range Ffchan {
		go func(f Ffmsg) {
			err := Ffmpeg(f.VideoName, f.ImageName)
			if err != nil {
				Ffchan <- f
				log.Println("重新派遣")
			}
		}(ffmsg)
	}
}

func Ffmpeg(videoName string, imageName string) error {
	session, err := ClientSSH.NewSession()
	if err != nil {
		log.Println("ssh session创建失败", err)
	}
	defer session.Close()
	command , err := session.CombinedOutput("ls;/usr/local/ffmpeg/bin/ffmpeg -ss 00:00:01 -i /home/ftpjc/videos/" + videoName + ".mp4 -vframes 1 /home/ftpjc/images/" + imageName + ".jpg")
	if err != nil {
		log.Println("session.CombinedOutput() 失败", command)
		return err
	}
	return nil
}

func keepAlive() {
	time.Sleep(time.Duration(config.SSH_HEARTBEAT_TIME) * time.Second)
	session, _ := ClientSSH.NewSession()
	session.Close()
}