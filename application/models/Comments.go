package models

// Comments 评论表.
type Comments struct {
	Id int `json:"id" gorm:"primaryKey"`
	DynamicId int `json:"dynamic_id" gorm:"index:idx_DynamicId;type:int(11);not null;default:0;comment:动态id"`
	ParentCommentId int `json:"parent_comment_id" gorm:"index:idx_ParentCommentId;type:int(11);not null;default:0;comment:父评论id"`
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