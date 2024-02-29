package main

//接受者的主程序，实现程序的解耦
import (
	"encoding/json"
	"fyoukuapi/model"
	"github.com/go-redis/redis/v8"

	"fyoukuapi/pkg/logger"
	"fyoukuapi/serve"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

const (
	MysqlDb = "test:123456@tcp(127.0.0.1:3306)/fyouku?charset=utf8"
)

var (
	Db  *gorm.DB
	err error
)

func main() {
	Db, err = gorm.Open("mysql", MysqlDb)
	Db.SingularTable(true)
	if err != nil {
		logger.Error(map[string]interface{}{"mysql connect error": err.Error()})
	}
	if Db.Error != nil {
		logger.Error(map[string]interface{}{"database error": Db.Error})
	}

	// ----------------------- 连接池设置 -----------------------
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	Db.DB().SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	Db.DB().SetMaxOpenConns(100)

	serve.Consumer("", "fyouku_top", callback)
}
func callback(s string) {
	type Data struct {
		VideoId int
	}
	var data Data
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		log.Println("json error in mq top json :", err)
	}
	videoInfo, err := model.RedisGetVideoInfo(data.VideoId)
	conn := serve.Rdb
	defer func(conn *redis.Client) {
		_ = conn.Close()
	}(conn)
	redisChannelKey := "video:top:channel:channelId:" + strconv.Itoa(videoInfo.ChannelId)
	redisTypeKey := "video:top:type:typeid:" + strconv.Itoa(videoInfo.TypeId)
	ctx := serve.Rctx
	conn.ZIncrBy(ctx, redisChannelKey, 1, strconv.Itoa(data.VideoId))
	conn.ZIncrBy(ctx, redisTypeKey, 1, strconv.Itoa(data.VideoId))
	log.Println("msg is :", s)

}
