package main

import (
	"context"
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
	fmt.Println("Starting a Unary RPC request. . .")
	doUnary(c)
}
func doUnary(c greetpb.GreetServiceClient){
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{FirstName: "Thong", LastName: "Tran Quoc"},
	}

	res, err := c.Greet(context.Background(),req)
	if err!=nil {
		log.Fatalf("error while calling greet RPC : %v",err)
	}

	log.Printf("Response from greet RPC server : %v",res)
}
