package controller

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"travel/MySQLTavelDate"
	"travel/TravelModel"
)

type Post struct {
	ID      uuid.UUID `gorm:"type:varchar(36); primary_key"`
	UserId  uint      `gorm:"not null"`
	Title   string    `gorm:"type:varchar(15); not null"`
	HeadImg string    `gorm:"type:text"`
	Content string    `gorm:"type:text; not null"`
}

type PostRequest struct {
	Title   string `json:"title"`
	HeadImg string `json:"head_img"`
	Content string `json:"content"`
}

func PostCreate(ctx *gin.Context) {
	var p PostRequest
	if err := ctx.ShouldBind(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	//创建文章
	user, _ := ctx.Get("user")
	post := Post{
		UserId:  user.(TravelModel.TraUser).ID,
		Title:   p.Title,
		HeadImg: p.HeadImg,
		Content: p.Content,
	}

	db := MySQLTavelDate.GetDB()
	if err := db.Create(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "创建失败"})
		return
	}
}
