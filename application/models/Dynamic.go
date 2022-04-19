package models

import "SunProject/config"

// Dynamic 动态表.
type Dynamic struct {
	Id int `json:"id" gorm:"primaryKey"`
	UserId int `json:"user_id" gorm:"index:idx_userId;type:int(11);not null;default:0;comment:用户id"`
	Text string `json:"text" gorm:"type:varchar(1024);not null;default:'';comment:动态文案;"`
	Images string `json:"images" gorm:"type:varchar(512);not null;default:'';comment:动态图片/视频;"`
	Tag string `json:"tag" gorm:"type:varchar(64);not null;default:'';comment:标签;"`
	ThumbUp int `json:"thumb_up" gorm:"type:int(0);not null;default:0;comment:点赞数"`
	CommentNum int `json:"comment_num" gorm:"type:int(0);not null;default:0;comment:评论数"`
	GoodId int `json:"good_id" gorm:"type:int(0);not null;default:0;comment:关联商品ID"`
	Deleted int `json:"deleted" gorm:"type:int(0);not null;default:0;comment:是否删除"`
	Date `gorm:"embedded"`
}

func (d Dynamic) TableName() string {
	return "keep_dynamic"
}

func (d Dynamic) CreateDynamic() bool {
	return config.DB.Create(&d).Error == nil
}

