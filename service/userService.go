package service

import (
	"TikTok/dao"
)

type UserService interface {
	// GetTableUserList 获得全部TableUser对象
	GetUserList() ([]dao.User, error)

	// GetTableUserByUsername 根据username获得TableUser对象
	GetUserByUsername(name string) (dao.User, error)

	// GetTableUserById 根据user_id获得TableUser对象
	GetUserById(id uint) (dao.User, error)

	// InsertTableUser 将tableUser插入表内
	InsertTUser(user *dao.User) bool

	//对token进行解析,拿到的是user_id
	GetparseTokens(token string) (uint, error)
}
