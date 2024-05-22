package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"os"
	"os/signal"
	"server/config"
	"server/dao/mysql"
	"server/dao/redis"
	"server/handler"
	"server/logger"
	"server/proto"
	"syscall"
)

func main() {
	var cfn string
	//命令行解析配置文件
	// good_service -conf="./conf/config.yaml"
	flag.StringVar(&cfn, "conf", "./config.yaml", "指定配置文件路径")
	flag.Parse()

	// 加载配置文件
	err := config.Init(cfn)
	if err != nil {
		// 配置文件加载失败直接退出
		panic(err)
	}

	// 加载日志
	err = logger.Init(config.Conf.LogConfig, config.Conf.Mode)
	if err != nil {
		// 日志加载失败直接退出
		panic(err)
	}

	// MySQL初始化
	err = mysql.Init(config.Conf.MySQLConfig)
	if err != nil {
		// MySQL初始化失败直接退出
		panic(err)
	}

	// redis
	err = redis.Init(config.Conf.RedisConfig)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Conf.RpcPort))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	proto.RegisterServerServer(s, &handler.ServerSrv{})
	go func() {
		err := s.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()
	zap.L().Sugar().Infof("Service start at: %s-%s-%d", config.Conf.Name, config.Conf.Ip, config.Conf.RpcPort)
	// 创建grpc客户端
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("%s:%d", config.Conf.Ip, config.Conf.RpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		zap.L().Fatal("Failed to dial server:", zap.Error(err))
	}
	gwmux := runtime.NewServeMux()
	err = proto.RegisterServerHandler(context.Background(), gwmux, conn)
	if err != nil {
		zap.L().Fatal("Failed to register gatewary:", zap.Error(err))
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.HttpPort),
		Handler: gwmux,
	}
	zap.L().Sugar().Infof("Serving gRPC-GateWay on http: 0.0.0.0%s", gwServer.Addr)

	go func() {
		err := gwServer.ListenAndServe()
		if err != nil {
			fmt.Printf("gwServer.ListenAndServe failed, err:%s", err)
		}
		return
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
