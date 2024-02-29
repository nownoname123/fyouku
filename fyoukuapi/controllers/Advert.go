package controllers

import (
	"fyoukuapi/model"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func ChannelAdvert(c *gin.Context) {
	cid := c.Query("channelId")
	id, err1 := strconv.Atoi(cid)
	if err1 != nil {
		log.Println("error in ATOI:", err1)
		return
	}
	if id == 0 {
		ReturnError(c, 4001, "必须指定频道")
		return
	}
	num, videos, err := model.GetChannelAdvert(id)
	if err != nil {
		ReturnError(c, 4002, "请求数据失败请稍后重试")
		return
	}
	k := int64(num)
	ReturnSuccess(c, 0, "success", videos, k)
}
