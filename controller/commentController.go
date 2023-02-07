package controller

import (
	"TikTok/dao"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// 评论
type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       user   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

// 评论查询的返回列表
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// 对于评论操作的返回json
type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

func GetList(videoId int64, userId int64) ([]Comment, error) {
	log.Println("GetList: running") //函数已运行
	//法一、使用SQL语句查询评论列表及用户信息，嵌套user信息。且导致提高耦合性。
	//1.查找CommentData结构体的信息
	commentData := make([]Comment, 1)
	result := dao.DB.Raw("select T.cid id,T.user_id user_id,T.`name`,T.follow_count,T.follower_count,"+
		"\nif(f.cancel is null,'false','true') is_follow,"+
		"\nT.comment_text content,T.create_date"+
		"\nfrom follows f right join\n("+
		"\n\tselect cid,vid,id user_id,`name`,comment_text,create_date,"+
		"\n\tcount(if(tag = 'follower' and cancel is not null,1,null)) follower_count,"+
		"\n\tcount(if(tag = 'follow' and cancel is not null,1,null)) follow_count"+
		"\n\tfrom\n\t("+
		"\n\t\tselect c.id cid,u.id,c.video_id vid,`name`,f.cancel,comment_text,create_date,'follower' tag"+
		"\n\t\tfrom comments c join users u on c.user_id = u.id and c.cancel = 0"+
		"\n\t\tleft join follows f on u.id = f.user_id and f.cancel = 0"+
		"\n\t\tunion all"+
		"\n\t\tselect c.id cid,u.id,c.video_id vid,`name`,f.cancel,comment_text,create_date,'follow' tag"+
		"\n\t\tfrom comments c join users u on c.user_id = u.id and c.cancel = 0"+
		"\n\t\tleft join follows f on u.id = f.follower_id and f.cancel = 0"+
		"\n\t\t) T\n\t\tgroup by cid,vid,id,`name`,comment_text,create_date"+
		"\n) T on f.follower_id = T.user_id and f.cancel = 0 and f.user_id = ?"+
		"\nwhere vid = ? group by cid order by create_date desc", userId, videoId).Scan(&commentData)

	err := result.Error

	if nil != err {
		log.Println("CommentService-GetList: sql error") //sql查询出错
		return nil, err
	}
	//当前有0条评论
	if result.RowsAffected == 0 {
		return nil, nil
	}
	//2.拼接
	commentInfoList := make([]Comment, 0, len(commentData))
	for _, comment := range commentData {
		userData := user{
			Id:   uint(comment.Id),
			Name: comment.User.Name,
		}
		_commentInfo := Comment{
			Id:      comment.Id,
			User:    userData,
			Content: comment.Content,
		}
		//3.组装list
		commentInfoList = append(commentInfoList, _commentInfo)
	}
	//-----------------------法一结束--------------------------

	sort.Sort(CommentSlice(commentInfoList))
	return commentInfoList, nil
}

// CommentSlice 此变量以及以下三个函数都是做排序-准备工作
type CommentSlice []Comment

func (a CommentSlice) Len() int { //重写Len()方法
	return len(a)
}
func (a CommentSlice) Swap(i, j int) { //重写Swap()方法
	a[i], a[j] = a[j], a[i]
}
func (a CommentSlice) Less(i, j int) bool { //重写Less()方法
	return a[i].Id > a[j].Id
}

// CommentAction
// 发表 or 删除评论 comment/action/
func CommentAction(c *gin.Context) {
	log.Println("CommentController-Comment_Action: running") //函数已运行
	//获取userId
	id, _ := c.Get("userId")
	userid, _ := id.(string)
	userId, err := strconv.ParseInt(userid, 10, 64)
	log.Printf("err:%v", err)
	log.Printf("userId:%v", userId)
	//错误处理
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "comment userId json invalid",
			},
		})
		log.Println("CommentController-Comment_Action: return comment userId json invalid") //函数返回userId无效
		return
	}
	//获取videoId
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	//错误处理
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{-1, "comment videoId json invalid"},
		})
		log.Println("CommentController-Comment_Action: return comment videoId json invalid") //函数返回视频id无效
		return
	}
	//获取操作类型
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	//错误处理
	if err != nil || actionType < 1 || actionType > 2 {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				-1,
				"comment actionType json invalid"},
		})
		log.Println("CommentController-Comment_Action: return actionType json invalid") //评论类型数据无效
		return
	}

	if actionType == 1 { //actionType为1，则进行发表评论操作
		content := c.Query("comment_text")

		//发表评论数据准备
		var sendComment dao.Comment
		sendComment.UserId = userId
		sendComment.VideoId = videoId
		sendComment.CommentText = content
		timeNow := time.Now()
		sendComment.CreateDate = timeNow
		//发表评论
		CommentInfo, err := dao.InsertComment(sendComment)
		//发表评论失败
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{
					-1,
					"send comment failed",
				},
			})
			log.Println("CommentController-Comment_Action: return send comment failed") //发表失败
			return
		}

		//发表评论成功:
		//返回结果
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "send comment success",
			},
			Comment: Comment{
				Id: CommentInfo.Id,
				User: user{
					Id: uint(userId),
				},
				Content:    content,
				CreateDate: timeNow.GoString(),
			},
		})
		log.Println("CommentController-Comment_Action: return Send success") //发表评论成功，返回正确信息
		return
	} else { //actionType为2，则进行删除评论操作
		//获取要删除的评论的id
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{
					StatusCode: -1,
					StatusMsg:  "delete commentId invalid",
				},
			})
			log.Println("CommentController-Comment_Action: return commentId invalid") //评论id格式错误
			return
		}
		//删除评论操作
		//err = commentService.DelComment(commentId)
		err = dao.DeleteComment(commentId)
		if err != nil { //删除评论失败
			str := err.Error()
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{
					StatusCode: -1,
					StatusMsg:  str,
				},
			})
			log.Println("CommentController-Comment_Action: return delete comment failed") //删除失败
			return
		}
		//删除评论成功
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "delete comment success",
			},
		})

		log.Println("CommentController-Comment_Action: return delete success") //函数执行成功，返回正确信息
		return
	}
}

// CommentList
// 查看评论列表 comment/list/
func CommentList(c *gin.Context) {
	log.Println("CommentController-Comment_List: running") //函数已运行
	//获取userId
	id, _ := c.Get("userId")
	userid, _ := id.(string)
	userId, err := strconv.ParseInt(userid, 10, 64)
	//log.Printf("err:%v", err)
	//log.Printf("userId:%v", userId)

	//获取videoId
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	//错误处理
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: -1,
			StatusMsg:  "comment videoId json invalid",
		})
		log.Println("CommentController-Comment_List: return videoId json invalid") //视频id格式有误
		return
	}
	log.Printf("videoId:%v", videoId)

	//commentService := new(service.CommentServiceImpl)
	commentList, err := GetList(videoId, userId)

	//commentList, err := commentService.GetListFromRedis(videoId, userId)
	if err != nil { //获取评论列表失败
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			},
		})
		log.Println("CommentController-Comment_List: return list false") //查询列表失败
		return
	}

	//获取评论列表成功
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "get comment list success",
		},
		CommentList: commentList,
	})
	log.Println("CommentController-Comment_List: return success") //成功返回列表
	return
}
