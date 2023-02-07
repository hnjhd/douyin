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
	// u := User{
	// 	UserName: "cylxxx",
	// 	Password: "1234567",
	// }
	// fg, err := u.InsertUser()
	// fmt.Println(u, "\n", fg, " ", err)
	// fmt.Println(u.NewToken())
	return nil
}
