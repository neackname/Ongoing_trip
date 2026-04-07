package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"travel/TravelDate"
	"travel/config"
	"travel/pkg/snowflake"
	"travel/router"
)

func main() {
	config.InitConfig()
	TravelDate.InitDB()

	// 雪花算法生成分布式ID
	if err := snowflake.Init(1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	r := gin.Default()
	r = router.NewRouter(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	} else {
		panic(r.Run())
	}
}
