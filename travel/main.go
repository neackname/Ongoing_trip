package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"travel/MySQLTavelDate"
	"travel/config"
	"travel/router"
)

func main() {
	config.InitConfig()
	MySQLTavelDate.InitDB()
	r := gin.Default()
	r = router.NewRouter(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	} else {
		panic(r.Run())
	}
}
