package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
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


