package controllers

import (
	"fyoukuapi/model"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
)

type SendData struct {
	UserId    int
	MessageId int64
}

func Message(c *gin.Context) {
	id := c.PostForm("uids")
	content := c.PostForm("content")
	if id == "" {
		ReturnError(c, 4001, "请填写接收人uid")
		return
	}
	if content == "" {
		ReturnError(c, 4002, "发送内容不能为空")
		return
	}
	mid, err := model.SendMessageSave(content)
	if err != nil {
		ReturnError(c, 5000, "发送失败请重试")
	}
	uidConfig := strings.Split(id, ",")
	n := len(uidConfig)
	sendChannel := make(chan SendData, n)
	closeChannel := make(chan bool, n)
	for _, v := range uidConfig {

		uid, err := strconv.Atoi(v)
		if err != nil {
			log.Println("atoi error 1 :", err)
			ReturnError(c, 4001, "格式错误")
			return
		}
		var data SendData
		data.MessageId = mid
		data.UserId = uid
		sendChannel <- data
	}
	close(sendChannel)
	for i := 0; i < 5; i++ {
		go sendMessageFunc(closeChannel, sendChannel)
	}
	for i := 0; i < 5; i++ {
		<-closeChannel
	}
	ReturnSuccess(c, 0, "success", "", 0)
}
func sendMessageFunc(closeChannel chan bool, sendChannel chan SendData) {
	for i := range sendChannel {
		err := model.SendMessage(i.UserId, i.MessageId)
		if err != nil {
			log.Println("send error in :", i.MessageId, ": ", err)
		}
	}
	closeChannel <- true
}
