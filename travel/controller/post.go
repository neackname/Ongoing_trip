package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"travel/TravelModel"
	"travel/logic"
	"travel/vo"
)

func PostCreate(ctx *gin.Context) {
	//参数验证
	var postR vo.PostRequest
	if err := ctx.ShouldBind(&postR); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	//检查用户权限
	authInfo, exit := ctx.Get("authInfo")
	if !exit {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "登录已过期，请重新登录"})
		return
	}

	//逻辑层创建文章，返回文章ID
	postID, err := logic.PostCreate(authInfo.(TravelModel.AuthInformation).ID, postR)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统错误，文章创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "postID": postID, "msg": "文章创建成功"})
	return
}

func PostUpdate(ctx *gin.Context) {
	//参数验证
	var postR vo.PostRequest
	if err := ctx.ShouldBind(&postR); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	//检查用户权限
	authInfo, exit := ctx.Get("authInfo")
	if !exit {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "登录已过期，请重新登录"})
		return
	}

	//在path中寻找文章ID
	postID := ctx.Param("id")

	//用户修改文章信息
	post, err := logic.PostUpdate(authInfo.(TravelModel.AuthInformation).ID, postID, postR)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "文章修改失败", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "文章信息修改成功", "post": post})
	return
}

func PostShow(ctx *gin.Context) {
	//在path中寻找文章ID
	postID := ctx.Param("id")

	//逻辑层查找商品信息
	post, err := logic.GetPostInfo(postID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户权限错误", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "文章查询成功", "post": post})
	return
}

func PostDelete(ctx *gin.Context) {
	// 处理商品的删除逻辑
	//在path中寻找商品ID
	postID := ctx.Param("id")

	//检查用户权限
	authInfo, exit := ctx.Get("authInfo")
	if !exit {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "登录已过期，请重新登录"})
		return
	}

	//逻辑层删除商品信息
	if err := logic.PostDelete(authInfo.(TravelModel.AuthInformation).ID, postID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "系统错误", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除文章成功，Delete Post Success"})
}

func PostPageList(ctx *gin.Context) {
	// 处理获取商品列表的逻辑，通常需要分页和筛选
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageList, _ := strconv.Atoi(ctx.DefaultQuery("pageList", "20"))

	posts, total := logic.PageList(pageNum, pageList)

	data := map[string]interface{}{
		"msg":         "商品信息获取成功",
		"commodities": posts,
		"total":       total,
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "date": data, "message": "删除文章成功，Delete Post Success"})
	return
}
