package main

import (
	"context"
	"log"
	"raspberry_chrome_controller/pkg/config"
	"raspberry_chrome_controller/pkg/logger"
	"raspberry_chrome_controller/pkg/utils"

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

	// dail server
	conn, err := grpc.Dial(appConf.GRPCPath, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("can not connect with server", err)
	}

	// create stream
	client := pb.NewCommandStreamClient(conn)
	stream, err := client.ListenCommands(context.Background(), &pb.Request{TargetChat: appConf.ListenChat})
	if err != nil {
		logger.Fatal("open stream error", err)
	}

	var resp *pb.Response
	commando := utils.NewCommandor()
	for {
		resp, err = stream.Recv()
		if err != nil {
			logger.Fatal("can not receive", err)
		}
		log.Printf("Resp received: %s - %s", resp.ActionID, resp.Cmd)
		commando.HandleCommand(&utils.Command{
			Cmd:      resp.Cmd,
			ActionID: resp.ActionID,
		})
	}
}
