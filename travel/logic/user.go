package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"travel/vo"
)

// 使用code获取sessionkey和openID的函数
func GetIdentify(wxCode string) (key, ID string, err error) {
	//向系统获取appID、appSecret
	var appID = viper.GetString("wx.appID")
	var appSecret = viper.GetString("wx.appSecret")
	if appID == "" || appSecret == "" {
		errMsg := "服务器系统错误"
		return "", "", errors.New(errMsg)
	}

	//向微信API发送请求，获取用户openID
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appID, appSecret, wxCode)
	resp, err1 := http.Get(url)
	if err1 != nil {
		errMsg := "系统错误"
		return "", "", errors.New(errMsg)
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		errMsg := "系统错误"
		return "", "", errors.New(errMsg)
	}

	var sessionResponse vo.Code2SessionResponse
	if err := json.Unmarshal(body, &sessionResponse); err != nil {
		errMsg := "系统错误"
		return "", "", errors.New(errMsg)
	}

	return sessionResponse.SessionKey, sessionResponse.OpenID, nil
}
