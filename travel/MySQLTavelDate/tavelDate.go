package MySQLTavelDate

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/url"
	"travel/TravelModel"
)

var DB *gorm.DB

// @title InitDB
// @description	初始化数据库
// @auth	Snactop	2023-11-27	20:07
// @param	void	没有传入值
// @return	void	没有返回值
func InitDB() {
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	database := viper.GetString("mysql.database")
	root := viper.GetString("mysql.root")
	password := viper.GetString("mysql.password")
	loc := viper.GetString("mysql.loc")

	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		root,
		password,
		host,
		port,
		database,
		url.QueryEscape(loc))

	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("fail to connect database, err:" + err.Error() + host)
	}
	if err := db.AutoMigrate(&TravelModel.TraUser{}, &TravelModel.TraUserFoot{}, &TravelModel.TraUserFootStart{}, TravelModel.TraUserPostStart{}); err != nil {
		panic(err)
		return
	}

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
