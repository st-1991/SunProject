package service

import (
	"SunProject/application/models"
	"SunProject/config"
	"SunProject/libary/request"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"time"
)


type Message struct {
	Role string `json:"role"` // 消息作者 One of system, user, or assistant.
	Content string `json:"content"` // 文本内容
	//Name string `json:"name"`
}

type Parameters struct {
	Model string `json:"model"`
	Stream bool `json:"stream"` // 是否流的形式返回
	Temperature float64 `json:"temperature"` // 回答性格 介于 0 和 2 之间。较高的值（如 0.8）将使输出更加随机，而较低的值（如 0.2）将使输出更加集中和确定
	PresencePenalty int `json:"presence_penalty"` // 惩罚机制
	FrequencyPenalty int `json:"frequency_penalty"`
	Messages []Message `json:"messages"` // 所有描述对话的消息列表
}

const HOST = "http://154.40.59.110:8686/"

type OpenAI struct {
	Host string
	ApiKey string
	Router string
}

func NewOpenAI(apiKey, Router string) OpenAI {
	return OpenAI{
		HOST,
		apiKey,
		Router,
	}
}

func optionsDefault() Parameters {
	return Parameters{
		Model: "gpt-3.5-turbo",
		Stream: true,
		Temperature: 0.8,
		PresencePenalty: 1,
		FrequencyPenalty: 1,
	}
}

func ChatCompletions(messages []Message) (*http.Response, error) {
	apiKey := models.GetKey()
	if apiKey == "" {
		return nil, fmt.Errorf("未获取的可用key")
	}
	openAi := NewOpenAI(apiKey, "chat/completions")
	p := optionsDefault()
	p.Messages = messages
	headers := map[string]string{
		"Api-Key":openAi.ApiKey,
	}
	config.Logger().Info("header", headers)
	params, _ := json.Marshal(p)
	res := request.JsonPost(openAi.Host + openAi.Router, params, headers)
	if res.Err != nil {
		return nil, res.Err
	}
	if res.StatusCode() != http.StatusOK {
		body, _ := res.Body()
		errBody := struct {
			Status int `json:"status"`
			Msg string `json:"msg"`
		}{}
		_ = json.Unmarshal(body, &errBody)
		config.Logger().Error(fmt.Sprintf("发送请求失败：%+v", errBody))
		return nil, fmt.Errorf(errBody.Msg)
	}
	return res.Response, nil
}

// CompletionRaw 原始返回信息：{"id":"chatcmpl-77KJEVys7qgcl9v8T9WhTNI5DQksM","object":"chat.completion.chunk","created":1681980720,"model":"gpt-3.5-turbo-0301","choices":[{"delta":{},"index":0,"finish_reason":"stop"}]}
type CompletionRaw struct {
	Id string `json:"id"`
	Object string `json:"object"`
	Created int64 `json:"created"`
	Model string `json:"model"`
	Choices []struct{
		Delta struct{
			Role string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type Completion struct {
	Id string `json:"id"`
	Role string `json:"role"`
	Text string `json:"text"`
	DateTime string `json:"dateTime"`
	Segment string `json:"segment"`
	ParentMessageId string `json:"parentMessageId"`
}

func ParseEventStreamFields(p []byte, parentMessageId, Role string) Completion {
	var compRaw CompletionRaw
	_ = json.Unmarshal(p, &compRaw)
	comp := Completion{
		Id: compRaw.Id,
	}
	comp.DateTime = time.Unix(compRaw.Created, 0).Format("2006-01-02 15:04:05")
	var segment string
	if len(compRaw.Choices) > 0 {
		comp.Role = compRaw.Choices[0].Delta.Role
		comp.Text = compRaw.Choices[0].Delta.Content
		comp.ParentMessageId = parentMessageId
		if compRaw.Choices[0].FinishReason == "stop" {
			comp.Segment = "stop"
			return comp
		}
	}
	if segment == "" && comp.Role != "" {
		comp.Segment = "start"
		return comp
	}
	comp.Segment = "text"
	comp.Role = Role
	return comp
}

func SaveCompletionRaw(chanComp chan Completion, messageNo string) {
	//redisKey := config.RedisKey("CompletionRaw:" + c.Id)
	//cStr, _ := json.Marshal(c)
	//if _, err := config.Redo("lpush", redisKey, cStr); err != nil {
	//	config.Logger().Error("原始信息保存失败" + err.Error())
	//}
	var (
		role string
		content string
	)

	for c := range chanComp{
		if c.Segment == "stop" {
			continue
		}
		if role == "" {
			role = c.Role
		}
		content += c.Text
	}
	mDB := models.Messages{
		MessageNo: messageNo,
		Role: role,
		Content: content,
	}
	mDB.CreateMessage()
}

func MessageComplete(MessageNo, CompletionId string) {
	var Content string
	var role string
	redisKey := config.RedisKey("CompletionRaw:" + CompletionId)
	for {
		m, err := redis.String(config.Redo("rpop", redisKey))
		if err != nil {
			config.Logger().Error(err.Error())
			break
		}
		var completionRaw CompletionRaw
		_ = json.Unmarshal([]byte(m), &completionRaw)
		config.Logger().Info("数据",completionRaw)
		if role == "" && completionRaw.Choices[0].Delta.Role != "" {
			role = completionRaw.Choices[0].Delta.Role
		}
		if completionRaw.Choices[0].Delta.Content != "" {
			Content += completionRaw.Choices[0].Delta.Content
		}
	}

	mDB := models.Messages{
		MessageNo: MessageNo,
		Role: role,
		Content: Content,
	}
	mDB.CreateMessage()
}

func UserMessageComplete(MessageNo, message string)  {
	mDB := models.Messages{
		MessageNo: MessageNo,
		Role: "user",
		Content: message,
	}
	mDB.CreateMessage()
}

func send()  {

}