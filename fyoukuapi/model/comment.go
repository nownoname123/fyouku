package model

import (
	"fyoukuapi/dao"
	"log"
	"time"
)

type Comment struct {
	Id          int
	Content     string
	AddTime     int64
	UserId      int
	Stamp       int
	Status      int
	PraiseCount int
	EpisodesId  int
	VideoId     int
}

func GetCommentList(eid int, offset int, limit int) (int, []Comment, error) {
	db := dao.Db
	db.AutoMigrate(&Comment{})
	var comment []Comment
	err := db.Where("episodes_id = ?", eid).Limit(limit).Offset(offset).Order("add_time desc").Find(&comment)
	if err.Error != nil {
		log.Println("comment ist error:", err.Error)
		return 0, comment, err.Error
	}
	return len(comment), comment, nil
}

func SaveComment(content string, uid int, eid int, vid int) error {
	db := dao.Db
	db.AutoMigrate(&Comment{})
	var comment Comment
	comment.Content = content
	comment.UserId = uid
	comment.VideoId = vid
	comment.AddTime = time.Now().Unix()
	comment.Stamp = 0
	comment.Status = 1
	comment.EpisodesId = eid
	comment.PraiseCount = 0

	result := db.Create(&comment)
	if result.Error == nil {
		AddVideoComment(vid)
		AddEpisodesComment(eid)
	}
	return result.Error
}
