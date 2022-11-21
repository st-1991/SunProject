package models

import (
	"SunProject/config"
	"gorm.io/gorm"
)

// Comments 评论表.
type Comments struct {
	Id int `json:"id" gorm:"primaryKey"`
	DynamicId int `json:"dynamic_id" gorm:"index:idx_DynamicId;type:int(11);not null;default:0;comment:动态id"`
	TopCommentId int `json:"top_comment_id" gorm:"index:idx_TopCommentId;type:int(11);not null;default:0;comment:顶级评论id"`
	ParentCommentId int `json:"parent_comment_id" gorm:"type:int(11);not null;default:0;comment:父评论id"`
	ParentCommentUserId int `json:"parent_comment_user_id" gorm:"type:int(11);not null;default:0;comment:父评论用户id"`
	Level int `json:"level" gorm:"type:int(2);not null;default:0;comment:评论层级"`
	UserId int `json:"user_id" gorm:"index:idx_userId;type:int(11);not null;default:0;comment:用户id"`
	Text string `json:"text" gorm:"type:varchar(255);not null;default:'';comment:评论文案;"`
	Images string `json:"images" gorm:"type:varchar(255);not null;default:'';comment:动态图片/视频;"`
	ThumbUp int `json:"thumb_up" gorm:"type:int(0);not null;default:0;comment:点赞数"`
	CommentNum int `json:"comment_num" gorm:"type:int(0);not null;default:0;comment:评论数"`
	Deleted int `json:"deleted" gorm:"type:int(0);not null;default:0;comment:是否删除"`
	Date `gorm:"embedded"`
}

func (c Comments) TableName() string {
	return "keep_comments"
}

func (c *Comments) Create() bool {
	return config.DB.Create(&c).Error == nil
}

// CommenterInfo 评论人信息.
type CommenterInfo struct {
	Id int `json:"id"`
	TopCommentId int `json:"top_comment_id"`
	Level int `json:"level"`
	UserId int `json:"user_id"`
}

func GetCommenterInfoById(commentId int) CommenterInfo {
	cInfo := CommenterInfo{}
	config.DB.Model(&Comments{}).Select("id, top_comment_id, level, user_id").Where("id = ?", commentId).First(&cInfo)
	return cInfo
}

func (c Comments) IncrColumn(DB *gorm.DB, column string) bool {
	if c.Id == 0 {
		return false
	}
	allowColumns := []string{"thumb_up", "comment_num"}
	if !config.InArray(column, allowColumns) {
		return false
	}
	return DB.Model(&c).Where("id = ?", c.Id).Update(column, gorm.Expr(column + " + ?", 1)).Error == nil
}


func GetCommentsById(id, level, page, pageSize int) []Comments {
	var comments []Comments
	db := config.DB.Model(&Comments{}).Where("deleted = ?", 0).Where("level = ?", level)
	if level == 1 {
		db = db.Where("dynamic_id = ?", id)
	} else {
		db.Where("top_comment_id = ?", id)
	}
	db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&comments)
	return comments
}