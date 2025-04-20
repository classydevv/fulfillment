package main

import (
	"context"
	"log"
	"time"

	pb "github.com/classydevv/fulfillment/pkg/api/providers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	cli := pb.NewProvidersServiceClient(conn)

	// SaveProvider
	{
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		res, err := cli.SaveProvider(ctx, &pb.SaveProviderRequest{Id: "kuper", Name: "Купер"})
		if err != nil {
			log.Fatalf("SaveProvider failed: %v", err)
		}
		log.Printf("SaveProvider success: providerId: %v", res.GetId())
	}

	// ListProviders
	{
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		res, err := cli.ListProviders(ctx, &pb.ListProvidersRequest{})
		if err != nil {
			log.Fatalf("ListProviders failed: error response: %v", err)
		}
		providers, err := protojson.Marshal(res)
		if err != nil {
			log.Fatalf("ListProviders failed: error marshaling: %v", err)
		}
		log.Printf("ListProviders success: providers: %v", string(providers))
	}
}
