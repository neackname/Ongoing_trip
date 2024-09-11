package router

import (
	"github.com/gin-gonic/gin"
	"travel/controller"
)

func NewRouter(r *gin.Engine) *gin.Engine {
	r.POST("/travel/login", controller.GetSessionId)
	r.POST("travel/register", controller.AuthLogin)

	return r
}
