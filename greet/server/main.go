package main

import (
	"fmt"
	pg "github.com/namrahov/ms-ecourt-go/greet/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var addr = "0.0.0.0:50051"

type Server struct {
	pg.GreetServiceServer
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("lis=", err)
	}

	log.Printf("Listeninng on %s", addr)
	s := grpc.NewServer()
	pg.RegisterGreetServiceServer(s, &Server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Fatal %v\n", err)
	}
}
