package main

import (
	"context"
	"fmt"
	"log"

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
	fmt.Println("Starting a Unary RPC request. . .")
	doUnary(c)
}
func doUnary(c calculatorpb.CalculatorServiceClient){
	req := &calculatorpb.SumRequest{
		FirstNumber: 10,
		SecondNumber: 40,
	}

	res, err := c.Sum(context.Background(),req)
	if err!=nil {
		log.Fatalf("error while calling greet RPC : %v",err)
	}

	log.Printf("Response from greet RPC server : %v",res)
}
