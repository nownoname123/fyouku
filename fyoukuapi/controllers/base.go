package controllers

import (
	"fyoukuapi/model"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func ChannelRegion(c *gin.Context) {
	id := c.Query("channelId")
	cid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("cid error:", err)
		return
	}
	if cid == 0 {
		ReturnError(c, 4001, "必须指定频道")
		return
	}
	num, region, err := model.GetChannelRegion(cid)
	if err != nil {
		ReturnError(c, 4002, "请求错误请稍后重试")
		return
	}
	k := int64(num)
	ReturnSuccess(c, 0, "success", region, k)
}

func ChannelType(c *gin.Context) {
	id := c.Query("channelId")
	cid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("cid error:", err)
		return
	}
	if cid == 0 {
		ReturnError(c, 4001, "必须指定频道")
		return
	}
	num, type1, err := model.GetChannelType(cid)
	if err != nil {
		ReturnError(c, 4002, "请求错误请稍后重试")
		return
	}
	k := int64(num)
	ReturnSuccess(c, 0, "success", type1, k)
}
func ChannelVideo(c *gin.Context) {
	hh := c.Query("channelId")
	cid, err := strconv.Atoi(hh)
	if err != nil {
		log.Println("atoi error in channelVideo :", err)
		return
	}
	if cid == 0 {
		ReturnError(c, 4001, "必须指定频道")
		return
	}
	hh = c.Query("regionId")
	rid, err := strconv.Atoi(hh)
	if hh == "" {
		err = nil
	}
	if err != nil {
		log.Println("atoi error in channelVideo :", err)
		return
	}
	hh = c.Query("typeId")
	tid, err := strconv.Atoi(hh)
	if hh == "" {
		err = nil
	}
	if err != nil {
		log.Println("atoi error in channelVideo :", err)
		return
	}
	end := c.Query("end")
	sort := c.Query("sort")

	hh = c.Query("channelId")
	limit, err := strconv.Atoi(hh)
	if hh == "" {
		err = nil
		limit = 0
	}
	if err != nil {
		log.Println("atoi error in channelVideo :", err)
		return
	}
	hh = c.Query("channelId")
	offset, err := strconv.Atoi(hh)
	if hh == "" {
		err = nil
		offset = 0
	}
	if err != nil {
		log.Println("atoi error in channelVideo :", err)
		return
	}
	if limit == 0 {
		limit = 12
	}
	num, video, err := model.GetChannelVideo(cid, rid, tid, end, sort, limit, offset)
	if err != nil {
		ReturnError(c, 4002, "请求错误请稍后重试")
		return
	}
	k := int64(num)
	ReturnSuccess(c, 0, "success", video, k)
}
