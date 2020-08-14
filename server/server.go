package server

import (
	"fmt"
	"net"

	"github.com/grpc-demo/message"
	"github.com/grpc-demo/provider"
	"google.golang.org/grpc"
)

func Start() {
	server := grpc.NewServer()
	message.RegisterOrderServiceServer(server, new(provider.OrderServiceImpl))
	listen, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("server starting on port:8090")
	server.Serve(listen)
}
