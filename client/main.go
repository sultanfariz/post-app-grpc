package main

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"github.com/sultanfariz/simple-grpc/client/cmd"
)

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if viper.GetBool("debug") {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	ctx := context.Background()
	grpcServer := cmd.NewgRPCServer(ctx, ":"+viper.GetString("SERVER_PORT"), ":"+viper.GetString("CLIENT_PORT"))
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	httpServer := cmd.NewHTTPServer(ctx, ":"+viper.GetString("CLIENT_PORT"), ":"+viper.GetString("GATEWAY_PORT"))
	if err := httpServer.Start(); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
