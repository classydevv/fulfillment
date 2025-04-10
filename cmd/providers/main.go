package main

import (
	"sync"

	pb "github.com/classydevv/metro-fulfillment/pkg/api/providers"
)

type providerId string

type server struct {
	pb.UnimplementedNotesServiceServer
	
	mu sync.RWMutex
	providers map[providerId]*pb.Provider
}

func NewServer() *server {
	return &server{providers: make(map[providerId]*pb.Provider)}
}

func main() {

	
}
