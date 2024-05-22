package metrics

import (
	pb "collector/proto"
	"context"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
	"time"
)

func getCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func getMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func getHostInfo() string {
	info, _ := host.Info()
	return info.Hostname
}

func getSystemStats() *pb.UploadMetricReq {
	var body []*pb.UploadMetricBody
	timeStamp := time.Now().Unix()
	endpoint := getHostInfo()
	cpuPercent := getCpuPercent()
	memPercent := getMemPercent()
	body = append(body, &pb.UploadMetricBody{
		Metric:    "cpu.used.percent",
		Endpoint:  endpoint,
		Timestamp: timeStamp,
		Step:      60,
		Value:     float32(cpuPercent),
	})
	body = append(body, &pb.UploadMetricBody{
		Metric:    "mem.used.percent",
		Endpoint:  endpoint,
		Timestamp: timeStamp,
		Step:      60,
		Value:     float32(memPercent),
	})
	return &pb.UploadMetricReq{
		Body: body,
	}
}

func CollectAndSendData(client pb.ServerClient) {
	for {

		resp, err := client.UploadMetric(context.Background(), getSystemStats())
		if err != nil {
			log.Printf("Error sending data to server: %s", err)
			continue
		}

		log.Printf("Server response: %s", resp.GetCommon())
		time.Sleep(1 * time.Minute)
	}
}
