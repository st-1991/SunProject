package models

import "SunProject/config"

type Tabs struct {
	Id int `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"type:varchar(16);not null;default:'';comment:标题"`
	Tag string `json:"tag" gorm:"type:varchar(16);not null;default:'';comment:标签"`
	Icon string `json:"icon" gorm:"type:varchar(255);default:'';not null;comment:图标"`
	FocusIcon string `json:"focus_icon" gorm:"type:varchar(255);default:'';not null;comment:选中图标"`
	Status int `json:"status" gorm:"type:tinyint(1);default:1;not null;comment:状态 1启用0禁用"`
	Sort int `json:"sort" gorm:"type:int(0);default:0;not null;comment:倒序排序"`
	Date `gorm:"embedded"`
}

func (t Tabs) TableName() string {
	return "keep_tabs"
}

type ApiTab struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Tag string `json:"tag"`
	Icon string `json:"icon"`
	FocusIcon string `json:"focus_icon"`
	Status int `json:"status"`
}

func (t *Tabs) GetTabs() []ApiTab {
	var tabs []ApiTab
	config.DB.Model(&t).Where("status = ?", 1).Order("sort desc").Find(&tabs)
	return tabs
}