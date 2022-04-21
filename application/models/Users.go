package models

import (
	"SunProject/config"
	"fmt"
	"gorm.io/gorm"
)


type User struct {
	Id int `json:"id" gorm:"primaryKey"`
	Phone string `json:"phone" gorm:"index:id_phone;type:char(11);not null;default:'';comment:手机号"`
	Nickname string `json:"nickname" gorm:"type:varchar(32);not null;default:'';comment:用户昵称"`
	Sex string `json:"sex" gorm:"type:enum(0,1);not null;comment:性别：0女1男"`
	Avatar string `json:"avatar" gorm:"type:varchar(256);not null;default:'';comment:头像"`
	Birthday string `json:"birthday"  gorm:"type:char(10);not null;default:'';comment:生日"`
	Profile string `json:"profile" gorm:"type:varchar(512);not null;default:'';comment:个人简介"`
	ThumbUp int `json:"thumb_up" gorm:"type:int(0);not null;default:0;comment:点赞数"`
	Fans int `json:"fans" gorm:"type:int(0);not null;default: 0;comment:粉丝数"`
	Focus int `json:"focus" gorm:"type:int(0);not null;default: 0;comment:关注数"`
	Balance float64 `json:"balance" gorm:"type:decimal(11, 2);not null;default:0.00;comment:余额"`
	Ip string `json:"ip" gorm:"type:varchar(32);not null;default:'';comment:登录ip"`
	Date `gorm:"embedded"`
}

func (User) TableName() string {
	return "keep_users"
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

func (u *User) CreateUser() bool {
	return config.DB.Create(u).Error == nil
}

func (u *User) EditUser() bool {
	res := config.DB.Model(u).Updates(u)
	if res.Error != nil {
		return false
	}
	return true
}

// GetUser 获取用户详情
func GetUser(phone string, id int) (User, bool) {
	var user User
	query := config.DB.Table(user.TableName())
	if phone != "" {
		query.Where("phone = ?", phone)
	}
	if id != 0 {
		query.Where("id = ?", id)
	}
	query.First(&user)
	if (user).Id == 0 {
		return user, false
	}
	return user, true
}

type UserBase struct {
	Id int `json:"id"`
	Avatar string `json:"avatar"`
	Nickname string `json:"nickname"`
}

func (u *User) GetUsersByIds(ids []int) []UserBase {
	var result []UserBase
	config.DB.Model(u).Select("id", "avatar", "nickname").Where(ids).Find(&result)
	return result
}

// CreateNickname 生成用户昵称
func CreateNickname() string {
	return fmt.Sprintf("小可爱-%s", config.CreateCode())
}

func CreateAvatar() string {
	return "http://www.gravatar.com/avatar/"
}

// SetThumbUp 点赞
func (u User) SetThumbUp(DB *gorm.DB) bool {
	if u.Id == 0 {
		return false
	}
	return DB.Model(&u).Where("id = ?", u.Id).Update("thumb_up", gorm.Expr("thumb_up + ?", 1)).Error == nil
}