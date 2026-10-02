package tests

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	pb "travelraft/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestConcurrentBookingRequests(t *testing.T) {
	if os.Getenv("TRAVELRAFT_LIVE_TESTS") != "1" {
		t.Skip("set TRAVELRAFT_LIVE_TESTS=1 to run against a live server")
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	client := pb.NewBookingServiceClient(conn)
	const requests = 20
	var wait sync.WaitGroup
	responses := make(chan *pb.BookResponse, requests)
	for index := 0; index < requests; index++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			response, err := client.Book(ctx, &pb.BookRequest{
				PassengerName: fmt.Sprintf("Load User %d", index),
				Email:         fmt.Sprintf("load-%d@example.com", index),
				VehicleId:     1,
				SeatId:        int32(2 + index),
			})
			if err != nil {
				responses <- &pb.BookResponse{Message: err.Error()}
				return
			}
			responses <- response
		}(index)
	}
	wait.Wait()
	close(responses)

	successes := 0
	for response := range responses {
		if response.GetSuccess() {
			successes++
		}
	}
	if successes != requests {
		t.Fatalf("successful concurrent bookings = %d, want %d", successes, requests)
	}
}
