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

	vi := r.Group("/travel")
	vi.Use(middleware.AuthMiddleware()) // 应用JWT认证中间件
	{
		//用户
		UserControllers := vi.Group("/user")
		UserControllers.PATCH("/update", controller.Update)
		UserControllers.GET("/info", controller.GetUserInformation)
		UserControllers.GET("/postCreate", controller.GetUserCreatedPosts)

		//用户收藏文章
		UserStartControllers := vi.Group("/user/start")
		UserStartControllers.POST("/add/:id", controller.AddPostStart)
		UserStartControllers.DELETE("/remove/:id", controller.RemovePostStart)
		UserStartControllers.GET("/list", controller.GetPostStart) // TODO 获取用户收藏列表，这里和PageList有点像，但是这里好像有点问题，之后解决一下

		//文章路由
		PostRouters := vi.Group("/post")
		PostRouters.POST("/create", controller.PostCreate)
		PostRouters.PATCH("/update/:id", controller.PostUpdate)
		PostRouters.GET("/show/:id", controller.PostShow)
		PostRouters.DELETE("/delete/:id", controller.PostDelete)
		PostRouters.GET("/page/list", controller.PostPageList)

	}

	return r
}
