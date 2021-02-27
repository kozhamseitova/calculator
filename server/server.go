package server

import (
	calculator "com.grpc.aisha/calculatorpb"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct{}

func (*server) PrimeNumberDecomposition(req *calculator.PrimeNumberDecompositionRequest, stream calculator.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Received PrimeNumberDecomposition RPC: %v\n", req)
	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculator.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			log.Printf("Dvisor has increased to %v", divisor)
		}
	}

	return nil
}

func (*server) ComputeAverage(stream calculator.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Received ComputeAverage RPC")

	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			stream.SendAndClose(&calculator.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		sum += req.GetNumber()
		count++
	}
}

func main() {
	fmt.Println("Running gRPC on port 50051...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculator.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
