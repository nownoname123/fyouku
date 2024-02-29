package controllers

import "github.com/gin-gonic/gin"
import "time"

type JsonStruct struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Items interface{} `json:"items"`
	Count int64       `json:"count"`
}

func ReturnSuccess(c *gin.Context, code int, msg interface{}, items interface{}, count int64) {
	json := &JsonStruct{Code: code, Msg: msg, Items: items, Count: count}
	c.JSON(200, json)

}

func ReturnError(c *gin.Context, code int, msg interface{}) {
	json := &JsonStruct{Code: code, Msg: msg}
	c.JSON(200, json)
}
func DateFormat(times int64) string {
	videoTime := time.Unix(times, 0)
	return videoTime.Format("2006-01-02")
}
