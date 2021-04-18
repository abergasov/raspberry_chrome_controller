package main

import (
	"context"
	"log"
	"raspberry_chrome_controller/pkg/config"
	"raspberry_chrome_controller/pkg/logger"
	"raspberry_chrome_controller/pkg/utils"
	"time"

	"google.golang.org/grpc"

	pb "raspberry_chrome_controller/pkg/grpc"

	"go.uber.org/zap"
)

func main() {
	logger.NewLogger()
	appConf := config.InitConf("/etc/commando.yaml")
	logger.Info(
		"Starting app",
		zap.String("token", appConf.KeyToken),
		zap.String("host", appConf.HostURL),
		zap.String("grpc", appConf.GRPCPath),
		zap.Int64("chat", appConf.ListenChat),
	)

	for {
		readStream(appConf)
		time.Sleep(3 * time.Second)
	}
}

func readStream(appConf *config.AppConfig) {
	logger.Info("connect to server", zap.String("path", appConf.GRPCPath))
	conn, err := grpc.Dial(appConf.GRPCPath, grpc.WithInsecure())
	if err != nil {
		logger.Error("can not connect with server", err)
		return
	}

	// create stream
	client := pb.NewCommandStreamClient(conn)
	stream, err := client.ListenCommands(context.Background(), &pb.Request{TargetChat: appConf.ListenChat})
	if err != nil {
		logger.Error("open stream error", err)
		return
	}

	var resp *pb.Response
	commando := utils.NewCommandor()
	for {
		resp, err = stream.Recv()
		if err != nil {
			logger.Error("can not receive", err)
			return
		}
		log.Printf("Resp received: %s - %s", resp.ActionID, resp.Cmd)
		commando.HandleCommand(&utils.Command{
			Cmd:      resp.Cmd,
			ActionID: resp.ActionID,
		})
	}
}
