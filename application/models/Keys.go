package models

import "SunProject/config"

type Keys struct {
	Id int `json:"id" gorm:"primaryKey"`
	Key string `json:"key" gorm:"index:id_account;type:varchar(128);not null;default:'';comment:key"`
	Status int8 `json:"status" gorm:"type:tinyint(1);not null;default:0;comment:状态 1启用0禁用"`
	Date `gorm:"embedded"`
}

func (Keys) TableName() string {
	return "keep_keys"
}

func getKey(k Keys) string {
	var res map[string]string
	result := config.DB.Table(k.TableName()).Select([]string{"key"}).Where("status = ?", 1).First(&res)
	if result.Error != nil {
		return ""
	}
	return res["key"]
}