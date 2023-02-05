package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/sultanfariz/simple-grpc/client/pkg/grpc/post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type httpServer struct {
	ctx         context.Context
	clientPort  string
	gatewayPort string
}

func NewHTTPServer(ctx context.Context, clientPort, gatewayPort string) *httpServer {
	return &httpServer{
		ctx:         ctx,
		clientPort:  clientPort,
		gatewayPort: gatewayPort,
	}
}

func (h *httpServer) Start() error {
	// start connection to gRPC server
	gatewayConn, err := grpc.DialContext(h.ctx, h.clientPort, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		errors.New("did not connect: %v" + err.Error())
	}

	muxServer := runtime.NewServeMux()
	if err = pb.RegisterPostClientServiceHandler(h.ctx, muxServer, gatewayConn); err != nil {
		errors.New("did not register: %v" + err.Error())
	}
	httpServer := &http.Server{
		Addr:    h.gatewayPort,
		Handler: muxServer,
	}
	fmt.Println("Serving gRPC-Gateway on", h.gatewayPort)
	return httpServer.ListenAndServe()
}
