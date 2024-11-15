package router

import (
	"github.com/gin-gonic/gin"
	"travel/controller"
	"travel/middleware"
)

func NewRouter(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RecoveryMiddleware())
	r.POST("/travel/login", controller.Login)
	r.POST("travel/GetUserProfile", controller.GetUserProfile)

	PostRouters := r.Group("/post")
	PostRouters.Use(middleware.AuthMiddleware())
	postController := controller.NewIPostController()
	PostRouters.POST("/create", postController.Create)
	PostRouters.PUT("/update/:id", postController.Update)
	PostRouters.GET("/show/:id", postController.Show)
	PostRouters.DELETE("/delete/:id", postController.Delete)
	PostRouters.GET("page/list", postController.PageList)

	return r
}
