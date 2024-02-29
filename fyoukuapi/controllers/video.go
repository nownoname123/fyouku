package controllers

import (
	"fyoukuapi/model"
	els "fyoukuapi/serve/es"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func ChannelHotList(c *gin.Context) {

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
	num, video, err := model.GetChannelHotList(id)
	if err != nil {
		ReturnError(c, 4002, "请求错误请稍后重试")
		return
	}
	k := int64(num)
	ReturnSuccess(c, 0, "success", video, k)
}
func Recommend(c *gin.Context) {
	cid := c.Query("channelId")
	id, err1 := strconv.Atoi(cid)
	if err1 != nil {
		log.Println("error in ATOI:", err1)
		return
	}
	rid := c.Query("regionId")
	rrid, err1 := strconv.Atoi(rid)
	if err1 != nil {
		log.Println("error in ATOI:", err1)
		return
	}
	if id == 0 {
		ReturnError(c, 4001, "必须指定频道")
		return
	}
	if rrid == 0 {
		ReturnError(c, 4002, "必须指定频道地区")
		return
	}
	num, video, err := model.GetChannelRecommendRegionList(id, rrid)
	if err != nil {
		ReturnError(c, 4002, "请求错误请稍后重试")
		return
	}
	k := int64(num)
	ReturnSuccess(c, 0, "success", video, k)
}
func RecommendByType(c *gin.Context) {
	id := c.Query("channelId")
	cid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("cid error:", err)
		return
	}
	id = c.Query("typeId")
	tid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("tid error:", err)
		return
	}
	if cid == 0 {
		ReturnError(c, 4001, "必须指定频道")
		return
	}
	if tid == 0 {
		ReturnError(c, 4002, "必须指定频道类型")
		return
	}
	num, video, err := model.GetChannelRecommendTypeList(cid, tid)
	if err != nil {
		ReturnError(c, 4002, "请求错误请稍后重试")
		return
	}
	k := int64(num)
	ReturnSuccess(c, 0, "success", video, k)
}

// VideoInfo 获取视频信息
func VideoInfo(c *gin.Context) {
	id := c.Query("videoId")
	vid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("vid error:", err)
		return
	}
	if vid == 0 {
		ReturnError(c, 4001, "必须指定视频id")
		return
	}
	video, err := model.RedisGetVideoInfo(vid)
	if err != nil {
		log.Println("获取视频信息失败：", err)
		ReturnError(c, 4004, "获取失败请稍后")
	}
	ReturnSuccess(c, 0, "success", video, 1)
}

// VideoEpisodesList 获取视频剧集列表
func VideoEpisodesList(c *gin.Context) {
	id := c.Query("videoId")
	vid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("vid error2:", err)
		return
	}
	if vid == 0 {
		ReturnError(c, 4001, "必须指定视频id")
		return
	}
	num, video, err := model.GetVideoEpisodesList(vid)
	if err != nil {
		ReturnError(c, 4002, "请求错误请稍后重试")
		return
	}
	k := int64(num)
	ReturnSuccess(c, 0, "success", video, k)
}
func UserVideo(c *gin.Context) {
	id := c.Query("uid")
	uid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("uid error in user video:", err)
		ReturnError(c, 4001, "必须指定用户")
		return
	}

	num, video, err := model.GetUserVideo(uid)
	if err != nil {
		ReturnError(c, 5000, "查询失败请重试")
		return
	}
	k := int64(num)
	ReturnSuccess(c, 0, "success", video, k)

}
func VideoSave(c *gin.Context) {
	playUrl := c.PostForm("playUrl")
	title := c.PostForm("title")
	subTitle := c.PostForm("subTitle")

	hh := c.PostForm("channelId")
	cid, err := strconv.Atoi(hh)
	if err != nil {
		log.Println("atoi error in cid :", err)
		return
	}
	if cid == 0 {
		ReturnError(c, 4001, "必须指定频道")
		return
	}

	hh = c.PostForm("regionId")
	rid, err := strconv.Atoi(hh)

	if err != nil {
		log.Println("atoi error in rid :", err)
		return
	}

	hh = c.PostForm("typeId")
	tid, err := strconv.Atoi(hh)

	if err != nil {
		log.Println("atoi error in tid :", err)
		return
	}
	id := c.PostForm("uid")
	uid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("uid error in user video:", err)
		ReturnError(c, 4001, "必须指定用户")
		return
	}
	err = model.SaveVideo(title, subTitle, cid, rid, tid, playUrl, uid)
	if err != nil {
		ReturnError(c, 5000, "保存失败请重试")
	}
	ReturnSuccess(c, 0, "success", "", 0)
}

// 导入ES脚本
func SendEs(c *gin.Context) {
	_, data, err := model.GetAllList()
	if err != nil {
		ReturnError(c, 4004, "查找视频数据错误")
		return
	}
	for _, v := range data {
		body := map[string]interface{}{
			"id":                   v.Id,
			"title":                v.Title,
			"sub_title":            v.SubTitle,
			"add_time":             v.AddTime,
			"img":                  v.Img,
			"img1":                 v.Img1,
			"episodes_count":       v.EpisodesCount,
			"is_end":               v.IsEnd,
			"channel_id":           v.ChannelId,
			"status":               v.Status,
			"region_id":            v.RegionId,
			"type_id":              v.TypeId,
			"episodes_update_time": v.EpisodesUpdateTime,
			"comment":              v.Comment,
			"user_id":              v.UserId,
			"is_recommend":         v.IsRecommend,
		}
		els.EsAdd("fyouku_video", "video-"+strconv.Itoa(v.Id), body)
	}
}
