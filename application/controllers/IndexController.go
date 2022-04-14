package controllers

import (
	"encoding/json"
	"regexp"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"

	"SunProject/application/models"
	"SunProject/application/service"
	"SunProject/config"
)

type UserLogin struct {
	Phone string `form:"phone" binding:"required"`
	Code string `form:"code" binding:"required"`
}

type LoginResult struct {
	UserDetail models.User `json:"user_detail"`
	Token      service.Token      `json:"token"`
}

func SendSms(c *gin.Context) {
	phone := c.Query("phone")

	rex := regexp.MustCompile(`^(1[3-9][0-9]{9})$`)
	if res := rex.MatchString(phone); !res {
		ApiResponse(c, &Response{-1, "手机号不正确！", new([]string)})
		return
	}
	code := config.CreateCode()

	//if ok := libary.SendSms(phone, code); !ok {
	//	ApiResponse(c, &Response{Code: -1, Msg: "短信发送失败"})
	//	return
	//}

	redisKey := config.RedisKey("sms:" + phone)
	redisKey.PrefixKey().Set(code).Expire(300)
	ApiResponse(c, &Response{Msg: "发送成功，请注意查收！"})
}

func Login(c *gin.Context)  {
	var param UserLogin
	if err := c.Bind(&param); err != nil {
		ApiResponse(c, &Response{Code: -1, Msg: err.Error()})
		return
	}

	rex := regexp.MustCompile(`^(1[3-9][0-9]{9})$`)
	if res := rex.MatchString(param.Phone); !res {
		ApiResponse(c, &Response{-1, "手机号不正确！", new([]string)})
		return
	}

	redisKey := config.RedisKey("sms:" + param.Phone)
	code, err := redis.String(redisKey.PrefixKey().Get())
	if err != nil {
		ApiResponse(c, &Response{Code: -1, Msg: "请发送验证码！"})
		return
	}
	if code != param.Code {
		ApiResponse(c, &Response{Code: -1, Msg: "验证码错误！"})
		return
	}

	User := models.User{Phone: param.Phone}
	userDetails, ok := models.GetUser(param.Phone, 0)
	if !ok {
		User = models.User{Phone: param.Phone, Nickname: models.CreateNickname(), Avatar: models.CreateAvatar(), Sex: "1", Ip: c.ClientIP()}
		ok := User.CreateUser()
		if !ok {
			ApiResponse(c, &Response{Code: -1, Msg: "登陆失败，请重试！"})
			return
		}
		userDetails, _ = models.GetUser("", User.Id)
	} else {
		go func(userId int, ip string) {
			user := models.User{Id: userId, Ip: ip}
			user.EditUser()
		}(userDetails.Id, c.ClientIP())
	}

	userData := service.UserData{ID: userDetails.Id, Phone: userDetails.Phone}
	userJson, _ := json.Marshal(&userData)
	j := service.Jwt{}
	token, err := j.CreateToken("user", string(userJson), 3600 * 1)
	if err != nil {
		config.Logger().Error("生成token失败！" + err.Error())
		ApiResponse(c, &Response{Code: -1, Msg: "系统错误！"})
	}

	redisKey.PrefixKey().Del()
	ApiResponse(c, &Response{Data: &LoginResult{UserDetail: userDetails, Token: token}})
}