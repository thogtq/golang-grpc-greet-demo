package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/thogtq/golang-grpc-greet-demo/m/v2/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	// fmt.Println("Hello i'm client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to connect server %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	// fmt.Println("Starting a Unary RPC request. . .")
	// doUnary(c)

	// fmt.Println("Starting a Server Streaming RPC request. . .")
	//doServerStreaming(c)

	fmt.Println("Starting a Client Streaming RPC request. . .")
	doClientStreaming(c)
}
func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 40,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greet RPC : %v", err)
	}

	log.Printf("Response from greet RPC server : %v", res)
}
func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 523320,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greet RPC : %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream from RPC server")
		}
		fmt.Println(res.GetPrimeFactor())
	}
}
func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	reqs := []*calculatorpb.ComputeAverageRequest{
		{Number: 32},
		{Number: 66},
		{Number: 17},
		{Number: 43},
		{Number: 55},
	}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverate from RPC server : %v", err)
	}
	for _, req := range reqs {
		fmt.Printf("sending number : %v\n", req)
		stream.Send(req)
		time.Sleep(time.Millisecond * 500)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error when receiving res %v", err)
	}
	fmt.Printf("RPC ComputeAverage res : %v", res.GetAverage())
}
