package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var DB *gorm.DB

func v(name string) string {
	env := "test."
	return viper.GetString(env + name)
}

type Writer struct {
}

func (w Writer) Printf(format string, args ...interface{})  {
	logFields := [4]string{
		"file", "runtime", "row", "sql",
	}
	fields := logrus.Fields{}
	for key, val := range args {
		fields[logFields[key]] = val
	}
	Logger().WithFields(fields).Info("SQL")
}

func init() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("database")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	newLogger := logger.New(
		Writer{},
		logger.Config{
			SlowThreshold: time.Second * 3,
			LogLevel: logger.Info,
			Colorful: false,
		},
	)

	dns := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local&timeout=2s", v("username"), v("password"), v("hostname"), v("database"))
	var err error
	DB, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Printf("mysql connect error %v", err)
		panic("数据库链接失败")
	}

	if DB.Error != nil {
		fmt.Printf("database error %v", DB.Error)
		panic("数据库查询失败")
	}
}