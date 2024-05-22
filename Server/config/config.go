package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Conf = new(Config)

type Config struct {
	Name     string `mapstructure:"name"`
	Mode     string `mapstructure:"mode"`
	Ip       string `mapstructure:"ip"`
	RpcPort  int    `mapstructure:"rpcPort"`
	HttpPort int    `mapstructure:"httpPort"`
	Version  string `mapstructure:"version"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init(filePath string) (err error) {
	viper.SetConfigFile(filePath)

	err = viper.ReadInConfig()
	if err != nil {
		zap.L().Info("viper.ReadInConfig failed", zap.Error(err))
		return
	}

	if err := viper.Unmarshal(Conf); err != nil {
		zap.L().Info("viper.Umarshal failed", zap.Error(err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Info("配置文件修改了！")
		if err := viper.Unmarshal(Conf); err != nil {
			zap.L().Info("viper.Umarshal failed", zap.Error(err))
		}
	})
	return
}
