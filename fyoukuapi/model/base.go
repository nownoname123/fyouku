package model

import (
	"fyoukuapi/dao"
	"log"
)

type ChannelRegion struct {
	Id   int
	Name string
}
type ChannelType struct {
	Id   int
	Name string
}

func GetChannelRegion(cid int) (int, []ChannelRegion, error) {
	db := dao.Db
	var region []ChannelRegion
	err1 := db.AutoMigrate(&ChannelRegion{})
	if err1.Error != nil {
		log.Println("数据库关联错误in GEtChannelRegion:", err1.Error)
		return 0, region, err1.Error
	}

	err := db.Where("channel_id =?  and status = ?", cid, 1).Order("sort desc").Find(&region)
	if err.Error != nil {
		log.Println("查询错误in GEtChannelRegion:", err1.Error)
		return 0, region, err.Error
	}
	return len(region), region, nil
}

func GetChannelType(cid int) (int, []ChannelType, error) {
	db := dao.Db
	var Type []ChannelType
	err1 := db.AutoMigrate(&ChannelType{})
	if err1.Error != nil {
		log.Println("数据库关联错误in GEtChannelType:", err1.Error)
		return 0, Type, err1.Error
	}

	err := db.Where("channel_id =?  and status = ?", cid, 1).Order("sort desc").Find(&Type)
	if err.Error != nil {
		log.Println("查询错误in GEtChannelType:", err1.Error)
		return 0, Type, err.Error
	}
	return len(Type), Type, nil
}

func GetChannelVideo(cid int, rid int, tid int, end string, sort string, limit int, offset int) (int, []Video, error) {
	db := dao.Db
	query := db.Model(&Video{})
	query = query.Where("channel_id =?", cid)
	query = query.Where("status =?", 1)
	if rid > 0 {
		query = query.Where("region_id =?", rid)
	}
	if tid > 0 {
		query = query.Where("type_id =?", tid)
	}
	if end == "n" {
		query = query.Where("is_end =?", 0)
	} else if end == "y" {
		query = query.Where("is_end =?", 1)
	}
	if sort == "episodesUpdateTime" {
		query = query.Order("episodes_update_time desc")
	} else if sort == "comment" {
		query = query.Order("comment desc")
	} else {
		query = query.Order("add_time desc")
	}
	query.Limit(limit).Offset(offset)
	var video []Video
	result := query.Find(&video)
	if result.Error != nil {
		log.Println("error in channel video db find :", result.Error)
		return 0, video, result.Error
	}
	return len(video), video, nil
}
