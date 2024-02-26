package main

import (
	"context"
	"fmt"
	"github.com/NubeIO/nogrpc"
	"github.com/NubeIO/nogrpc/demo/clients/go/pb"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

var (
	sharePort    int
	shutdownFunc func()
)

func init() {
	sharePort = 8081
	shutdownFunc = func() {
		fmt.Println("Server shutting down")
	}
}

func main() {
	// add the /test endpoint
	route := nogrpc.Route{
		Method: "GET",
		Path:   "/test",
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Write([]byte("Hello!"))
		},
	}

	// test Option func
	s := nogrpc.NewService(
		nogrpc.WithRouteOpt(route),
		nogrpc.WithShutdownFunc(shutdownFunc),
		nogrpc.WithPreShutdownDelay(1*time.Second),
		nogrpc.WithShutdownTimeout(1*time.Second),
		nogrpc.WithHandlerFromEndpoint(pb.RegisterGreeterServiceHandlerFromEndpoint),
		nogrpc.WithLogger(nogrpc.LoggerFunc(log.Printf)),
		nogrpc.WithRequestAccess(true),
		nogrpc.WithGRPCServerOption(grpc.ConnectionTimeout(10*time.Second)),
		nogrpc.WithGRPCNetwork("tcp"), // grpc server start network
		nogrpc.WithStaticAccess(true), // enable static file access,if use http gw
	)
	// register grpc service
	pb.RegisterGreeterServiceServer(s.GRPCServer, &greeterService{})

	newRoute := nogrpc.Route{
		Method: "GET",
		Path:   "health",
		Handler: func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	}

	s.AddRoute(newRoute)

	// you can start grpc server and http gateway on one port
	//log.Fatalln(s.StartGRPCAndHTTPServer(sharePort))

	// you can also specify ports for grpc and http gw separately
	log.Fatalln(s.Start(sharePort, 50051))

	// you can start server without http gateway
	// log.Fatalln(s.StartGRPCWithoutGateway(50051))
}

// rpc service entry
type greeterService struct {
	pb.UnimplementedGreeterServiceServer
}

func (s *greeterService) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloReply, error) {
	log.Println("req data: ", in)
	time.Sleep(12 * time.Millisecond)
	md := nogrpc.GetIncomingMD(ctx)
	log.Println("request md: ", md)
	return &pb.HelloReply{
		Name:    "hello," + in.Name,
		Message: "call ok",
	}, nil
}
