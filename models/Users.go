package models

import (
	"SunProject/config"
	"encoding/json"
	"time"
)


var TableName = "keep_users"

type User struct {
	Id int `gorm:"primaryKey"`
	Phone string `gorm:"index:id_phone;type:char(11);not null;default:'';comment:手机号"`
	Nickname string `gorm:"type:varchar(32);not null;default:'';comment:用户昵称"`
	Sex int `gorm:"type:enum('0','1');not null;comment:性别：0女1男"`
	Avatar string `gorm:"type:varchar(256);not null;default:'';comment:头像"`
	ThumbUp int `gorm:"type:int(0);not null;default:0;comment:点赞数"`
	Fans int `gorm:"type:int(0);not null;default: 0;comment:粉丝数"`
	Focus int `gorm:"type:int(0);not null;default: 0;comment:关注数"`
	Balance float64 `gorm:"type:decimal(11, 2);not null;default:0.00;comment:余额"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return TableName
}

type UserDetail struct {
	Id int `json:"id"`
	Phone string `json:"phone"`
	json.Token
}

type ApiUser struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Sex int `json:"sex"`
}

func Users() []ApiUser {
	var users []ApiUser
	config.DB.Table(TableName).Select([]string{"id", "email", "sex"}).Find(&users)
	return users
}

func CreateUser(phone string, sex int, email string) bool {
	res := config.DB.Table(TableName).Create(&User{
		Phone: phone,
		Sex: sex,
	})
	if res.Error != nil {
		return false
	}
	return true
}

// UserDetails 获取用户详情
func UserDetails(userId int) ApiUser {
	var user ApiUser
	config.DB.Table(TableName).Select([]string{"id", "email", "sex"}).First(user)
	return user
}