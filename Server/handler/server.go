package handler

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"server/dao/mysql"
	"server/dao/redis"
	"server/dao/storage"
	"server/model"
	"server/proto"
)

type ServerSrv struct {
	proto.UnimplementedServerServer
}

func (s *ServerSrv) UploadMetric(ctx context.Context, req *proto.UploadMetricReq) (*proto.UploadMetricResp, error) {
	var (
		datas = make([]*model.UsedPercent, 0, 2)
	)

	reqBody := req.GetBody()
	for _, body := range reqBody {
		datas = append(datas, &model.UsedPercent{
			Metric:    body.Metric,
			Endpoint:  body.Endpoint,
			Timestamp: int(body.Timestamp),
			Step:      int(body.Step),
			Value:     body.Value,
		})
	}
	err := mysql.CreateMetric(ctx, datas)
	if err != nil {
		zap.L().Error("UploadMetric.mysql.CreateMetric error:", zap.Error(err))

		return &proto.UploadMetricResp{
			Common: &proto.CommonResp{
				Code:    500,
				Message: "失败！",
			},
			Data: "",
		}, err
	}
	rdskey := datas[0].Endpoint
	b, err := json.Marshal(datas)
	if err != nil {
		zap.L().Error("UploadMetric.mysql.CreateMetric error:", zap.Error(err))
	}
	redis.Rc.LPush(ctx, rdskey, b)
	redis.Rc.LTrim(ctx, rdskey, 0, 9)

	return &proto.UploadMetricResp{
		Common: &proto.CommonResp{
			Code:    200,
			Message: "Success",
		},
		Data: "Success",
	}, nil
}

func (s *ServerSrv) QueryMetric(ctx context.Context, req *proto.QueryMetricReq) (*proto.QueryMetricResp, error) {
	var (
		respData  []*proto.QueryMetricRespData
		cpuMetric []*proto.MetricValue
		memMetric []*proto.MetricValue
	)
	datas, err := mysql.QueryByEndpoint(ctx, req.GetEndpoint(), int(req.GetStartTs()), int(req.GetEndTs()))
	if err != nil {
		zap.L().Error("QueryMetric.mysql.QueryByEndpoint:", zap.Error(err))
		return &proto.QueryMetricResp{
			Common: &proto.CommonResp{
				Code:    500,
				Message: err.Error(),
			},
			Data: nil,
		}, nil
	}
	for _, d := range datas {
		if d.Metric == "cpu.used.percent" {
			cpuMetric = append(cpuMetric, &proto.MetricValue{Timestamp: int64(d.Timestamp), Value: d.Value})
		} else {
			memMetric = append(memMetric, &proto.MetricValue{Timestamp: int64(d.Timestamp), Value: d.Value})
		}
	}
	if req.Metric == "cpu.used.percent" {
		respData = append(respData, &proto.QueryMetricRespData{Metric: req.Metric, MetricValues: cpuMetric})
	} else if req.Metric == "mem.used.percent" {
		respData = append(respData, &proto.QueryMetricRespData{Metric: req.Metric, MetricValues: memMetric})
	} else {
		respData = append(respData, &proto.QueryMetricRespData{Metric: "cpu.used.percent", MetricValues: cpuMetric})
		respData = append(respData, &proto.QueryMetricRespData{Metric: "mem.used.percent", MetricValues: memMetric})
	}
	resp := &proto.QueryMetricResp{Data: respData, Common: &proto.CommonResp{Code: 200, Message: "查询成功！"}}
	return resp, nil
}

func (s *ServerSrv) UploadLog(ctx context.Context, req *proto.UploadLogReq) (*proto.UploadLogResp, error) {
	reqBody := req.GetBody()
	for _, body := range reqBody {
		instance := storage.GetStorageInstance("file", body.Hostname, body.File)
		for _, l := range body.Logs {
			log := &model.Logs{Hostname: body.Hostname, File: body.File, Log: l}
			//err := mysql.CreateLog(ctx, log)
			err := instance.UploadLog(ctx, log)
			if err != nil {
				zap.L().Error("UploadLog.mysql.CreateLog error:", zap.Error(err))
				return &proto.UploadLogResp{
					Common: &proto.CommonResp{
						Code:    500,
						Message: err.Error(),
					},
					Data: "Error!",
				}, nil
			}
		}
	}
	return &proto.UploadLogResp{
		Common: &proto.CommonResp{
			Code:    200,
			Message: "Success",
		},
		Data: "",
	}, nil
}

func (s *ServerSrv) QueryLog(ctx context.Context, req *proto.QueryLogReq) (*proto.QueryLogResp, error) {
	var (
		logs []string
	)
	instance := storage.GetStorageInstance("mysql", req.GetHostname(), req.GetFile())
	//datas, err := mysql.QueryLogs(ctx, req.GetHostname(), req.GetFile())
	datas, err := instance.QueryLog(ctx, req.GetHostname(), req.GetFile())
	if err != nil {
		zap.L().Error("QueryMetric.mysql.QueryByEndpoint:", zap.Error(err))
		return &proto.QueryLogResp{
			Common: &proto.CommonResp{
				Code:    500,
				Message: err.Error(),
			},
			Data: nil,
		}, nil
	}
	for _, l := range datas {
		logs = append(logs, l.Log)
	}
	return &proto.QueryLogResp{
		Common: &proto.CommonResp{
			Code:    200,
			Message: "Ok",
		},
		Data: &proto.QueryLogRespData{
			Hostname: req.GetHostname(),
			File:     req.GetFile(),
			Logs:     logs,
		},
	}, nil
}
