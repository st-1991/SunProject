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