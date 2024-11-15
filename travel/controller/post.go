package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"travel/MySQLTavelDate"
	"travel/TravelModel"
	"travel/vo"
)

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

// @title	NewIPostController
// @description	control层创建文章接口
// @auth	Snactop	2023-12-1	19:53
// @param     void			没有入参
// @return	void  没有回参
func NewIPostController() IPostController {
	db := MySQLTavelDate.GetDB()
	db.AutoMigrate(&TravelModel.Post{})

	return PostController{DB: db}
}

// @title	Create
// @description	control层创建帖子
// @auth	Snactop	2023-12-1	19:53
// @param	ctx *gin.Context	传入一个上下文
// @return	void 无返回值
func (po PostController) Create(ctx *gin.Context) {
	//数据验证
	var p vo.PostRequest
	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "数据验证失败"})
		return
	}

	//创建文章
	user, _ := ctx.Get("authInfo")
	post := TravelModel.Post{
		UserId:  user.(TravelModel.AuthInformation).ID,
		Title:   p.Title,
		HeadImg: p.HeadImg,
		Content: p.Content,
	}

	if err := po.DB.Create(&post); err != nil {
		panic(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "文章创建成功"})
}

// @title	Update
// @description	control层更新帖子，此接口只适合放在用户创建的文章当中
// @auth	Snactop	2023-12-1	19:53
// @param	ctx *gin.Context	传入一个上下文
// @return	void     无返回值
func (po PostController) Update(ctx *gin.Context) {
	//TODO 数据验证
	var postR vo.PostRequest
	if err := ctx.ShouldBindJSON(&postR); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "数据验证失败"})
		return
	}
	//TODO 在path中寻找文章ID
	id, _ := strconv.Atoi(ctx.Param("id"))

	//TODO 在数据库中根据ID判断文章是否存在
	var post TravelModel.Post
	if err := po.DB.First(&post, id).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "文章不存在"})
		return
	}

	//TODO 重要：  判断删除者是否为文章作者
	user, exit := ctx.Get("authInfo")
	if !exit {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无更改权限"})
		return
	}
	if user.(TravelModel.AuthInformation).ID != post.UserId {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "文章不属于您， 操作非法"})
		return
	}

	//TODO 更新文章
	if err := po.DB.Model(&post).Updates(postR).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "文章更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "post": post, "msg": "更新成功"})
}

// @title	Show
// @description	control层查找帖子
// @auth	Snactop	2023-12-1	19:53
// @param	ctx *gin.Context	传入一个上下文
// @return	void   无返回值
func (po PostController) Show(ctx *gin.Context) {

}

func (po PostController) Delete(ctx *gin.Context) {
}

// @title	PageList
// @description	control层展示帖子列表
// @auth	Snactop	2023-12-1	19:53
// @param	ctx *gin.Context	传入一个上下文
// @return	void  无返回值
func (po PostController) PageList(ctx *gin.Context) {
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageList, _ := strconv.Atoi(ctx.DefaultQuery("pageList", "20"))

	var posts []TravelModel.Post
	po.DB.Order("created_at desc").Offset((pageNum - 1) * pageList).Limit(pageList).Find(&posts)

	//前端渲染分页需要知道总数
	var total int64
	po.DB.Model(TravelModel.Post{}).Count(&total)

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": posts, "total": total})
}
