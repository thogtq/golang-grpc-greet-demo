package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	// fmt.Println("Starting a Server Streaming RPC request. . .")
	// doServerStreaming(c)

	// fmt.Println("Starting a Client Streaming RPC request. . .")
	// doClientStreaming(c)

	fmt.Println("Starting a BiDi Streaming RPC request. . .")
	doBiDiStreaming(c)
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
func doClientStreaming(c greetpb.GreetServiceClient) {
	reqs := []*greetpb.LongGreetRequest{
		{Greeting: &greetpb.Greeting{FirstName: "Thong", LastName: "Tran Quoc"}},
		{Greeting: &greetpb.Greeting{FirstName: "Anh", LastName: "Nguyen"}},
		{Greeting: &greetpb.Greeting{FirstName: "Bao", LastName: "Nguyen"}},
		{Greeting: &greetpb.Greeting{FirstName: "Duy", LastName: "Nguyen"}},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet from RPC server")
	}
	for _, req := range reqs {
		fmt.Printf("Sending request %v\n", req)
		stream.Send(req)
		time.Sleep(time.Millisecond * 500)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receive server res")
	}
	fmt.Printf("Server LongGreeet response : %v", res.GetResult())
}
func doBiDiStreaming(c greetpb.GreetServiceClient) {
	//Create a stream by invokeing the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreetEveryone from RPC server")
		return
	}
	waitChannel := make(chan struct{})
	reqs := []*greetpb.GreetEveryoneRequest{
		{Greeting: &greetpb.Greeting{FirstName: "Thong", LastName: "Tran Quoc"}},
		{Greeting: &greetpb.Greeting{FirstName: "Anh", LastName: "Nguyen"}},
		{Greeting: &greetpb.Greeting{FirstName: "Bao", LastName: "Nguyen"}},
		{Greeting: &greetpb.Greeting{FirstName: "Duy", LastName: "Nguyen"}},
	}
	//Send a bunch of messages to client (gorountine)
	go func() {
		for _, req := range reqs {
			fmt.Printf("Sending request %v\n", req)
			stream.Send(req)
			time.Sleep(time.Millisecond * 500)
		}
		stream.CloseSend()
	}()
	//Receive a bunch of messages from client (gorountine)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving server stream: %v", err)
				break
			}
			log.Printf("Response from GreetEveryone rpc server : \n%s\n", res.GetResult())
		}
		close(waitChannel)
	}()
	//Block until everything is done
	<-waitChannel
}
