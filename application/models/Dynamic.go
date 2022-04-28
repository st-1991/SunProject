package models

import (
	"SunProject/config"
	"gorm.io/gorm"
)

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

type ApiDynamic struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	Text string `json:"text"`
	Images string `json:"images"`
}

func (d Dynamic) TableName() string {
	return "keep_dynamic"
}

func (d Dynamic) CreateDynamic() bool {
	return config.DB.Create(&d).Error == nil
}

func (d Dynamic) GetDynamics(page, pageSize int) []ApiDynamic {
	var dynamics []ApiDynamic
	config.DB.Model(&d).Where("deleted = ?", 0).Offset((page - 1) * pageSize).Limit(pageSize).Order("id desc").Find(&dynamics)
	return dynamics
}

func (d Dynamic) IncrColumn(DB *gorm.DB, column string) bool {
	if d.Id == 0 {
		return false
	}
	allowColumns := []string{"thumb_up", "comment_num"}
	if !config.InArray(column, allowColumns) {
		return false
	}
	return DB.Model(&d).Where("id = ?", d.Id).Update(column, gorm.Expr(column + " + ?", 1)).Error == nil
}

