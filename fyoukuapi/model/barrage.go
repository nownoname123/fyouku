package model

import (
	"fyoukuapi/dao"
	"log"
	"time"
)

type Barrage struct {
	Id          int
	Content     string
	CurrentTime int
	AddTime     int64
	UserID      int
	Status      int
	EpisodesId  int
	VideoId     int
}

type BarrageData struct {
	Id          int    `json:"id"`
	Content     string `json:"content"`
	CurrentTime int    `json:"currentTime"`
}

func GetBarrageList(eid int, startTime int, endTime int) (int64, []BarrageData, error) {
	db := dao.Db
	db = db.Table("barrage")
	var barrages []BarrageData
	var barrage []Barrage

	//反引号用于将标识符（例如列名、表名）转义
	result := db.Where("episodes_id = ? AND `current_time`>= ? AND `current_time`<= ? ", eid, startTime, endTime).Order("`current_time` asc").Find(&barrage)
	var b BarrageData
	for _, k := range barrage {
		b.Id = k.Id
		b.Content = k.Content
		b.CurrentTime = k.CurrentTime
		barrages = append(barrages, b)
		//log.Println(b.CurrentTime, "-->", b.Content)
	}
	if result.Error != nil {
		log.Println("db error in GEtBarrageList:", result.Error)
		return 0, barrages, result.Error
	}
	k := int64(len(barrages))
	return k, barrages, nil
}
func SaveBarrage(eid int, vid int, ct int, uid int, content string) error {
	db := dao.Db
	db.AutoMigrate(&Barrage{})
	var barrage Barrage
	barrage.EpisodesId = eid
	barrage.AddTime = time.Now().Unix()
	barrage.Status = 1
	barrage.VideoId = vid
	barrage.UserID = uid
	barrage.CurrentTime = ct
	barrage.Content = content
	log.Println("收到弹幕：,", content)
	result := db.Create(&barrage)
	return result.Error
}
