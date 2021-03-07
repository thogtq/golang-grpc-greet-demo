package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/thogtq/golang-grpc-greet-demo/m/v2/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	// fmt.Println("Hello World")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("fail to listen port %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("fail to serve %v", err)
	}

}
func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Println("Received RPC request")
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	return &calculatorpb.SumResponse{
		SumResult: sum,
	}, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Println("Received PrimeNumberDecomposition RPC request")
	number := req.GetNumber()
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number /= divisor
		} else {
			divisor++
			log.Printf("Divisor increase to %v", divisor)
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Received ComputeAverate RPC request")
	sum := 0
	count := 0
	average := 0.0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading client stream : %v", err)
		}
		number := req.GetNumber()
		count++
		sum += int(number)
	}
	average = float64(sum) / float64(count)
	stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
		Average: average,
	})
	return nil
}
