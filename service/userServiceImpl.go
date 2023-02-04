package service

import (
	"TikTok/dao"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

type UserServiceImpl struct {
	UserService
}

// 获得全部User对象
func (usi *UserServiceImpl) GetUserList() ([]dao.User, error) {
	return dao.GetUserList()
}

// 根据username获得User对象
func (usi *UserServiceImpl) GetUserByUsername(name string) (dao.User, error) {
	return dao.GetUserByUserName(name)
}

// 根据user_id获得User对象
func (usi *UserServiceImpl) GetUserById(id uint) (dao.User, error) {
	return dao.GetUserByUserID(id)
}

// 将User插入表内
func (usi *UserServiceImpl) InsertTUser(user *dao.User) bool {
	flag, _ := user.InsertUser()
	return flag
}

// 对password进行加密
func (usi *UserServiceImpl) GetSha256(Password string) string {
	return dao.GetSha256(Password)
}

// 对token进行解析
func parseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte("tiktok"), nil
	})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}

// 返回id
func (usi *UserServiceImpl) GetparseTokens(token string) (uint, error) {
	tokens, err := parseToken(token)
	if err != nil {
		return 0, err
	}
	userid, err := strconv.ParseInt(tokens.Id, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(userid), nil

}
