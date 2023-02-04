package dao

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	Model    gorm.Model `gorm:"embedded"`
	UserName string     `gorm:"unique"` //用户名
	//NickName string     //昵称
	Password string //密码
}

// 保证一些数据的合法性
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserName == "" {
		return errors.New("用户名为空")
	}
	if len(u.UserName) > 32 {
		return errors.New("用户名过长")
	}
	if u.Password == "" {
		return errors.New("密码不符合规范")
	}

	return nil
}

// 生成相应的权限
func (u *User) NewToken() string {
	expiresTime := time.Now().Unix() + int64(86400)
	fmt.Printf("expiresTime: %v\n", expiresTime)
	id64 := int64(u.Model.ID)
	//fmt.Printf("id: %v\n", strconv.FormatInt(id64, 10))
	claims := jwt.StandardClaims{
		Audience:  u.UserName,
		ExpiresAt: expiresTime,
		Id:        strconv.FormatInt(id64, 10),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "tiktok",
		NotBefore: time.Now().Unix(),
		Subject:   "token",
	}
	var jwtSecret = []byte("tiktok")
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
		println("generate token success!\n")
		return token
	} else {
		println("generate token fail\n")
		return "fail"
	}
}

// 判定密码的合法性
func cheak(str string) error {
	if len(str) > 32 {
		return errors.New("密码过长")
	}
	if len(str) == 0 {
		return errors.New("密码为空")
	}
	return nil
}

// 对密码进行SHA256加密
func GetSha256(str string) string {
	if err := cheak(str); err != nil {
		return ""
	}
	srcByte := []byte(str)
	sha256New := sha256.New()
	sha256Bytes := sha256New.Sum(srcByte)
	sha256String := hex.EncodeToString(sha256Bytes)
	return sha256String
}

// 使用UserName(用户名进行查询)
func GetUserByUserName(UserName string) (User, error) {
	user := User{}
	M := DB.Where(User{UserName: UserName}).First(&user)
	err := M.Error
	if err != nil {
		//log.Println(err)
		return user, err
	}
	return user, nil
}

// 使用ID进行查询
func GetUserByUserID(UserID uint) (User, error) {
	user := User{}
	M := DB.Where(User{Model: gorm.Model{ID: UserID}}).First(&user)
	if err := M.Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
	// return tableUser, nil
}

// 将User插入数据库中
func (u *User) InsertUser() (bool, error) {
	u.Password = GetSha256(u.Password)
	M := DB.Create(&u)
	if M.Error != nil {
		return false, M.Error
	}
	if M.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

// GetUserList 获取全部TableUser对象
func GetUserList() ([]User, error) {
	tableUsers := []User{}
	if err := DB.Find(&tableUsers).Error; err != nil {
		log.Println(err.Error())
		return tableUsers, err
	}
	return tableUsers, nil
}
