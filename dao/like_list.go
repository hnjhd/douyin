package dao

type LikeList struct {
	Id      int64
	VideoId int64
	UserId  int64
}

func (u *LikeList) InsertLikeList() (bool, error) {
	M := DB.Create(&u)
	if M.Error != nil {
		return false, M.Error
	}
	if M.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (u *LikeList) DeleteLikeList() (bool, error) {
	M := DB.Delete(&u)
	if M.Error != nil {
		return false, M.Error
	}
	if M.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func FindUserLike(userid int64) []int64 {
	var LikeVideoList []int64
	DB.Where(LikeList{UserId: userid}).Pluck("user_id", &LikeVideoList)
	return LikeVideoList
}
