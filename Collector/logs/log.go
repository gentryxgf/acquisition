package logs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	Files      []string `json:"files"`
	LogStorage string   `json:"log_storage"`
}

var (
	conf       Config
	files      = make([]string, 0, 10)
	logStorage string
)

func getConfig() error {

	fmt.Println(os.Getwd())
	viper.SetConfigFile("./config.json")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := viper.Unmarshal(&conf); err != nil {
		zap.L().Info("viper.Umarshal failed", zap.Error(err))
	}
	logStorage = viper.GetString("log_storage")
	files = viper.GetStringSlice("files")
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了！")
		logStorage = viper.GetString("log_storage")
		files = viper.GetStringSlice("files")
		fmt.Println(files)
		fmt.Println(logStorage)
	})
	viper.WatchConfig()
	return nil
}

func Watch() {
	getConfig()
	fmt.Println(files)
	fmt.Println(logStorage)
}
