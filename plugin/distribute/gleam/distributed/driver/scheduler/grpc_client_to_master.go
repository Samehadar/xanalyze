package scheduler

import (
	"log"
	"time"

	"github.com/sniperkit/xanalyze/plugin/distribute/gleam/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func getResources(master string, request *pb.ComputeRequest) (*pb.AllocationResult, error) {

	grpcConection, err := grpc.Dial(master, grpc.WithInsecure())
	if err != nil {
		log.Printf("fail to dial %s: %v", master, err)
	}
	defer func() {
		time.Sleep(50 * time.Millisecond)
		grpcConection.Close()
	}()

	client := pb.NewGleamMasterClient(grpcConection)

	return client.GetResources(context.Background(), request)
}
