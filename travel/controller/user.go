package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"travel/TravelDate"
	"travel/TravelModel"
	"travel/logic"
	"travel/pkg/jwt"
	"travel/pkg/snowflake"
	"travel/vo"
)

// Login @title  Login
// @description	调用方式：POST； 提交表单额方式：x-www-form-urlencoded；获取微信用户获取用户的openID和SessionKey以计入系统，
// @auth	Snactop	2023-11-27	20:07
// @param	ctx *gin.Context  传入一个上下文
// @return	void	没有返回值
func Login(ctx *gin.Context) {
	db := TravelDate.GetDB()
	wxCode := ctx.PostForm("code")

	//参数验证
	if validInputPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`); wxCode == "" || !validInputPattern.MatchString(wxCode) {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	SessionKey, OpenID, err := logic.GetIdentify(wxCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err})
		return
	}

	//检查此openID的用户是否注册
	var user TravelModel.TraUser
	db.Where("  open_id  = ?", OpenID).First(&user)
	if user.ID == 0 {
		//用户未注册，为用户注册，创建用户信息
		//TODO 生成用户ID
		userId, err := snowflake.GetID()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "系统繁忙，登录失败"})
			return
		}

		userInfo := TravelModel.TraUser{
			ID:         userId,
			OpenID:     OpenID,
			SessionKey: SessionKey,
		}
		db.Create(&userInfo)
	} else { //更新用户SessionKey
		if err := db.Model(&user).Update("session_key", SessionKey).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "系统错误，session_key更新失败"})
			return
		}
	}

	//发放token
	token, err := jwt.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "token发放错误"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "token": token, "SessionKey": SessionKey, "msg": "登录成功"})
	return
}

// @title  GetUserProfile
// @description	获取用户信息,暂时没有作用
// @auth	Snactop	2023-11-27	20:07
// @param	ctx *gin.Context  传入一个上下文（ EncryptedData、 Iv 、SessionKey）
// @return	void	没有返回值
func GetUserProfile(ctx *gin.Context) {
	var identifyCode vo.IdentifyCode
	if err := ctx.ShouldBind(&identifyCode); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	//todo 解密数据
	//var userInfo TravelModel.TraUser
	plainText, err := logic.DecryptUserInfo(identifyCode.SessionKey, identifyCode.EncryptedData, identifyCode.Iv)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "plainText": plainText, "msg": "操作成功"})
	return
}

// @title  Update
// @description	获取方式：POST；用户自行更改用户信息
// @auth	Snactop	2024-9-20	15:13
// @param	ctx *gin.Context  传入一个上下文
// @return	void	没有返回值
func Update(ctx *gin.Context) {
	var userP vo.UpdateUserRequest
	if err := ctx.ShouldBind(&userP); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求错误或数据格式错误"})
		return
	}

	//TODO 查找用户的openID
	authInfo, exit := ctx.Get("authInfo")
	if !exit {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "登录已过期，请重新登录"})
		return
	}

	if authInfo.(TravelModel.AuthInformation).OpenID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "登录已过期，请重新登录"})
		return
	}

	//更改用户信息
	if err := logic.UpdateUserInformation(userP, authInfo.(TravelModel.AuthInformation).OpenID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "msg": "操作成功"})
	return
}

// @title  GetUserInformation
// @description	获取方式：GET；展示用户信息用户信息
// @auth	Snactop	2024-9-20	15:13
// @param	ctx *gin.Context  传入一个上下文
// @return	void	没有返回值
func GetUserInformation(ctx *gin.Context) {
	//TODO 查找用户openID
	authInfo, exit := ctx.Get("authInfo")
	if !exit {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "登录已过期，请重新登录"})
		return
	}

	OpenID := authInfo.(TravelModel.AuthInformation).OpenID
	if OpenID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "登录已过期，请重新登录"})
		return
	}

	//查找用户信息
	user, err := logic.GetUserInformation(OpenID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "information": user, "msg": "获取用户信息成功"})
	return
}

// AddPostStart @title  AddPostStart
// @description	获取方式：POST；用户收藏文章
// @auth	Snactop	2024-9-20	15:13
// @param	ctx *gin.Context  传入一个上下文
// @return	void	没有返回值
func AddPostStart(ctx *gin.Context) {

	// 获取当前登录的用户ID（假设是通过JWT或Session获取的）
	authInfo, exists := ctx.Get("authInfo")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	// 获取文章ID
	postID := ctx.Param("id") // 从URL中获取 /users/favorites/:id
	if len(postID) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 逻辑层收藏商品
	code, err := logic.AddFavoritePost(authInfo.(TravelModel.AuthInformation).ID, postID)
	if err != nil {
		ctx.JSON(code, gin.H{"message": "收藏失败", "error": err.Error()})
		return
	}

	ctx.JSON(code, gin.H{"message": "收藏成功"})
	return

}

// RemovePostStart @title  RemovePostStart
// @description	获取方式：GET；用户删除收藏文章
// @auth	Snactop	2024-9-20	15:13
// @param	ctx *gin.Context  传入一个上下文
// @return	void	没有返回值
func RemovePostStart(ctx *gin.Context) {
	// 获取当前登录的用户ID
	authInfo, exists := ctx.Get("authInfo")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	// 获取商品ID
	postID := ctx.Param("id")
	if len(postID) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的商品ID"})
		return
	}

	// 删除用户的收藏记录
	code, err := logic.RemovePostStart(authInfo.(TravelModel.AuthInformation).ID, postID)
	if err != nil {
		ctx.JSON(code, gin.H{"message": "取消收藏失败", "error": err.Error()})
		return
	}

	ctx.JSON(code, gin.H{"message": "取消收藏成功"})
	return
}

func GetPostStart(ctx *gin.Context) {
	// 获取当前登录的用户ID
	authInfo, exists := ctx.Get("authInfo")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	//获取用户收藏列表
	posts, err := logic.GetUserPostStart(authInfo.(TravelModel.AuthInformation).ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "获取收藏失败", "err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": posts, "msg": "获取用户收藏列表成功", "err": nil})
	return
}

func GetUserCreatedPosts(ctx *gin.Context) {
	// 获取当前登录的用户ID
	authInfo, exists := ctx.Get("authInfo")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未登录"})
		return
	}

	posts, err := logic.GetUserCreatedPosts(authInfo.(TravelModel.AuthInformation).ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "获取收藏失败", "err": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": posts, "msg": "获取用户创建商品列表成功", "err": nil})
	return
}

// @title  Favorite
// @description	获取方式：POST；用户收藏足迹
// @auth	Snactop	2024-9-20	15:13
// @param	ctx *gin.Context  传入一个上下文
// @return	void	没有返回值
func FavoriteF(ctx *gin.Context) {
	//TODO 获取前端发送过来的足迹信息（不知道该以什么样的存储形式）

	//TODO 获取用户ID（系统发放的ID）

	//TODO 将足迹信息ID加入到用户文章收藏夹中

}

// @title  ShowUserFoot
// @description	获取方式：GET；展示用户历史足迹
// @auth	Snactop	2024-10-12	0:13
// @param	ctx *gin.Context  传入一个上下文
// @return	void	没有返回值
func ShowUserFoot(ctx *gin.Context) {

}

// @title  ShowP
// @description	获取方式：GET；展示用户收藏的文章信息
// @auth	Snactop	2024-9-20	15:13
// @param	ctx *gin.Context  传入一个上下文
// @return	void	没有返回值
func ShowP(ctx *gin.Context) {
	//TODO 获取用户ID

	//TODO 展示用户收藏的文章列表（要专门写一个文章列表的格式）
}

// @title  ShowF
// @description	获取方式：GET；展示用户收藏的足迹信息
// @auth	Snactop	2024-9-20	15:13
// @param	ctx *gin.Context  传入一个上下文
// @return	void	没有返回值
func ShowF(ctx *gin.Context) {
	//TODO 获取用户ID

	//TODO 展示用户收藏的足迹列表（要专门写一个足迹列表的格式）
}
