package logic

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"travel/TravelDate"
	"travel/TravelModel"
	"travel/vo"
)

func PostCreate(userID uint64, p vo.PostRequest) (uuid.UUID, error) {
	//创建文章
	post := TravelModel.Post{
		UserID:  userID,
		Title:   p.Title,
		HeadImg: p.HeadImg,
		Content: p.Content,
	}

	postID, err := TravelDate.InsertPost(post)

	if err != nil {
		return uuid.Nil, errors.New("系统错误，文章创建失败")
	}

	return postID, nil
}

func PostUpdate(userID uint64, postID string, postR vo.PostRequest) (TravelModel.Post, error) {
	//查找commodityID是否存在
	post, err := TravelDate.GetPostInfo(postID)
	if err != nil {
		return TravelModel.Post{}, err
	}

	//判断是否为商品商家
	if userID != post.UserID {
		err := errors.New("用户不是文章作者，用户无修改权限")
		return TravelModel.Post{}, err
	}

	//更新商品信息
	if err := TravelDate.UpdatePost(&post, postR); err != nil {
		return TravelModel.Post{}, err
	}

	return post, nil
}

func PostDelete(userID uint64, postID string) error {
	//查找postID是否存在
	post, err := TravelDate.GetPostInfo(postID)
	if err != nil {
		return err
	}

	//判断是否为商品商家
	if userID != post.UserID {
		err := errors.New("删除失败，用户无修改权限")
		return err
	}

	//删除商品
	if err := TravelDate.DeletePost(post); err != nil {
		return err
	}

	return nil
}

func GetPostInfo(postID string) (TravelModel.Post, error) {
	post, err := TravelDate.GetPostInfo(postID)
	if err != nil {
		return TravelModel.Post{}, err
	}
	return post, nil
}

func PageList(pageNum int, pageList int) ([]TravelModel.Post, int64) {
	commodities, total := TravelDate.PageList(pageNum, pageList)
	return commodities, total
}
