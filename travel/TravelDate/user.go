package TravelDate

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"travel/TravelModel"
	"travel/vo"
)

// CheckUserExistWithOpenID
// @Description 数据库查找用户是否存在
// @auth Snactop 2024/12/10 0:25
// @param  userName string  传入一个用户名
// @return bool 返回一个布尔值，用户存在返回true，反之则false
func CheckUserExistWithOpenID(OpenID string) bool {
	db := GetDB()
	//TODO 判断用户是否存在
	var user TravelModel.TraUser
	db.Where("  open_id  = ?", OpenID).First(&user)

	if user.ID != 0 {
		return true
	}
	return false
}

func GetUserInformation(OpenID string) (TravelModel.TraUser, error) {
	db := GetDB()
	//TODO 判断用户是否存在
	var user TravelModel.TraUser
	db.Where("  open_id  = ?", OpenID).First(&user)
	if user.ID == 0 {
		return TravelModel.TraUser{}, errors.New("用户不存在")
	}

	return user, nil
}

func UpdateUserInformation(user TravelModel.TraUser, userUpdate vo.UpdateUserRequest) error {
	db := GetDB()

	if err := db.Model(&user).Updates(userUpdate).Error; err != nil {
		return errors.New("系统错误，用户信息更新失败")
	}
	return nil
}

func CheckUserPostStartExist(userID uint64, postID uuid.UUID) bool {
	db := GetDB()
	var existingFavorite TravelModel.TraUserPostStart
	if err := db.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingFavorite).Error; err == nil {
		return true
	}
	return false
}

func RemovePostStart1(userID uint64, postID uuid.UUID) error {
	db := GetDB()
	if err := db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&TravelModel.TraUserPostStart{}).Error; err != nil {
		return errors.New("取消收藏失败")
	}
	return nil
}

func AddPostStart(userID uint64, postID uuid.UUID) error {
	db := GetDB()
	favorite := TravelModel.TraUserPostStart{
		UserID: userID, // 转换为uint64类型
		PostID: postID,
	}
	if err := db.Create(&favorite).Error; err != nil {
		err = errors.New("收藏失败，请稍后重试")
		return err
	}
	return nil
}

func GetPostStartIDs(userID uint64) ([]uuid.UUID, error) {
	db := GetDB()
	var postIDs []uuid.UUID
	err := db.Table("tra_user_post_starts"). //收藏表的表名
							Select("post_id").
							Where("user_id = ?", userID).
							Pluck("post_id", &postIDs).Error

	if err != nil {
		return nil, err
	}

	return postIDs, nil
}

func GetPostsByIDs(ids []uuid.UUID) ([]TravelModel.Post, error) {
	db := GetDB()

	var posts []TravelModel.Post

	err := db.Table("posts").
		Where("id IN ?", ids).
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func GetUserCreatedPosts(userID uint64) (posts []TravelModel.Post, err error) {
	db := GetDB()

	// 查询用户创建的所有商品
	err = db.Table("posts").
		Where("user_id = ?", userID).
		Find(&posts).Error

	if err != nil {
		return []TravelModel.Post{}, errors.New("获取用户创建的文章失败，请稍后重试")
	}

	// 如果用户没有创建任何商品
	if len(posts) == 0 {
		return []TravelModel.Post{}, nil
	}

	return posts, nil
}
