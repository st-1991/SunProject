package models

import (
	"SunProject/application/models/custom"
	"SunProject/config"
)

type Date struct {
	CreatedAt custom.JTime `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt custom.JTime `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
}

func Tables() []interface{} {
	return []interface{}{
		&User{},
		&Tabs{},
		&Comments{},
		&Files{},
		&Dynamic{},
	}
}

func Create(model interface{}) bool {
	return config.DB.Create(&model).Error == nil
}