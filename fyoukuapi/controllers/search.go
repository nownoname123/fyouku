package controllers

import (
	"encoding/json"
	"fyoukuapi/model"
	els "fyoukuapi/serve/es"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func Search(c *gin.Context) {
	keyword := c.PostForm("keyword")
	num := c.Query("limit")
	var err error
	var limit int
	var offset int
	if num == "" {
		limit = 12
	} else {
		limit, err = strconv.Atoi(num)
		if err != nil {
			log.Println("atoi error in search :", err)
			limit = 12
		}
	}
	num = c.Query("offset")
	if num == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(num)
		if err != nil {
			log.Println("atoi error in search :", err)
			offset = 0
		}
	}
	if keyword == "" {
		ReturnError(c, 4001, "关键字不能为空")
		return
	}
	//定义排序条件
	sort := []map[string]string{{"id": "desc"}}
	//定义查询条件
	query := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": map[string]interface{}{
				"term": map[string]interface{}{
					"title": keyword,
				},
			},
		},
	}

	res := els.EsSearch("fyouku_video", query, offset, limit, sort)
	total := res.Total.Value
	var data []model.Video
	for _, v := range res.Hits {
		var itemData model.Video
		err := json.Unmarshal(v.Source, &itemData)
		if err != nil {
			log.Println("err in range res hits:", err)
			ReturnError(c, 5000, "发生错误请重试")
			return

		}
		data = append(data, itemData)
	}
	if total > 0 {
		ReturnSuccess(c, 0, "success", data, int64(total))
		return

	} else {
		ReturnError(c, 4004, "没有相关内容")
		return
	}
}
