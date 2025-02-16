package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"travel/TravelDate"
	"travel/TravelModel"
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

func GetUserInformation(OpenID string) (vo.ShowUserRequest, error) {
	user, err := TravelDate.GetUserInformation(OpenID)
	if err != nil {
		return vo.ShowUserRequest{}, err
	}

	userP := vo.ShowUserRequest{
		Telephone: user.Telephone,
		NickName:  user.NickName,
		Motto:     user.Motto,
		Gender:    user.Gender,
	}

	return userP, nil
}

func UpdateUserInformation(userUpdate vo.UpdateUserRequest, OpenID string) error {
	//查找用户是存在
	user, err := TravelDate.GetUserInformation(OpenID)
	if err != nil {
		return err
	}

	//更新用户数据
	if err := TravelDate.UpdateUserInformation(user, userUpdate); err != nil {
		return err
	}

	return nil
}

func AddFavoritePost(userID uint64, postID string) (code int, err error) {
	post, err := GetPostInfo(postID)
	if err != nil {
		code = http.StatusNotFound
		err = errors.New("商品不存在")
		return
	}

	// 检查用户是否已经收藏过该商品
	if exits := TravelDate.CheckUserPostStartExist(userID, post.ID); exits {
		code = http.StatusConflict
		return code, errors.New("该商品已收藏")
	}

	// 创建收藏记录
	if err := TravelDate.AddPostStart(userID, post.ID); err != nil {
		code = http.StatusInternalServerError
		return code, errors.New("收藏失败，请稍后重试")
	}

	return 200, nil
}

func RemovePostStart(userID uint64, postID string) (code int, err error) {
	post, err := GetPostInfo(postID)
	if err != nil {
		code = http.StatusNotFound
		err = errors.New("文章不存在")
		return
	}

	// 检查用户是否已经收藏过该商品
	if exits := TravelDate.CheckUserPostStartExist(userID, post.ID); !exits {
		code = http.StatusConflict
		err = errors.New("该文章未被收藏")
		return
	}

	if err := TravelDate.RemovePostStart1(userID, post.ID); err != nil {
		code = http.StatusInternalServerError
		return code, errors.New("取消收藏失败")
	}

	return 200, nil
}

func GetUserPostStart(userID uint64) (posts []TravelModel.Post, err error) {
	// 从数据库中获取用户收藏的商品ID列表
	postIDs, err := TravelDate.GetPostStartIDs(userID)
	if err != nil {
		return []TravelModel.Post{}, errors.New("获取收藏列表失败，请稍后重试")
	}

	// 如果用户没有收藏任何商品
	if len(postIDs) == 0 {
		return []TravelModel.Post{}, nil
	}

	// 查询所有收藏的商品信息
	posts, err = TravelDate.GetPostsByIDs(postIDs)
	if err != nil {
		return []TravelModel.Post{}, errors.New("获取文章信息失败，请稍后重试")
	}

	return posts, nil
}

func GetUserCreatedPosts(userID uint64) ([]TravelModel.Post, error) {
	commodities, err := TravelDate.GetUserCreatedPosts(userID)
	if err != nil {
		return []TravelModel.Post{}, err
	}

	return commodities, nil
}
