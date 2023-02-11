package config

// 每次获取视频的数量
const VideoCount = 5

// 存储视频连接
const PLAY_URL_PREFIX = "http://124.71.58.18/videos/"
const COVER_URL_PREFIX = "http://124.71.58.18/images/"

// FTP 服务器地址
const FTP_IP = "124.71.58.18:21"
const FTP_USER = "ftpjc"
const FTP_PASSWORD = "ftpJCLee"
const FTP_HEARTBEAT_TIME = 60

// SSH 配置
const SSH_HOST = "124.71.58.18"
const SSH_USER = "ftpjc"
const SSH_PASSOWRD = "ftpJCLee"
const SSH_TYPE = "password"
const SSH_PORT = 22
const SSH_MAX_MESSAGE_COUNT = 100
const SSH_HEARTBEAT_TIME = 5 * 60