package model

import (
	"fyoukuapi/dao"
	"time"
)

type Message struct {
	Id      int64
	Content string
	AddTime int64
}
type MessageUser struct {
	Id        int
	MessageId int64
	AddTime   int64
	UserId    int
	Status    int
}

// SendMessageSave 保存消息通知
func SendMessageSave(content string) (int64, error) {
	db := dao.Db
	db.AutoMigrate(&Message{})
	var message Message
	message.Content = content
	message.AddTime = time.Now().Unix()
	result := db.Create(&message)
	var nmessage Message
	db.Where("content = ? and add_time = ?", content, message.AddTime).First(&nmessage)
	id := nmessage.Id
	return id, result.Error

}

func SendMessage(uid int, mid int64) error {
	db := dao.Db
	db.AutoMigrate(&MessageUser{})
	var m MessageUser
	m.MessageId = mid
	m.UserId = uid
	m.Status = 1
	m.AddTime = time.Now().Unix()
	result := db.Create(&m)
	return result.Error
}
