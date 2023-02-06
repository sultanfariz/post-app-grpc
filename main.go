package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/viper"
	postDomain "github.com/sultanfariz/simple-grpc/domain/posts"
	userDomain "github.com/sultanfariz/simple-grpc/domain/users"
	"github.com/sultanfariz/simple-grpc/infrastructure/commons"
	"github.com/sultanfariz/simple-grpc/infrastructure/repository/mysql"
	postsRepository "github.com/sultanfariz/simple-grpc/infrastructure/repository/mysql/posts"
	userRepository "github.com/sultanfariz/simple-grpc/infrastructure/repository/mysql/users"
	grpcServerController "github.com/sultanfariz/simple-grpc/infrastructure/transport/grpc"
	rabbitmq "github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	lis, err := net.Listen("tcp", ":"+viper.GetString("SERVER_PORT"))
	timeoutContext := time.Duration(viper.GetInt("CONTEXT_TIMEOUT")) * time.Second

	if err != nil {
		log.Fatalf("failed listen: %v", err)
	}
	fmt.Println("Server running on port " + viper.GetString("SERVER_PORT"))

	configJWT := commons.ConfigJWT{
		SecretJWT:       viper.GetString("JWT_SECRET_KEY"),
		ExpiresDuration: viper.GetInt("JWT_EXPIRES_DURATION"),
	}

	// rabbitmq connection
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	publisher.NotifyReturn(func(r rabbitmq.Return) {
		log.Printf("message returned from server: %s", string(r.Body))
	})

	publisher.NotifyPublish(func(c rabbitmq.Confirmation) {
		log.Printf("message confirmed from server. tag: %v, ack: %v", c.DeliveryTag, c.Ack)
	})

	db := mysql.InitDB()
	userRepo := userRepository.NewUsersRepository(db)
	userUsecase := userDomain.NewUsersUsecase(userRepo, timeoutContext, &configJWT)
	postRepo := postsRepository.NewPostsRepository(db)
	postUsecase := postDomain.NewPostsUsecase(postRepo, userRepo, timeoutContext, publisher)

	// grpc server
	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpcServerController.JWTInterceptor),
	}
	grpcServer := grpc.NewServer(serverOpts...)

	grpcServerController.NewUserServerGrpc(grpcServer, *userUsecase)
	grpcServerController.NewPostServerGrpc(grpcServer, *postUsecase)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
