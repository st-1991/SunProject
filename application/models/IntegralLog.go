package models

import (
	"SunProject/application/models/custom"
	"SunProject/config"
)

type IntegralLog struct {
	Id int `json:"id" gorm:"type:int(11);primaryKey;"`
	UserId int `json:"user_id" gorm:"type:int(11);noll null;default:0;"`
	Title string `json:"title" gorm:"type:varchar(126);not null;default:'';comment:变动内容"`
	Integral int `json:"integral" gorm:"type:int(11);unsigned;not null;default:0;comment:积分"`
	Type int `json:"type" gorm:"type:tinyint(1);not null;default:1;comment:1增加2减少"`
	Date `gorm:"embedded"`
}

func (IntegralLog) TableName() string {
	return "keep_integral_logs"
}

type IntegralLogItem struct {
	Title string `json:"title"`
	Integral int `json:"integral"`
	Type int `json:"type"`
	CreatedAt custom.JTime `json:"created_at"`
}

func GetIntegralLogsByUserId(userId, page, pageSize int) []IntegralLogItem {
	var logs []IntegralLogItem
	config.DB.Model(&IntegralLog{}).
		Where("user_id = ?", userId).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id desc").
		Find(&logs)
	return logs
}