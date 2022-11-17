package main

import (
	pg "github.com/namrahov/ms-ecourt-go/greet/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = "localhost:50051"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	defer conn.Close()

	c := pg.NewGreetServiceClient(conn)

	doGreet(c)

}
