package main

import (
	"fmt"
	"log"
	"net"
	"time"

	userDomain "github.com/sultanfariz/simple-grpc/domain/users"
	"github.com/sultanfariz/simple-grpc/infrastructure/repository/mysql"
	userRepository "github.com/sultanfariz/simple-grpc/infrastructure/repository/mysql/users"
	grpcServerController "github.com/sultanfariz/simple-grpc/infrastructure/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50001")
	timeoutContext := time.Duration(10 * time.Second)

	if err != nil {
		log.Fatalf("failed listen: %v", err)
	}
	fmt.Println("Server running on port :50001")

	db := mysql.InitDB()
	userRepo := userRepository.NewUsersRepository(db)
	userUsecase := userDomain.NewUsersUsecase(userRepo, timeoutContext)

	
	grpcServer := grpc.NewServer()
	grpcServerController.NewUserServerGrpc(grpcServer, *userUsecase)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}