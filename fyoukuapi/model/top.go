package model

import (
	"fmt"
	"fyoukuapi/dao"
	"fyoukuapi/serve"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

/*func GetChannelTop(cid int) (int, []Video, error) {
	db := dao.Db
	db.AutoMigrate(&Video{})
	var video []Video
	result := db.Where("channel_id = ?", cid).Limit(10).Order("comment desc").Find(&video)
	return len(video), video, result.Error

}*/

func GetChannelTop(cid int) (int, []Video, error) {
	rdb := serve.Rdb
	ctx := serve.Rctx
	var video []Video
	log.Println("Get channel top running")
	var redisKey = "video:top:channel:channelId" + strconv.Itoa(cid)
	exist, err := rdb.Exists(ctx, redisKey).Result()
	if err != nil {
		log.Println("error in rdb exists:", err)
	}
	log.Println("redisKey is :", redisKey)
	if exist == 1 {
		var num int64
		log.Println("排行榜中有redis key:", redisKey)
		num, err = rdb.ZCard(ctx, redisKey).Result()
		if err != nil {
			log.Println("z card error :", err)
		}
		log.Println("找到了：", num)
		/*	op := redis.ZRangeBy{

			Offset: 0,  // 类似sql的limit, 表示开始偏移量
			Count:  10, // 一次返回多少数据
		}*/

		res, err := rdb.ZRevRange(ctx, redisKey, 0, -1).Result()
		if err != nil {
			log.Println("redis error in rdb range :", err)
		}
		log.Println("res = ", res)
		for _, v := range res {
			fmt.Println(v)
			vid, err := strconv.Atoi(v)
			if err != nil {
				log.Println("atoi error :", err)
			}
			val, err := RedisGetVideoInfo(vid)
			if err == nil {
				video = append(video, val)
			} else {
				log.Println("range redis error:", err)
			}
		}
		k := int(num)
		return k, video, nil
	} else {
		db := dao.Db
		db.AutoMigrate(&Video{})

		result := db.Where("channel_id = ?", cid).Limit(10).Order("comment desc").Find(&video)
		if result.Error != nil {
			log.Println("db error in channel top :", result.Error)
		}
		for _, v := range video {
			score := float64(v.Comment)
			rdb.ZAdd(ctx, redisKey, &redis.Z{Score: score, Member: v.Id})
			rdb.Expire(ctx, redisKey, 24*time.Hour)
		}
		return len(video), video, result.Error
	}

}

func GetTypeTop(tid int) (int, []Video, error) {
	db := dao.Db
	db.AutoMigrate(&Video{})
	var video []Video
	result := db.Where("type_id = ?", tid).Limit(10).Order("comment desc").Find(&video)
	return len(video), video, result.Error

}
