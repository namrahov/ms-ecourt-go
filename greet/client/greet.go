package main

import (
	"context"
	pg "github.com/namrahov/ms-ecourt-go/greet/proto"
	log "github.com/sirupsen/logrus"
)

func doGreet(c pg.GreetServiceClient) {

	log.Println("doGreet was invoke")
	res, err := c.Greet(context.Background(), &pg.GreetRequest{
		FirstName: "Clement",
	})

	if err != nil {
		log.Fatalf("couldn't great %v\n", err)
	}

	log.Printf("Greeting: %s\n", res.Result)

}
