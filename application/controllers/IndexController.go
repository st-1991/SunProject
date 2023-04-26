package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"os"
	"regexp"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"

	"SunProject/application/models"
	"SunProject/application/service"
	"SunProject/config"
	"SunProject/libary/email"
)

type UserLogin struct {
	Account string `form:"account" binding:"required"`
	Code string `form:"code" binding:"required"`
}

type LoginResult struct {
	UserDetail models.User `json:"user_detail"`
	Token      service.Token      `json:"token"`
}

func SendSms(c *gin.Context) {
	account := c.Query("account")
	if account == "" {
		ApiError(c, &Response{Code: -1, Msg: "请输入用户账号~"})
		return
	}

	rex := regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
	if res := rex.MatchString(account); !res {
		ApiError(c, &Response{-1, "请使用邮箱登录！", nil})
		return
	}
	code := config.CreateCode()

	gm := gomail.NewMessage()
	d := email.InitMail()
	m := email.Message{
		To:        account,
		GoMessage: gm,
	}
	m.Title("openAI验证码")
	fileC, err := os.ReadFile(config.ProjectPath + "/libary/email/temp/code.html")
	if err != nil {
		config.Logger().Error(fmt.Sprintf("文件打开失败：%s", err))
		ApiError(c, &Response{Code: -1, Msg: "发送失败-open file error"})
		return
	}
	content := strings.Replace(string(fileC), "[code]", code, -1)
	m.Content(content)
	if err := email.Send(d, m); err != nil {
		ApiError(c, &Response{Code: -1, Msg: "发送失败"})
		return
	}

	redisKey := config.RedisKey("sms:" + account)
	redisKey.PrefixKey().Set(code).Expire(300)
	ApiResponse(c, &Response{Msg: "发送成功，请注意查收！"})
}

func Login(c *gin.Context)  {
	var param UserLogin
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		ApiResponse(c, &Response{Code: -1, Msg: err.Error()})
		return
	}
	config.Logger().Info(param)

	//rex := regexp.MustCompile(`^(1[3-9][0-9]{9})$`)
	//if res := rex.MatchString(param.Account); !res {
	//	ApiResponse(c, &Response{-1, "手机号不正确！", new([]string)})
	//	return
	//}

	redisKey := config.RedisKey("sms:" + param.Account)
	code, err := redis.String(redisKey.PrefixKey().Get())
	if err != nil {
		ApiResponse(c, &Response{Code: -1, Msg: "请发送验证码！"})
		return
	}
	if code != param.Code {
		ApiResponse(c, &Response{Code: -1, Msg: "验证码错误！"})
		return
	}

	User := models.User{Account: param.Account}
	userDetails, ok := models.GetUser(param.Account, 0)
	if !ok {
		User = models.User{Account: param.Account, Nickname: models.CreateNickname(), Avatar: models.CreateAvatar(), Ip: c.ClientIP()}
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

	userData := service.UserData{ID: userDetails.Id, Account: userDetails.Account}
	userJson, _ := json.Marshal(&userData)
	j := service.Jwt{}
	token, err := j.CreateToken("user", string(userJson), 3600 * 24 * 30)
	if err != nil {
		config.Logger().Error("生成token失败！" + err.Error())
		ApiResponse(c, &Response{Code: -1, Msg: "系统错误！"})
	}

	redisKey.PrefixKey().Del()
	ApiResponse(c, &Response{Data: &LoginResult{UserDetail: userDetails, Token: token}})
}