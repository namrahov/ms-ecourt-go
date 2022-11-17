package main

import (
	"context"
	pg "github.com/namrahov/ms-ecourt-go/greet/proto"
	"log"
)

func (s *Server) Greet(ctx context.Context, in *pg.GreetRequest) (*pg.GreetResponse, error) {
	log.Printf("Greet function was invoked with %v\n", in)
	return &pg.GreetResponse{
		Result: in.FirstName,
	}, nil
}
