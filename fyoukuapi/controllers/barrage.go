package controllers

import (
	"encoding/json"
	"fyoukuapi/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsData struct {
	CurrentTime int
	EpisodesId  int
	ViedoID     int
}

func BarrageWs(c *gin.Context) {
	var (
		conn     *websocket.Conn
		err      error
		data     []byte
		barrages []model.BarrageData
	)
	conn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		goto ERR
	}
	//轮询
	for {
		if _, data, err = conn.ReadMessage(); err != nil {
			goto ERR

		}
		var wsData WsData
		err = json.Unmarshal(data, &wsData)
		if err != nil {
			log.Println("json error :", err)
			goto ERR

		}
		endTime := wsData.CurrentTime + 60
		//获取弹幕数据
		_, barrages, err = model.GetBarrageList(wsData.EpisodesId, wsData.CurrentTime, endTime)
		if err == nil {
			if err := conn.WriteJSON(barrages); err != nil {
				goto ERR

			}
		}

	}
ERR:
	log.Println("socket error : ", err)
	err = conn.Close()
	if err != nil {
		log.Println("close connect error:", err)
		return
	}

}

func BarrageSave(c *gin.Context) {
	id := c.PostForm("uid")
	uid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("strconv error in barrage save1 :", err)
		return
	}
	content := c.PostForm("content")
	time1 := c.PostForm("currentTime")
	ct, err := strconv.Atoi(time1)
	log.Println(ct)
	if err != nil {
		log.Println("strconv error in barrage save2 :", err)
		return
	}
	id = c.PostForm("episodesId")
	eid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("strconv error in barrage save3 :", err)
		return
	}
	id = c.PostForm("videoId")
	vid, err := strconv.Atoi(id)
	if err != nil {
		log.Println("strconv error in barrage save4 :", err)
		return
	}
	if content == "" {
		ReturnError(c, 4001, "弹幕不能为空")
		return
	}
	if uid == 0 {
		ReturnError(c, 4002, "请先登陆")
		return
	}
	if eid == 0 {
		ReturnError(c, 4003, "请先选择剧集")
		return
	}
	if vid == 0 {
		ReturnError(c, 4005, "必须指定视频id")
		return
	}
	if ct == 0 {
		ReturnError(c, 4006, "必须指定当前视频时间")
		return
	}
	err = model.SaveBarrage(eid, vid, ct, uid, content)
	if err != nil {
		ReturnError(c, 5000, "保存失败请重试")
		return
	}
	ReturnSuccess(c, 0, "success", "", 0)
}
