package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"regexp"
	"travel/MySQLTavelDate"
	"travel/TravelModel"
	"travel/logic"
	"travel/vo"
)

func GetSessionId(ctx *gin.Context) {
	db := MySQLTavelDate.GetDB()
	wxCode := ctx.PostForm("code")

	//参数验证
	if validInputPattern := regexp.MustCompile(`^[a-zA-Z0-9_]+$`); wxCode == "" || !validInputPattern.MatchString(wxCode) {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	//向系统获取appID、appSecret
	var appID = viper.GetString("wx.appID")
	var appSecret = viper.GetString("wx.appSecret")
	if appID == "" || appSecret == "" {
		fmt.Println("appID或appSecret错误， 参数不能为空")
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统出错"})
		return
	}

	//向微信API发送请求，获取用户openID
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appID, appSecret, wxCode)
	resp, err1 := http.Get(url)
	if err1 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "微信系统出错"})
		return
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	var sessionResponse vo.Code2SessionResponse
	if err := json.Unmarshal(body, &sessionResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	//检查此openID的用户是否注册
	var user TravelModel.TraUser
	db.Where("  OpenID  = ?", sessionResponse.OpenID).First(&user)
	if user.ID == 0 {
		//用户未注册，为用户注册，创建用户信息
		userInfo := TravelModel.TraUser{
			OpenID:     sessionResponse.OpenID,
			SessionKey: sessionResponse.SessionKey,
		}
		db.Create(&userInfo)
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "SessionKey": sessionResponse.SessionKey, "msg": "登录成功"})
	return
}

func AuthLogin(ctx *gin.Context) {
	var identifyCode vo.IdentifyCode
	if err := ctx.ShouldBind(&identifyCode); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	//解密数据
	var userInfo TravelModel.TraUser
	if err := logic.DecryptUserInfo(&userInfo, identifyCode.SessionKey, identifyCode.EncryptedData, identifyCode.Iv); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误"})
		return
	}

	//验证openID是否一致
	var user TravelModel.TraUser
	db := MySQLTavelDate.GetDB()

	if err := db.Where(" SessionKey = ?", identifyCode.SessionKey).First(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "数据库错误"})
		return
	}
	if user.OpenID != userInfo.OpenID {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "OpenID不一致"})
		return
	}

	token, err := MySQLTavelDate.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "token发放错误"})
		return
	}

	//更新用户信息
	userInfo.SessionKey = identifyCode.SessionKey
	if err := db.Model(&user).Updates(userInfo).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户信息更新失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "token": token, "msg": "登录成功"})
	return
}
