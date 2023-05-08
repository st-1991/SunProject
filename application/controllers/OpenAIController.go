package controllers

import (
	"SunProject/application/models"
	"SunProject/application/service"
	"SunProject/config"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strings"
)

type Prompt struct {
	Prompt string `form:"prompt"`
	Options struct{
		Temperature float64 `form:"temperature"` // 回答性格 介于 0 和 2 之间。较高的值（如 0.8）将使输出更加随机，而较低的值（如 0.2）将使输出更加集中和确定
		PresencePenalty float64 `form:"presence_penalty"` // 惩罚机制
		FrequencyPenalty float64 `form:"frequency_penalty"`
	} `form:"options"`
	ParentMessageId string `form:"parentMessageId"`
}

const (
	CompletionIntegral = 10 // 对话消耗积分
	LeastImageIntegral = 80 // 最小绘画消耗积分
	MediumImageIntegral = 90 // 中等绘画消耗积分
	MaxImageIntegral = 100 // 最大绘画消耗积分
)

func Completions(c *gin.Context) {
	if !IsLogin(c) {
		return
	}
	var p Prompt
	if c.Bind(&p) != nil {
		ApiError(c, &Response{Code: -1, Msg: "参数绑定错误"})
		return
	}
	if p.Prompt == "" || strings.TrimSpace(p.Prompt) == "" {
		ApiError(c, &Response{Code:-1, Msg: "请输入您要咨询的问题"})
		return
	}

	userId := c.MustGet("userId").(int)
	// 用户积分余额判断
	if models.GetUserIntegral(userId) < CompletionIntegral {
		ApiError(c, &Response{Code:3, Msg: "积分不足请充值～"})
		return
	}

	var messages []service.Message
	if p.ParentMessageId != "" {
		var mDB models.Messages
		dbMessages := mDB.GetMessageByNo(p.ParentMessageId, 20)
		for _, dbM := range dbMessages {
			messages = append(messages, service.Message{Role: dbM.Role, Content: dbM.Content})
		}
	}
	messages = append(messages, service.Message{Role: "user", Content: p.Prompt})

	resp, err := service.ChatCompletions(messages)
	if err != nil {
		ApiError(c, &Response{Code:-1, Msg: err.Error()})
		return
	}
	// 余额扣减
	go service.ChangeIntegral(config.DB, userId, -10, 2, "对话消耗")
	// 设置为流
	c.Header("Content-Type", "text/event-stream")
	if p.ParentMessageId == "" {
		id := uuid.New()
		//p.ParentMessageId = strconv.FormatInt(time.Now().Unix(), 10)
		p.ParentMessageId = id.String()
	}
	go service.UserMessageComplete(p.ParentMessageId, p.Prompt)

	//scanner := bufio.NewScanner(resp.Body)
	//scanner.Buffer()
	reader := bufio.NewReaderSize(resp.Body, 100000)
	role := ""

	chanComp := make(chan service.Completion, 5)
	// 保存信息信道
	go service.SaveCompletionRaw(chanComp, p.ParentMessageId)
	//for scanner.Scan() {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("请求不完整")
			break
		}
		if line == "\n" {
			continue
		}
		config.Logger().Info("接收参数 ：" + line)
		if line == "[DONE]\n" {
			break
		}
 		completion := service.ParseEventStreamFields([]byte(line), p.ParentMessageId, role)
		chanComp <- completion
		if completion.Segment == "start" {
			role = completion.Role
		}
		completionB, _ := json.Marshal(completion)
		fmt.Fprint(c.Writer, string(completionB) + "\n")
		c.Writer.(http.Flusher).Flush()
	}
	// 关闭信道
	close(chanComp)
}

func CreateImages(c *gin.Context) {
	if !IsLogin(c) {
		return
	}
	var i service.ImageParams
	if err := c.Bind(&i); err != nil {
		ApiError(c, &Response{Code: -1, Msg: "参数绑定错误"})
		return
	}
	if i.Prompt == "" {
		ApiError(c, &Response{Code: -1, Msg: "请输入绘画修饰词"})
		return
	}
	if i.N < 1 || i.N > 10 {
		ApiError(c, &Response{Code: -1, Msg: "图片数量为 1～10张"})
		return
	}

	IntegralList := map[string]int{
		"256x256": LeastImageIntegral,
		"512x512": MediumImageIntegral,
		"1024x1024": MaxImageIntegral,
	}
	depleteIntegral, ok := IntegralList[i.Size]
	if !ok {
		ApiError(c, &Response{Code: -1, Msg: "图片尺寸错误"})
		return
	}

	userId := c.MustGet("userId").(int)
	// 用户积分余额判断
	if models.GetUserIntegral(userId) < depleteIntegral {
		ApiError(c, &Response{Code:3, Msg: "积分不足请充值～"})
		return
	}
	// 消耗余额
	go service.ChangeIntegral(config.DB, userId, -depleteIntegral, 2, "绘画消耗-" + i.Size)
	resp, err := service.ImagesGenerations(i)
	if err != nil {
		ApiError(c, &Response{Code:-1, Msg: err.Error()})
		return
	}
	respStruct := struct {
		Created int `json:"created"`
		Data []struct{
			Url string `json:"url"`
		} `json:"data"`
	}{}
	_ = json.Unmarshal(resp, &respStruct)
	ApiResponse(c, &Response{Data: respStruct.Data})
}