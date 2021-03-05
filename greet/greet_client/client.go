package main

import (
	"fmt"
	"log"

	"github.com/thogtq/golang-grpc-greet-demo/m/v2/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello i'm client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to connect server %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("Created client : %f", c)

}
