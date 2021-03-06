package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/thogtq/golang-grpc-greet-demo/m/v2/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type server struct{}

func main() {
	// fmt.Println("Hello World")
	const certFile = "../../ssl/server.crt"
	const keyFile = "../../ssl/server.pem"
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("fail to load cerificate key %v", err)
	}
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("fail to listen port %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer(grpc.Creds(creds))
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("fail to serve %v", err)
	}

}
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet func invoked with %v", req)
	firstName := req.Greeting.GetFirstName()
	result := "Hello " + firstName
	_ = result
	return &greetpb.GreetResponse{
		Result: result,
	}, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes func invoked with %v", req)
	firstName := req.Greeting.GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet func invoked")
	result := "Hello "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//Reached end of stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading client stream")
		}
		firstName := req.Greeting.GetFirstName()
		result += firstName + "! "
	}
	stream.SendAndClose(&greetpb.LongGreetResponse{
		Result: result,
	})
	return nil
}
func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone func invoked")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receive client stream : %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "hello " + firstName + "!"
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("error while sending to client : %v", sendErr)
		}
	}
	return nil
}
func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline func invoked\n")
	time.Sleep(3 * time.Second)
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Printf("The client has cancelled the request\n")
		return nil, status.Error(codes.Canceled, "The client has cancelled the request")
	}
	firstName := req.Greeting.GetFirstName()
	result := "Hello " + firstName
	return &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}, nil
}
