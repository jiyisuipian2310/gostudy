package main

import "fmt"

//发送消息的基类
type ISendMsg interface {
	SendMessage(msg string) (length int, err error)
}

//通过邮件发送
type EmailSendMsg struct{}

func (s *EmailSendMsg) SendMessage(msg string) (length int, err error) {
	fmt.Printf("邮件发送消息: \n    %s\n", msg)
	length = len(msg)
	err = nil
	return
}

//通过微信发送
type WeixinSendMsg struct{}

func (s *WeixinSendMsg) SendMessage(msg string) (length int, err error) {
	fmt.Printf("微信发送消息: \n    %s\n", msg)
	length = len(msg)
	err = nil
	return
}

//通过短信发送消息
type SmsSendMsg struct{}

func (s *SmsSendMsg) SendMessage(msg string) (length int, err error) {
	fmt.Printf("短信发送消息: \n    %s\n", msg)
	length = len(msg)
	err = nil
	return
}

type Student struct {
	Name    string
	Age     int
	SendMap map[string]ISendMsg
}

func (s *Student) AddSendMode(mode string, sendor ISendMsg) {
	s.SendMap[mode] = sendor
}

func (s *Student) SendMessage(mode string) (length int, err error) {
	msg := fmt.Sprintf("name: %s, age: %d", s.Name, s.Age)
	sendor := s.SendMap[mode]
	if sendor == nil {
		return 0, fmt.Errorf("未找到发送模式：%s", mode)
	}

	return s.SendMap[mode].SendMessage(msg)
}

//设计模式中的开闭原则：对修改关闭，对扩展开放
func main() {
	s := &Student{
		Name:    "zhangsan",
		Age:     30,
		SendMap: make(map[string]ISendMsg),
	}

	var sendlen int
	var err error
	s.AddSendMode("EmailSend", &EmailSendMsg{})
	s.AddSendMode("WeixinSend", &WeixinSendMsg{})

	sendlen, err = s.SendMessage("EmailSend")
	if err != nil {
		fmt.Printf("通过邮件方式发送数据失败：%s\n\n", err.Error())
	} else {
		fmt.Printf("邮件方式发送了 %d 长度的数据\n\n", sendlen)
	}

	sendlen, err = s.SendMessage("WeixinSend")
	if err != nil {
		fmt.Printf("通过微信方式发送数据失败：%s\n\n", err.Error())
	} else {
		fmt.Printf("微信方式发送了 %d 长度的数据\n\n", sendlen)
	}

	//以上是原本的功能，现在需要添加新的发送功能
	//添加短信发送功能
	s.AddSendMode("SmsSend", &SmsSendMsg{})
	sendlen, err = s.SendMessage("SmsSend")
	if err != nil {
		fmt.Printf("通过短信方式发送数据失败：%s\n\n", err.Error())
	} else {
		fmt.Printf("短信方式发送了 %d 长度的数据\n\n", sendlen)
	}
}
