package main

import (
	"context"
	"fmt"
	"io"
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

	// fmt.Println("Starting a Unary RPC request. . .")
	// doUnary(c)

	fmt.Println("Starting a Server Streaming RPC request. . .")
	doServerStreaming(c)
}
func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{FirstName: "Thong", LastName: "Tran Quoc"},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling server streaming rpc : %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//reached the end of stream
			break
		} else if err != nil {
			log.Fatalf("error when reading stream from rpc server : %v", err)
		}
		log.Printf("Response from GreetManyTimes rpc server : \n%s", msg.GetResult())
	}

}
func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{FirstName: "Thong", LastName: "Tran Quoc"},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greet RPC : %v", err)
	}

	log.Printf("Response from greet RPC server : %v", res)
}
