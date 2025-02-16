package TravelDate

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"travel/TravelModel"
	"travel/vo"
)

func InsertPost(post TravelModel.Post) (uuid.UUID, error) {
	db := GetDB()
	if err := db.Create(&post).Error; err != nil {
		return uuid.Nil, err
	}
	return post.ID, nil
}

func UpdatePost(post *TravelModel.Post, postR vo.PostRequest) error {
	db := GetDB()
	if err := db.Model(post).Updates(postR).Error; err != nil {
		err := errors.New("商品信息更新失败")
		return err
	}
	return nil
}

func DeletePost(post TravelModel.Post) error {
	db := GetDB()
	if db.Delete(&post).Error != nil {
		return errors.New("系统错误，商品信息删除失败")
	}
	return nil
}

func GetPostInfo(postID string) (TravelModel.Post, error) {
	db := GetDB()
	// 判断商品是否存在
	var post TravelModel.Post
	db.Where("id = ?", postID).First(&post)
	if post.ID == uuid.Nil {
		err := errors.New("商品不存在")
		return TravelModel.Post{}, err
	}
	return post, nil
}

func PageList(pageNum int, pageList int) ([]TravelModel.Post, int64) {
	db := GetDB()
	var commodities []TravelModel.Post
	db.Order("created_at desc").Offset((pageNum - 1) * pageList).Limit(pageList).Find(&commodities)
	//前端渲染分页需要知道总数
	var total int64
	db.Model([]TravelModel.Post{}).Count(&total)
	return commodities, total
}
