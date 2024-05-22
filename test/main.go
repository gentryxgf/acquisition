package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	Files      []string `json:"files"`
	LogStorage string   `json:"log_storage"`
}

var conf = new(Config)

func Init() error {
	viper.SetConfigFile("./test/config.json")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := viper.Unmarshal(conf); err != nil {
		zap.L().Info("viper.Umarshal failed", zap.Error(err))
	}
	fmt.Println(viper.Get("files"))
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了！")
		fmt.Println(in.Name)
		fmt.Println(viper.Get("log_storage"))
	})
	viper.WatchConfig()
	return nil
}

func main() {

	//f, _ := os.Getwd()
	//path := f + "/config.json"
	//fmt.Println(path)
	err := Init()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.Get("log_storage"))
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

}
