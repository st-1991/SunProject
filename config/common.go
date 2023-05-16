package config

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var ProjectPath string

type Result struct {
	Data interface{} `json:"data"`
	Status int `json:"status"`
	Msg string `json:"msg"`
}

func LogFile() *os.File {
	logFileName := "logs/" + time.Now().Format("2006-01-02") + ".log"
	f, err := os.Create(logFileName)
	if err != nil {
		println(err)
	}
	return f
}

func (r *Result) Success(c *gin.Context) {
	if r.Msg == "" {
		r.Msg = "操作成功"
	}
	if r.Data == nil {
		var data []string
		r.Data = data
	}
	c.JSON(http.StatusOK, r)
}

func (r Result) Error(c *gin.Context) {
	httpStatus := r.Status
	if httpStatus == 0 {
		httpStatus = http.StatusInternalServerError
	}
 	c.JSON(httpStatus, r)
}

func CreateCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Intn(100000))
}

func GenHMACMd5(ciphertext, key []byte) string {
	mac := hmac.New(md5.New, key)
	mac.Write(ciphertext)
	hmac := mac.Sum(nil)
	return fmt.Sprintf("%x", hmac)
}

func GetFileMd5(file *multipart.FileHeader) string {
	src, err := file.Open()
	if err != nil {
		return ""
	}
	defer src.Close()
	if err != nil {
		fmt.Errorf("打开文件失败，filename=%v, err=%v", file.Filename, err)
		return ""
	}
	md5h := md5.New()
	io.Copy(md5h, src)
	return hex.EncodeToString(md5h.Sum(nil))
}

func SliceColumn(input []map[string]interface{}, columnKey string) []interface{} {
	columns := make([]interface{}, 0, len(input))
	for _, val := range input {
		if v, ok := val[columnKey]; ok {
			columns = append(columns, v)
		}
	}
	return columns
}

func CreateCardMi(length int) string {
	rand.Seed(time.Now().UnixNano()) // 设置种子值为当前时间戳的纳秒部分

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") // 可用字符集合
	maxIndex := len(letters) - 1 // 最大索引值

	result := make([]rune, length)
	result[0] = letters[rand.Intn(52)] // 随机选择大小写字母作为首字符

	for i := 1; i < len(result); i++ {
		result[i] = letters[rand.Intn(maxIndex+1)]
		//if i == (len(result)-1)/2 { // 在中间位置插入一个横线符号（可省略）
		//	result[i] = '-'
		//}
	}
	return string(result)
}

func InArray(needle interface{}, haystack interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range haystack.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range haystack.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range haystack.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}