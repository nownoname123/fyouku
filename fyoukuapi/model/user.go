package model

import (
	"encoding/json"
	"fyoukuapi/dao"
	"fyoukuapi/serve"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Usertype struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Id       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Status   int
	AddTime  int64
	Avatar   string //头像
	Name     string
}
type UserInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	AddTime int64  `json:"addTime"`
	Avatar  string `json:"avatar"`
}

func IsUserMobile(mobile string) (bool, error) {
	db := dao.Db

	err := db.AutoMigrate(&Usertype{})
	if err.Error != nil {

		log.Println("数据库迁移错误:", err.Error)
		return false, err.Error
	}
	var user []*Usertype
	err1 := db.Where("Mobile = ?", mobile).Find(&user)
	if err1.Error != nil {
		log.Println("IsUserMobile error", err1.Error)
		return false, err1.Error
	}
	if len(user) != 0 {
		return false, nil
	}
	return true, nil
}
func UserSave(mobile string, password string) error {
	st, err11 := HashPassword(password)
	if err11 != nil {
		log.Println("密码加密出错", err11)
		return err11
	}
	db := dao.Db
	err := db.AutoMigrate(&Usertype{})
	if err.Error != nil {

		log.Println("数据库迁移错误in usersave")
		return err.Error
	}
	ntr := rand.Intn(5)
	var head [6]string //头像地址

	head[0] = "https://cn.bing.com/images/search?view=detailV2&ccid=8Ka%2fQr0U&id=CD9C6AB1F9A394EB6BBB1E56267CC8A380700988&thid=OIP.8Ka_Qr0UZpBWSu1UF5RcQAHaHa&mediaurl=https%3a%2f%2fpic2.zhimg.com%2fv2-2d44b34343fadb3f01872fa244580bc1_r.jpg&exph=700&expw=700&q=%e4%ba%8c%e6%ac%a1%e5%85%83%e5%a4%b4%e5%83%8f%e5%a4%a7%e5%85%a8&simid=608008520324423137&FORM=IRPRST&ck=7C4805037328343D3156714186512B3E&selectedIndex=9"
	head[1] = "https://cn.bing.com/images/search?view=detailV2&ccid=ZeHJt4kY&id=88B1F3351D0F698B3806D2DE980A910EBCB7B01D&thid=OIP.ZeHJt4kYjExfeHk0cSEybwHaHa&mediaurl=https%3a%2f%2fimg.zcool.cn%2fcommunity%2f01e8ce5d89c4eaa8012060bed28708.png%402o.png&exph=2000&expw=2000&q=%e4%ba%8c%e6%ac%a1%e5%85%83%e5%a4%b4%e5%83%8f%e5%a4%a7%e5%85%a8&simid=607995360541499991&FORM=IRPRST&ck=F788604851AF30E61456499FB5A0A59F&selectedIndex=16"
	head[2] = "https://cn.bing.com/images/search?view=detailV2&ccid=kwBQODCZ&id=39EA1E612A979AC08F88B3DA09A1AF069DDE0D14&thid=OIP.kwBQODCZLDq2v2XpXcGxsQHaHa&mediaurl=https%3a%2f%2fpic1.zhimg.com%2fv2-7dc0fb227cc67e12f186b19858356567_r.jpg%3fsource%3d1940ef5c&exph=1080&expw=1080&q=%e4%ba%8c%e6%ac%a1%e5%85%83%e5%a4%b4%e5%83%8f%e5%a4%a7%e5%85%a8&simid=608046973159152076&FORM=IRPRST&ck=0A78B7F2274110AE0B09C655075037FA&selectedIndex=32"
	head[3] = "https://ts1.cn.mm.bing.net/th/id/R-C.abd68212f721eb2546840ba4735aa562?rik=WRkygXN7Qfgfeg&riu=http%3a%2f%2fimg4.a0bi.com%2fupload%2fttq%2f20200724%2f1595558809098.jpg%3fimageView2%2f0%2fw%2f600%2fh%2f800&ehk=OVWrnvEj%2b3opKgE9YBHy7BpWkSmaE0P0Lz2HTeijO7o%3d&risl=&pid=ImgRaw&r=0"
	head[4] = "https://img.zcool.cn/community/01e8ce5d89c4eaa8012060bed28708.png@2o.png"
	head[5] = "https://pic4.zhimg.com/v2-77fa45645098cec0aa17f57daffd7881_r.jpg"
	var user Usertype
	user.Password = st
	user.Name = ""
	user.AddTime = time.Now().Unix()
	user.Mobile = mobile
	user.Status = 1
	user.Avatar = head[ntr]
	err1 := db.Create(&user)
	if err1.Error != nil {
		return err1.Error
	}
	return nil
}
func UserLogin(mobile string, password string) (uid int64, uname string, err error) {
	db := dao.Db
	err1 := db.AutoMigrate(&Usertype{})
	if err1.Error != nil {

		log.Println("数据库迁移错误in userLogin:", err1.Error)
		return 0, "", err1.Error
	}
	var user Usertype
	err1 = db.Where("Mobile = ?", mobile).First(&user)
	if err1.Error != nil {
		log.Println("数据库读取错误in userlogin:", err1.Error)
		return 0, "", err1.Error
	}
	hashPassword := user.Password
	if ComparePasswords(hashPassword, password) {
		return user.Id, user.Name, nil
	}
	return 0, "密码错误", nil
}

/*func GetUserInfo(uid int) (UserInfo, error) {

db := dao.Db
	db.AutoMigrate(&Usertype{})
	var user Usertype
	var u UserInfo
	result :=	 db.Where("id = ?", uid).First(&user)
	if result.Error != nil {
		log.Println("user info error in db :", result.Error)
		return u, result.Error
	}

	u.Id = int(user.Id)
	u.Name = user.Name
	u.AddTime = user.AddTime
	u.Avatar = user.Avatar
	return u, nil

}*/

func GetUserInfo(uid int) (UserInfo, error) {
	rdb := serve.Rdb
	ctx := serve.Rctx
	redisKey := "user:id" + strconv.Itoa(uid)
	exists, err := rdb.Exists(ctx, redisKey).Result()

	if exists == 1 {
		var user UserInfo
		res, _ := rdb.Get(ctx, redisKey).Result()
		err = json.Unmarshal([]byte(res), &user)
		return user, err
	}
	db := dao.Db
	db.AutoMigrate(&Usertype{})
	var user Usertype
	var u UserInfo
	result := db.Where("id = ?", uid).First(&user)
	if result.Error != nil {
		log.Println("user info error in db :", result.Error)
		return u, result.Error
	}

	u.Id = int(user.Id)
	u.Name = user.Name
	u.AddTime = user.AddTime
	u.Avatar = user.Avatar
	jsonData, err := json.Marshal(u)
	err = rdb.Set(ctx, redisKey, jsonData, 24*time.Hour).Err()

	return u, err

}
