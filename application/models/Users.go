package models

import (
	"SunProject/config"
	"fmt"
	"time"
)


var TableName = "keep_users"

type User struct {
	Id int `json:"id" gorm:"primaryKey"`
	Phone string `json:"phone" gorm:"index:id_phone;type:char(11);not null;default:'';comment:手机号"`
	Nickname string `json:"nickname" gorm:"type:varchar(32);not null;default:'';comment:用户昵称"`
	Sex int `json:"sex" gorm:"type:enum('0','1');not null;comment:性别：0女1男"`
	Avatar string `json:"avatar" gorm:"type:varchar(256);not null;default:'';comment:头像"`
	ThumbUp int `json:"thumb_up" gorm:"type:int(0);not null;default:0;comment:点赞数"`
	Fans int `json:"fans" gorm:"type:int(0);not null;default: 0;comment:粉丝数"`
	Focus int `json:"focus" gorm:"type:int(0);not null;default: 0;comment:关注数"`
	Balance float64 `json:"balance" gorm:"type:decimal(11, 2);not null;default:0.00;comment:余额"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return TableName
}

type UserDetail struct {
	Id int `json:"id"`
	Phone string `json:"phone"`
	Nickname string `json:"nickname"`
	Sex int `json:"sex"`
	Avatar string `json:"avatar"`
	ThumbUp int `json:"thumb_up"`
	Fans int `json:"fans"`
	Focus int `json:"focus"`
	Balance float64 `json:"balance"`
}

type ApiUser struct {
	ID int `json:"id"`
	Phone string `json:"phone"`
	Sex int `json:"sex"`
}

func Users() []ApiUser {
	var users []ApiUser
	config.DB.Table(TableName).Select([]string{"id", "email", "sex"}).Find(&users)
	return users
}

func CreateUser(user *User) bool {
	res := config.DB.Table(TableName).Create(user)
	if res.Error != nil {
		return false
	}
	return true
}

// GetUser 获取用户详情
func (u *User) GetUser() (User, bool) {
	var user User
	config.DB.Where(u).First(&user)
	if (&user).Id == 0 {
		return user, false
	}
	return user, true
}

// CreateNickname 生成用户昵称
func CreateNickname() string {
	return fmt.Sprintf("小可爱-%s", config.CreateCode())
}

func CreateAvatar() string {
	return "http://www.gravatar.com/avatar/"
}