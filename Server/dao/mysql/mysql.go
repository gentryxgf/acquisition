package mysql

import (
	"context"
	"fmt"
	"server/config"
	"server/model"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(cfg *config.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlBD, err := db.DB()
	if err != nil {
		return
	}

	sqlBD.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlBD.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlBD.SetConnMaxLifetime(time.Hour)
	zap.L().Info("Init MySQL Success!")
	return
}

func CreateMetric(ctx context.Context, data []*model.UsedPercent) error {
	return db.WithContext(ctx).Model(&model.UsedPercent{}).Create(data).Error
}

func QueryByEndpoint(ctx context.Context, endpoint string, start_ts int, end_ts int) ([]*model.UsedPercent, error) {
	var res []*model.UsedPercent
	err := db.WithContext(ctx).Model(&model.UsedPercent{}).Where("endpoint = ? and timestamp >= ? and timestamp <= ?", endpoint, start_ts, end_ts).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func CreateLog(ctx context.Context, data *model.Logs) error {
	return db.WithContext(ctx).Model(&model.Logs{}).Create(data).Error
}

func QueryLogs(ctx context.Context, hostname string, file string) ([]*model.Logs, error) {
	var res []*model.Logs
	err := db.WithContext(ctx).Model(&model.Logs{}).Where("hostname = ? and file = ?", hostname, file).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
