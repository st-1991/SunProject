package controllers

import (
	"SunProject/models"
	"github.com/garyburd/redigo/redis"
	"regexp"

	"github.com/gin-gonic/gin"

	"SunProject/config"
)

type UserLogin struct {
	Phone string `form:"phone" binding:"required"`
	Code string `form:"code" binding:"required"`
}

type LoginResult struct {
	UserDetail models.User `json:"user_detail"`
	Token string `json:"token"`
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
	userDetails, err := User.GetUser()
	if err != nil {
		userDetails = models.User{Phone: param.Phone, Nickname: models.CreateNickname(), Avatar: models.CreateAvatar()}
		models.CreateUser(&userDetails)
	}


	redisKey.PrefixKey().Del()
	ApiResponse(c, &Response{Data: &LoginResult{UserDetail: userDetails}})
}