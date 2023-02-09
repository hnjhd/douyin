package dao

import "fmt"

// 初始化表格
func StructInit() error {
	fmt.Println(DB)
	M := DB.Migrator()
	if M.HasTable(&User{}) == false {
		M.CreateTable(&User{})
	}
	if M.HasTable(&LikeList{}) == false {
		M.CreateTable(&LikeList{})
	}
	if M.HasTable(&TableVideo{}) == false {
		M.CreateTable(&TableVideo{})
	}
	
	return nil
}
