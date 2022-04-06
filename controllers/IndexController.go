package controllers

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"SunProject/config"
)

type UserLogin struct {
	User string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserDetails struct {
	Phone string
	Nickname string
	Age int
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
	//	ApiResponse(c, &Response{Code: -1, Message: "短信发送失败"})
	//	return
	//}

	redisKey := config.RedisKey("sms:" + phone)
	redisKey.PrefixKey().Set(code).Expire(300)
	ApiResponse(c, &Response{Message: "发送成功，请注意查收！"})
}

func Login(c *gin.Context)  {
	var form UserLogin
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if form.User != "root" || form.Password != "000000" {
		result := config.Result{Data: "", Status: 200, Msg: "账号或密码不正确"}
		result.Error(c)
		return
	}
	userInfo := make(map[string]interface{})
	userInfo["token"] = "adfafdafd"
	userInfo["details"] = UserDetails{
		Phone: "13785925782",
		Nickname: "你大哥",
		Age: 11,
	}
	result := config.Result{Data: userInfo, Status: 200, Msg: "操作成功"}
	result.Success(c)
}