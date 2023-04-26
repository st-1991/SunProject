package models

import (
	"SunProject/config"
)

type Messages struct {
	Id int `json:"id" gorm:"primaryKey"`
	MessageNo string `json:"message_no" gorm:"type:varchar(64);not null;default:'';comment:唯一值"`
	Role string `json:"role" gorm:"type:varchar(64);not null;default:'';comment:发送人"`
	Content string `json:"content" gorm:"type:varchar(2048);not null;default:'';comment:内容"`
	Date `gorm:"embedded"`
}

func (m Messages) TableName() string {
	return "keep_messages"
}

func (m *Messages) CreateMessage() bool {
	return config.DB.Create(m).Error == nil
}

type Message struct {
	Role string `json:"role"` // 消息作者 One of system, user, or assistant.
	Content string `json:"content"` // 文本内容
	//Name string `json:"name"`
}

func (m *Messages) GetMessageByNo(messageNo string, pageSize int) []Message {
	var result []Message
	config.DB.Model(m).Select("role", "content").Where("message_no = ?", messageNo).Limit(pageSize).Find(&result)
	return result
}