package dao

type LikeList struct {
	Id      int64
	VideoId int64
	UserId  int64
}

func (u *LikeList) InsertLikeList() (bool, string) {
	var num int64
	DB.Model(LikeList{}).Where("video_id = ? and user_id = ?", u.VideoId, u.UserId).Count(&num)
	if num != 0 {
		return false, "点赞失败,已经点赞"
	}
	M := DB.Create(&u)
	if M.Error != nil {
		return false, "点赞失败"
	}
	return true, "已经成功点赞"
}

func (u *LikeList) DeleteLikeList() (bool, string) {
	M := DB.Where("video_id = ? and user_id = ?", u.VideoId, u.UserId).Delete(&LikeList{})
	if M.Error != nil {
		return false, "取消点赞失败"
	}
	if M.RowsAffected == 0 {
		return false, "取消点赞失败"
	}
	return true, "取消点赞成功"
}

func FindUserLike(userid int64) []int64 {
	var LikeVideoList []int64
	DB.Model(LikeList{}).Where("user_id = ?", userid).Pluck("video_id", &LikeVideoList)
	return LikeVideoList
}
