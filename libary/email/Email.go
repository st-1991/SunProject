package email

import (
	"SunProject/config"
	"fmt"
	"gopkg.in/gomail.v2"
)

const (
	HOST = "smtp.163.com"
	PORT = 465
	USERNAME = "noticecode@163.com"
	PASSWORD = "DQXHGOZAZLZCIUCZ"
)

type Email struct {
	Host string
	Port int
	UserName string
	PassWord string
}

func InitMail() *gomail.Dialer {
	return gomail.NewDialer(HOST, PORT, USERNAME, PASSWORD)
}

type Message struct {
	To string
	GoMessage *gomail.Message
}

func (m Message) Title(title string) {
	m.GoMessage.SetHeader("Subject", title)
}

func (m Message) Content(content string)  {
	m.GoMessage.SetBody("text/html", content)
}

func Send(d *gomail.Dialer, m Message) error {
	m.GoMessage.SetHeader("From", USERNAME)
	m.GoMessage.SetHeader("To", m.To)
	if err := d.DialAndSend(m.GoMessage); err != nil {
		config.Logger().Error(fmt.Sprintf("发送验证码失败：%s", err))
		return err
	}
	return nil
}

//func (m Message) ContentPath(path string) error {
//	c, err := os.ReadFile(path)
//	if err != nil {
//		return err
//	}
//	content := strings.Replace(string(c), "[code]", );
//}
