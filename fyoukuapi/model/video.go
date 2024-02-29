package model

import (
	"encoding/json"
	"fyoukuapi/dao"
	"fyoukuapi/serve"
	"log"
	"strconv"
	"time"
)

type Video struct {
	Id                 int    `redis:"id" json:"id"`
	Title              string `redis:"title" json:"title"`
	SubTitle           string `redis:"sub_title" json:"subTitle"`
	AddTime            int64  `redis:"add_time" json:"addTime"`
	Img                string `redis:"img" json:"img"`
	Img1               string `redis:"img1" json:"img1"`
	EpisodesCount      int    `redis:"episodes_count" json:"episodesCount"`
	IsEnd              int    `redis:"is_end" json:"isEnd"`
	Comment            int    `json:"comment"`
	ChannelId          int    `json:"channelID"`
	TypeId             int    `json:"typeId"`
	RegionId           int    `json:"regionId"`
	UserId             int    `json:"userId"`
	Status             int    `json:"status"`
	EpisodesUpdateTime int    `json:"episodesUpdateTime"`
	IsRecommend        int    `json:"isRecommend"`
}
type Video1 struct {
	Id                 int
	Title              string
	SubTitle           string
	AddTime            int64
	Img                string
	Img1               string
	EpisodesCount      int
	IsEnd              int
	ChannelId          int
	RegionId           int
	TypeId             int
	UserId             int
	Comment            int
	EpisodesUpdateTime int64
}

// GetChannelHotList 增加redis缓存-获取视频详情
func GetChannelHotList(cid int) (int, []Video, error) {
	db := dao.Db
	var video []Video
	err1 := db.AutoMigrate(&Video{})
	if err1.Error != nil {
		log.Println("数据库关联错误in GEtChannelHotList:", err1.Error)
		return 0, video, err1.Error
	}

	err := db.Where("channel_id =? and is_hot= ? and status = ?", cid, 1, 1).Order("episodes_" +
		"update_time desc").Limit(9).Find(&video)
	if err.Error != nil {
		log.Println("查询错误in GEtChannelHotList:", err1.Error)
		return 0, video, err.Error
	}
	return len(video), video, nil
}

func GetChannelRecommendRegionList(cid int, rid int) (int, []Video, error) {
	db := dao.Db
	err1 := db.AutoMigrate(&Video{})
	var video []Video
	if err1.Error != nil {
		log.Println("数据库关联错误in 2", err1.Error)
		return 0, video, err1.Error
	}
	err := db.Where("channel_id= ? and region_id= ? and is_recommend= ? and status= ?", cid, rid, 1, 1).Order("" +
		"episodes_update_time desc").Limit(9).Find(&video)
	if err.Error != nil {
		log.Println("查询错误in GEtChannelHotList:", err.Error)
		return 0, video, err.Error
	}
	return len(video), video, nil
}

func GetChannelRecommendTypeList(cid int, tid int) (int, []Video, error) {
	db := dao.Db
	err1 := db.AutoMigrate(&Video{})
	var video []Video
	if err1.Error != nil {
		log.Println("数据库关联错误in 3", err1.Error)
		return 0, video, err1.Error
	}
	err := db.Where("channel_id= ? and type_id= ? and is_recommend= ? and status= ?", cid, tid, 1, 1).Order("" +
		"episodes_update_time desc").Limit(9).Find(&video)
	if err.Error != nil {
		log.Println("查询错误in GEtChannelHotList:", err1.Error)
		return 0, video, err.Error
	}
	return len(video), video, nil
}

/*func GetVideoInfo(vid int) (Video, error) {
	db := dao.Db
	db.AutoMigrate(&Video{})
	var video Video
	result := db.Where("id =?", vid).First(&video)
	return video, result.Error
}*/

// RedisGetVideoInfo GetChannelRecommendRegionList 获取推荐视频根据频道id
func RedisGetVideoInfo(vid int) (Video, error) {
	var ctx = serve.Rctx
	var video Video
	conn := serve.Rdb
	//log.Println("redis start")
	redisKey := "video:id:" + strconv.Itoa(vid)

	//判断redis中是否存在
	exists, err := conn.Exists(ctx, redisKey).Result()
	if err != nil {
		log.Println("redis error1 :", err)

	}
	if exists == 1 {
		res, _ := conn.Get(ctx, redisKey).Result()
		err = json.Unmarshal([]byte(res), &video)

		//	log.Println("redis 有，从redis中拿的:", video.Id)
	} else {
		db := dao.Db
		db.AutoMigrate(&Video{})
		result := db.Where("id =?", vid).First(&video)
		if result.Error == nil {
			//保存redis
			jsonData, _ := json.Marshal(video)
			if err := conn.Set(ctx, redisKey, jsonData, 24*time.Hour).Err(); err != nil {
				log.Println("redis error:", err)
			}

			//	log.Println("redis 无，从数据库中拿的:", video.Id, " error is :", err)
		}
	}

	return video, err
}

type VideoEpisodes struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	AddTime int64  `json:"addTime"`
	Num     int    `json:"num"`
	PlayUrl string `json:"playUrl"`
	Comment int    `json:"comment"`

	Status  int         `json:"status"`
	VideoId interface{} `json:"videoId"`
}
type Ve struct {
	Id      int
	Title   string
	AddTime int64
	Num     int
	PlayUrl string
	Comment int

	Status  int
	VideoId int
}

/*func GetVideoEpisodesList(vid int) (int, []VideoEpisodes, error) {
	db := dao.Db
	db.AutoMigrate(&VideoEpisodes{})
	var video []VideoEpisodes
	err := db.Where("video_id= ?", vid).Find(&video)
	if err.Error != nil {
		log.Println("查询错误in GEtChannelHotList:", err.Error)
		return 0, video, err.Error
	}
	return len(video), video, nil
}*/

func GetVideoEpisodesList(vid int) (int, []VideoEpisodes, error) {
	rdb := serve.Rdb
	ctx := serve.Rctx
	var video []VideoEpisodes

	redisKey := "video:episodes:videoid:" + strconv.Itoa(vid)
	exist, err := rdb.Exists(ctx, redisKey).Result()
	if err != nil {
		log.Println("episodes error :", err)
	}
	if exist == 1 {
		num, err := rdb.LLen(ctx, redisKey).Result()
		var episodes VideoEpisodes
		if err == nil {
			res, _ := rdb.LRange(ctx, redisKey, 0, -1).Result()
			for _, v := range res {
				err = json.Unmarshal([]byte(v), &episodes)
				if err == nil {
					video = append(video, episodes)
				}
			}
		}
		return int(num), video, err
	}
	db := dao.Db
	db.AutoMigrate(&VideoEpisodes{})

	err1 := db.Where("video_id= ?", vid).Order("num asc").Find(&video)
	if err1.Error != nil {
		log.Println("查询错误in GEtChannelHotList:", err1.Error)
		return 0, video, err1.Error
	}
	for _, v := range video {
		jv, _ := json.Marshal(v)
		rdb.RPush(ctx, redisKey, jv)
	}
	return len(video), video, nil
}

type VideoComment struct {
	id      int
	Comment int
}

func AddVideoComment(vid int) {
	db := dao.Db
	db = db.Table("video")
	var video VideoComment
	result := db.Where("id =?", vid).First(&video)
	if result.Error != nil {
		log.Println("查询失败1:", result.Error)
		return
	}
	video.Comment++
	db.Save(video)

	//更新redis排行榜，通过MQ来实现

	videoObj := map[string]int{
		"VideoId": vid,
	}
	videoJson, _ := json.Marshal(videoObj)
	err := serve.Publish("", "fyouku_top", string(videoJson))
	if err != nil {
		log.Println("mq error:", err)
	}
}
func AddEpisodesComment(eid int) {
	db := dao.Db
	db = db.Table("video_episodes")
	var video VideoComment
	result := db.Where("id =?", eid).First(&video)
	if result.Error != nil {
		log.Println("查询失败1:", result.Error)
		return
	}
	video.Comment++
	db.Save(video)
}

func GetUserVideo(uid int) (int, []Video, error) {
	db := dao.Db
	var video []Video
	db.AutoMigrate(&Video{})
	result := db.Where("user_id = ?", uid).Order("add_time desc").Find(&video)
	return len(video), video, result.Error

}

func SaveVideo(title string, subTitle string, cid int, rid int, tid int, playUrl string, uid int) error {
	db := dao.Db
	db = db.Table("video")
	log.Println(cid, "  ", rid, "  ", tid, "  ", uid)
	var video1 Video1
	tm := time.Now().Unix()
	video1.AddTime = tm
	video1.Title = title
	video1.IsEnd = 1
	video1.EpisodesCount = 1
	video1.SubTitle = subTitle
	video1.UserId = uid
	video1.ChannelId = cid
	video1.TypeId = tid
	video1.RegionId = rid
	video1.Comment = 1
	video1.EpisodesUpdateTime = tm
	result := db.Create(&video1)
	if result.Error != nil {
		log.Println("error in create video:", result.Error)
		return result.Error
	}
	db.Where("add_time =? and user_id = ? and region_id = ?", tm, uid, rid).First(&video1)
	db2 := dao.Db
	db2 = db2.Table("video_episodes")
	var videoEpisodes Ve
	videoEpisodes.AddTime = tm
	videoEpisodes.Title = subTitle
	videoEpisodes.PlayUrl = playUrl
	videoEpisodes.Num = 1
	videoEpisodes.Comment = 0
	videoEpisodes.Status = 1
	videoEpisodes.VideoId = video1.Id
	result = db2.Create(&videoEpisodes)
	if result.Error != nil {
		log.Println("error in create video:", result.Error)
		return result.Error
	}
	return nil
}

// GetAllList 获取所有视频数据
func GetAllList() (int64, []Video, error) {
	db := dao.Db
	db = db.Table("video")
	var videos []Video
	result := db.Find(&videos)
	if result.Error != nil {
		log.Println("db error in Get all list:", result.Error)
		return 0, videos, result.Error
	}
	num := len(videos)
	k := int64(num)
	return k, videos, nil
}
