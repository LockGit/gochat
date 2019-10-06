/**
 * Created by lock
 * Date: 2019-10-06
 * Time: 23:09
 */
package router

import (
	"github.com/gin-gonic/gin"
	"gochat/api/handler"
	"gochat/config"
	"gochat/tools"
)

func Register() *gin.Engine {
	r := gin.Default()
	initUserRouter(r)
	initPushRouter(r)
	r.NoRoute(func(c *gin.Context) {
		tools.FailWithMsg(c, "please check request url !")
	})
	return r
}

func initUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.POST("/login", handler.Login)
	userGroup.POST("/register", handler.Register)
	userGroup.Use(CheckSessionId())
	{
		userGroup.POST("/checkAuth", handler.CheckAuth)
		userGroup.POST("/logout", handler.Logout)
	}

}

func initPushRouter(r *gin.Engine) {
	pushGroup := r.Group("/push")
	pushGroup.Use(CheckSessionId())
	{
		pushGroup.POST("/push", handler.Push)
		pushGroup.POST("/pushRoom", handler.PushRoom)
		pushGroup.POST("/count", handler.Count)
		pushGroup.POST("/getRoomInfo", handler.GetRoomInfo)
	}

}

func CheckSessionId() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.GetMode() == "dev" {
			c.Next()
			return
		}
		c.Next()
		return
	}
}
