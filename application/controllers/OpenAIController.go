package controllers

import (
	"SunProject/application/models"
	"SunProject/application/service"
	"SunProject/config"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Prompt struct {
	Prompt string `form:"prompt"`
	Options struct{
		Temperature float64 `form:"temperature"` // 回答性格 介于 0 和 2 之间。较高的值（如 0.8）将使输出更加随机，而较低的值（如 0.2）将使输出更加集中和确定
		PresencePenalty int `form:"presence_penalty"` // 惩罚机制
		FrequencyPenalty int `form:"frequency_penalty"`
	} `form:"options"`
	ParentMessageId string `form:"parentMessageId"`
}

func Completions(c *gin.Context) {
	if !IsLogin(c) {
		return
	}
	var p Prompt
	if c.Bind(&p) != nil {
		ApiError(c, &Response{Code: -1, Msg: "参数绑定错误"})
		return
	}
	if p.Prompt == "" {
		ApiError(c, &Response{Code:-1, Msg: "请输入您要咨询的问题"})
		return
	}
	var messages []service.Message
	if p.ParentMessageId != "" {
		var mDB models.Messages
		dbMessages := mDB.GetMessageByNo(p.ParentMessageId)
		for _, dbM := range dbMessages {
			messages = append(messages, service.Message{Role: dbM.Role, Content: dbM.Content})
		}
	}
	messages = append(messages, service.Message{Role: "user", Content: p.Prompt})
	//ApiResponse(c, &Response{Data: messages})
	//return
	resp, err := service.ChatCompletions(messages)
	if err != nil {
		ApiError(c, &Response{Code:-1, Msg: err.Error()})
		return
	}
	// 设置为流
	c.Header("Content-Type", "text/event-stream")
	scanner := bufio.NewScanner(resp.Body)
	if p.ParentMessageId == "" {
		p.ParentMessageId = strconv.FormatInt(time.Now().Unix(), 10)
	}
	go service.UserMessageComplete(p.ParentMessageId, p.Prompt)
	role := ""
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		config.Logger().Info("接收参数 ：" + string(scanner.Bytes()))
		if string(scanner.Bytes()) == "[DONE]" {
			break
		}
 		completion := service.ParseEventStreamFields(scanner.Bytes(), p.ParentMessageId, role)
		if completion.Segment == "start" {
			role = completion.Role
		}
		completionB, _ := json.Marshal(completion)
		fmt.Fprint(c.Writer, string(completionB) + "\n")
		c.Writer.(http.Flusher).Flush()
	}
}