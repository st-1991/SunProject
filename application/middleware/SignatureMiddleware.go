package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func VerifySign() gin.HandlerFunc {
	return func(c *gin.Context) {
		sign := c.GetHeader("Keep-Sign")
		if sign == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": -999,
				"message": "Invalid signature",
				"data": nil,
			})
			c.Abort()
			return
		}
		timestampInt, _ := strconv.ParseInt(c.GetHeader("Keep-Timestamp"), 10, 64)
		if time.Now().Unix() - timestampInt > 300 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": -998,
				"message": "Invalid signature",
				"data": nil,
			})
			c.Abort()
			return
		}

		newSign := computeSignature(c)
		if sign != newSign {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -997,
				"message": "Invalid signature",
				"data": nil,
			})
			c.Abort()
			return
		}
		//log.Println(c.Request.ContentLength)
		//computeSignature(c)
		c.Next()
	}
}

const SecretKey = "b4b5f02a1b4c925b1b1b4b5f02a1b4c9"

func computeSignature(c *gin.Context) string {
	//key := []byte(SecretKey)
	//h := hmac.New(sha256.New, key)

	// 添加时间戳和随机数
	timestamp := c.GetHeader("Keep-Timestamp")
	//if timestamp == "" {
	//	c.JSON(http.StatusUnauthorized, gin.H{
	//		"status": -999,
	//		"message": "Invalid signature",
	//		"data": nil,
	//	})
	//	c.Abort()
	//}

	//var buf bytes.Buffer
	//tee := io.TeeReader(c.Request.Body, &buf)
	//data, err := ioutil.ReadAll(tee)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"status": -999,
	//		"message": "参数读取失败",
	//		"data": nil,
	//	})
	//	c.Abort()
	//}
	//// 重新写入body
	//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	hash := md5.Sum([]byte(fmt.Sprintf("%s-%s", timestamp, SecretKey)))
	//dataWithTimeAndRand := fmt.Sprintf("%s-%s", timestamp, string(data))
	//dataWithTimeAndRand := "abc"

	//h.Write([]byte(dataWithTimeAndRand))

	// 对结果进行base64编码
	//signature := base64.StdEncoding.EncodeToString(h.Sum(nil))\
	signature := hex.EncodeToString(hash[:])
	return signature
}