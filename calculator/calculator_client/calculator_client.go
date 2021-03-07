package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/thogtq/golang-grpc-greet-demo/m/v2/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// fmt.Println("Starting a Client Streaming RPC request. . .")
	// doClientStreaming(c)

	// fmt.Println("Starting a BiDi Streaming RPC request. . .")
	// doBidiStreaming(c)

	doErrorUnary(c)
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
func doBidiStreaming(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while calling FindMaximum from RPC server : %v", err)
		return
	}
	waitChannel := make(chan struct{})
	reqs := []*calculatorpb.FindMaximumRequest{
		{Number: 2},
		{Number: 8},
		{Number: 4},
		{Number: 11},
		{Number: 1},
	}
	go func() {
		for _, req := range reqs {
			fmt.Printf("sending number : %v\n", req)
			stream.Send(req)
			time.Sleep(time.Second * 1)
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading server stream : %v", err)
				return
			}
			fmt.Printf("Response from FindMaximum RPC server : %v\n", res.GetMaximum())
		}
		close(waitChannel)
	}()
	<-waitChannel
}
func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	//This function implementation with gRPC error handling
	req := &calculatorpb.SquareRootRequest{
		Number: -16,
	}
	res, err := c.SquareRoot(context.Background(), req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			if respErr.Code() == codes.InvalidArgument {
				fmt.Printf("Don't try to send negative number to SquareRoot : %v\n", req.GetNumber())
				return
			}
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			return
		}
		log.Fatalf("error while calling SquareRoot RPC : %v", err)
		return
	}
	log.Printf("Response from SquareRoot RPC server : %v", res.GetNumberRoot())
}
