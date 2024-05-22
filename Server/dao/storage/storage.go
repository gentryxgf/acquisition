package storage

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"server/dao/mysql"
	"server/model"
	"strings"
)

// 抽象产品类
type LogsStorage interface {
	UploadLog(ctx context.Context, data *model.Logs) error
	QueryLog(ctx context.Context, hostname string, file string) ([]*model.Logs, error)
}

// 具体产品类
type MysqlStorage struct {
}

func (m *MysqlStorage) UploadLog(ctx context.Context, data *model.Logs) error {
	return mysql.CreateLog(ctx, data)
}

func (m *MysqlStorage) QueryLog(ctx context.Context, hostname string, file string) ([]*model.Logs, error) {
	return mysql.QueryLogs(ctx, hostname, file)
}

type FileStorage struct {
	Hostname string
	File     string
}

func (f *FileStorage) writeFile(data *model.Logs) error {
	index := strings.LastIndex(f.File, "/")
	dir := f.Hostname + "/" + f.File[:index]
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		zap.L().Error("getFile.MkdirAll error:", zap.Error(err))
		return err
	}
	file, err := os.OpenFile(f.Hostname+f.File, os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		zap.L().Error("UploadLog.f.getfile error:", zap.Error(err))
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(data.Log))
	file.WriteString("\n")
	if err != nil {
		zap.L().Error("UploadLog.file.Write error:", zap.Error(err))
		return err
	}
	return nil
}

func (f *FileStorage) readFile(hostname string, file string) ([]string, error) {
	ff, err := os.Open(hostname + file)
	fmt.Println(hostname + file)
	if err != nil {
		zap.L().Error("FileStorage.readFile.os.Open error:", zap.Error(err))
		return nil, err
	}
	stats, statsErr := ff.Stat()
	if statsErr != nil {
		zap.L().Error("FileStorage.readFile.ff.Stat error:", zap.Error(err))
		return nil, statsErr
	}
	size := stats.Size()
	lines := make([]string, 0, 10)
	buf := make([]byte, size)
	for {
		readSize, err := ff.ReadAt(buf, size-int64(len(buf)))
		if err != nil && err != io.EOF {
			break
		}

		size -= int64(readSize)

		for i := readSize - 1; i >= 0; i-- {
			if buf[i] == '\n' {
				lines = append(lines, string(buf[i+1:readSize]))
				readSize = i
				if len(lines) == 10 {
					return lines, nil
				}
			}
		}

		if size == 0 {
			lines = append(lines, string(buf[:readSize]))
			break
		}

		if len(lines) >= 10 {
			break
		}

		if size < int64(len(buf)) {
			buf = make([]byte, size)
		}
	}
	return lines, nil

}

func (f *FileStorage) UploadLog(ctx context.Context, data *model.Logs) error {
	err := f.writeFile(data)
	if err != nil {
		zap.L().Error("UploadLog.file.Write error:", zap.Error(err))
		return err
	}
	return nil
}

func (f *FileStorage) QueryLog(ctx context.Context, hostname string, file string) ([]*model.Logs, error) {
	var res = make([]*model.Logs, 0, 10)
	fmt.Println(hostname, file)
	data, err := f.readFile(hostname, file)
	if err != nil {
		zap.L().Error("QueryLog.file.readFile error:", zap.Error(err))
		return nil, err
	}
	for _, d := range data {
		res = append(res, &model.Logs{
			Hostname: hostname,
			File:     file,
			Log:      d,
		})
	}
	return res, nil
}

// 抽象工厂类
type StorageFactory interface {
	CreateFactory() LogsStorage
}

// 具体工厂类
type MysqlFactory struct {
}

func (m *MysqlFactory) CreateFactory() LogsStorage {
	return &MysqlStorage{}
}

type FileFactory struct {
}

func (m *FileFactory) CreateFactory() LogsStorage {
	return &FileStorage{}
}

func GetStorageInstance(method string, hostname string, file string) LogsStorage {
	if method == "mysql" {
		return &MysqlStorage{}
	}
	if method == "file" {
		return &FileStorage{Hostname: hostname, File: file}
	}
	return nil
}
