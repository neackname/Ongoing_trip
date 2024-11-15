package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"travel/MySQLTavelDate"
	"travel/TravelModel"
)

// @title	AuthMiddleware
// /@description	鉴权中间件，鉴定用户权限，鉴定用户是否登录
// @auth	Snactop	2023-11-30	19:54
// @param	无传入参数
// @return	gin.HandlerFunc	返回一个请求处理函数
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取 authorization Header
		tokenString := ctx.GetHeader("Authorization")

		//valid token format
		if len(tokenString) == 0 || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": "401", "msg": "权限不足"})
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:]

		token, claim, err := MySQLTavelDate.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": "401", "msg": "权限不足"})
			ctx.Abort()
			return
		}

		//验证通过获取后的userId
		user := TravelModel.TraUser{}
		userId := claim.UserId
		if userId == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": "401", "msg": "权限不足"})
			ctx.Abort()
			return
		}

		db := MySQLTavelDate.GetDB()
		if err := db.First(&user, userId).Error; err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": "401", "msg": "权限不足"})
			ctx.Abort()
			return
		}

		authInfo := TravelModel.AuthInformation{
			ID:         user.ID,
			OpenID:     user.OpenID,
			SessionKey: user.SessionKey,
		}

		//将用户信息写入上下文
		ctx.Set("authInfo", authInfo)
		ctx.Next()
	}
}
