package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func VerifySign() gin.HandlerFunc {
	return func(c *gin.Context) {
		//computeSignature(c)
		log.Println(computeSignature(c))
		c.Next()
	}
}

const SecretKey = "b4b5f02a1b4c925b1b1b4b5f02a1b4c9"

func computeSignature(c *gin.Context) string {
	key := []byte(SecretKey)
	h := hmac.New(sha256.New, key)

	// 添加时间戳和随机数
	timestamp := c.GetHeader("Keep-Timestamp")
	if timestamp == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": -999,
			"message": "Invalid signature",
			"data": nil,
		})
		c.Abort()
	}

	//var buf bytes.Buffer
	//tee := io.TeeReader(c.Request.Body, &buf)
	//data, err := ioutil.ReadAll(tee)
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{
	//		"status": -999,
	//		"message": "参数读取失败",
	//		"data": nil,
	//	})
	//	c.Abort()
	//}
	data, _ := c.GetRawData()
	// 重新写入body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	var md5Str string
	if c.Request.ContentLength > 0 {
		randomStr := strconv.FormatInt(c.Request.ContentLength, 10) // 数字转字符串
		hash := md5.Sum([]byte(randomStr)) // 转成md5
		md5Str = hex.EncodeToString(hash[:])
	}

	dataWithTimeAndRand := fmt.Sprintf("%s-%s-%s", timestamp, md5Str, string(data))
	//config.Logger().Info(dataWithTimeAndRand)

	//log.Println(dataWithTimeAndRand)

	h.Write([]byte(dataWithTimeAndRand))

	// 对结果进行base64编码
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}