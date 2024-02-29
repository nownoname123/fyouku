package controllers

import (
	"fyoukuapi/model"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func ChannelTop(c *gin.Context) {
	id := c.Query("channelId")
	if id == "" {
		ReturnError(c, 4001, "必须指定频道")
	}
	cid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("channel top atoi error:", err)
		return
	}
	num, video, err := model.GetChannelTop(cid)
	if err != nil {
		log.Println("db error in channel top :", err)
		ReturnError(c, 4005, "数据库出错请稍后再试")
		return
	}
	k := int64(num)
	ReturnSuccess(c, 0, "success", video, k)

}

func TypeTop(c *gin.Context) {
	id := c.Query("typeId")
	if id == "" {
		ReturnError(c, 4001, "必须指定频道")
	}
	tid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("type top atoi error:", err)
		return
	}
	num, video, err := model.GetTypeTop(tid)
	if err != nil {
		log.Println("db error in channel top :", err)
		ReturnError(c, 4005, "数据库出错请稍后再试")
		return
	}
	k := int64(num)
	ReturnSuccess(c, 0, "success", video, k)

}
