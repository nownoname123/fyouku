package model

import (
	"fyoukuapi/dao"
	"log"
)

type Advert struct {
	Id       int `gorm:"primaryKey;autoIncrement"`
	Title    string
	SubTitle string
	AddTime  int64
	Img      string
	Url      string
}

func GetChannelAdvert(channelId int) (int, []Advert, error) {
	db := dao.Db
	db.AutoMigrate(&Advert{})
	var advert []Advert
	err1 := db.Where("Id = ? and status = ?", channelId, 1).Order("sort desc").Find(&advert)
	if err1.Error != nil {
		log.Println("advert error :", err1.Error)
		return 0, advert, err1.Error
	}
	return len(advert), advert, nil
}
