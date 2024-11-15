package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

// @title     InitConfig
// @description     初始化配置文件
// @auth      Snactop            2023-12-5 18:01
// @return    void	没有回参
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("ReadInConfig Failed, err:%v", err))
	}
}
