package controllers

import (
	"fyoukuapi/model"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type CommentInfo struct {
	Id           int            `json:"id"`
	Content      string         `json:"content"`
	AddTime      int64          `json:"addTime"`
	AddTimeTitle string         `json:"addTimeTitle"`
	UserId       int            `json:"userId"`
	Stamp        int            `json:"stamp"`
	PraiseCount  int            `json:"praiseCount"`
	UserInfo     model.UserInfo `json:"userinfo"`
	EpisodesId   int            `json:"episodesId"`
}

func CommentList(c *gin.Context) {
	id := c.Query("episodesId")
	eid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("eid error1:", err)
		return
	}
	id = c.Query("offset")
	offset, err := strconv.Atoi(id)

	if err != nil {
		log.Println("offset error1:", err)
		return
	}
	id = c.Query("limit")
	limit, err := strconv.Atoi(id)

	if err != nil {
		log.Println("limit error1:", err)
		return
	}
	if eid == 0 {
		ReturnError(c, 4001, "必须指定剧集")
		return
	}
	if limit == 0 {
		limit = 12
	}
	log.Println(eid, " ?? ", offset, " ?? ", limit)
	num, comment, err := model.GetCommentList(eid, offset, limit)
	if err != nil {
		ReturnError(c, 4004, "查询失败请重试")
		return
	}

	var data []CommentInfo
	var cif CommentInfo

	//获取uid channel
	uidChannel := make(chan int, 12)

	closeChannel := make(chan bool, 5)

	resChannel := make(chan model.UserInfo, 12)

	go func() {
		for _, v := range comment {
			uidChannel <- v.UserId
		}
		close(uidChannel)
	}()
	//开五个协程去取数据
	for i := 0; i < 5; i++ {
		go chanGetUserInfo(uidChannel, resChannel, closeChannel)
	}
	//判断是否完成，信息聚合，如果五个协程都完成才视为完成

	go func() {
		for i := 0; i < 5; i++ {
			<-closeChannel
		}
		close(resChannel)
		close(closeChannel)

	}()
	userInfoMap := make(map[int]model.UserInfo)
	//利用map建立id和userinfo的映射关系
	for r := range resChannel {
		userInfoMap[r.Id] = r
	}

	for _, v := range comment {
		cif.Id = v.Id
		cif.Content = v.Content
		cif.AddTime = v.AddTime
		cif.Stamp = v.Stamp
		cif.AddTimeTitle = DateFormat(v.AddTime)
		cif.UserId = v.UserId
		cif.PraiseCount = v.PraiseCount
		cif.UserInfo = userInfoMap[v.UserId]
		if err != nil {
			log.Println("get user info error:", err)
			ReturnError(c, 4005, "获取用户信息失败")
		}
		data = append(data, cif)

	}

	k := int64(num)
	/*for _, v := range data {
		log.Println(v.UserInfo.Id, " ", v.UserInfo.Name, " ", v.UserInfo.Avatar, " ", v.UserInfo.AddTime)
	}*/
	ReturnSuccess(c, 0, "success", data, k)
}
func chanGetUserInfo(uidChannel chan int, resChannel chan model.UserInfo, closeChannel chan bool) {
	for uid := range uidChannel {
		res, err := model.GetUserInfo(uid)
		if err == nil {
			resChannel <- res
		} else {
			log.Println("get user info err in chanGetUserInfo :", err)
		}
	}
	closeChannel <- true
}
func CommentSave(c *gin.Context) {
	content := c.PostForm("content")
	id := c.PostForm("uid")
	uid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("atoi err 2 :", err)
		return
	}
	id = c.PostForm("episodesId")
	eid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("atoi err 2 :", err)
		return
	}
	id = c.PostForm("videoId")
	vid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("atoi err 2 :", err)
		return
	}
	if content == "" {
		ReturnError(c, 4001, "评论不能为空")
		return
	}
	if uid == 0 {
		ReturnError(c, 4002, "请先登陆")
		return
	}
	if eid == 0 {
		ReturnError(c, 4003, "请先指定剧集id")
		return
	}
	if vid == 0 {
		ReturnError(c, 4004, "请先指定视频id")
	}

	err = model.SaveComment(content, uid, eid, vid)
	if err != nil {
		log.Println("error in save comment:", err)
		ReturnError(c, 5000, "保存评论失败")
	}
	ReturnSuccess(c, 0, "success", "", 0)
}
