package router

import (
	"fyoukuapi/controllers"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.POST("/register/save", controllers.UserRegister)
	r.POST("/login/do", controllers.UserLogin)
	r.GET("/channel/advert", controllers.ChannelAdvert)
	r.GET("/channel/hot", controllers.ChannelHotList)
	r.GET("/channel/recommend/region", controllers.Recommend)
	r.GET("/channel/recommend/type", controllers.RecommendByType)
	r.GET("/channel/region", controllers.ChannelRegion)
	r.GET("/channel/type", controllers.ChannelType)
	r.GET("/channel/video", controllers.ChannelVideo)
	r.GET("/video/info", controllers.VideoInfo)
	r.GET("/video/episodes/list", controllers.VideoEpisodesList)
	r.GET("/comment/list", controllers.CommentList)
	r.POST("/comment/save", controllers.CommentSave)
	r.GET("/channel/top", controllers.ChannelTop)
	r.GET("/type/top", controllers.TypeTop)
	r.POST("/send/message", controllers.Message)
	r.GET("/barrage/ws", controllers.BarrageWs)
	r.POST("barrage/save", controllers.BarrageSave)
	r.GET("/user/video", controllers.UserVideo)
	r.POST("/video/save", controllers.VideoSave)
	r.GET("/video/send/es", controllers.SendEs)
	r.POST("video/search", controllers.Search)
	return r
}
