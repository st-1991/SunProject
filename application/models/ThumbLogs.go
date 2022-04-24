package models

type ThumbLogs struct {
	Id int `json:"id" gorm:"primary_key"`
	ThumbType int `json:"Type" gorm:"uniqueIndex:uniq_log;type:tinyint(1);not null;default:1;comment:点赞类型1动态2评论"`
	UserId int `json:"user_id" gorm:"uniqueIndex:uniq_log;type:int(11);not null;default:0;"`
	BaseId int `json:"dynamic_id" gorm:"uniqueIndex:uniq_log;type:int(11);not null;default:0;comment:动态id或评论id"`
	Deleted int `json:"deleted" gorm:"type:tinyint(1);not null;default:0;comment:是否删除1是0否"`
	Date `gorm:"embedded"`
}