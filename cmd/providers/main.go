package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/classydevv/metro-fulfillment/pkg/api/providers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type providerId string

type server struct {
	pb.UnimplementedProvidersServiceServer
	
	mu sync.RWMutex
	providers map[providerId]*pb.Provider
}

func NewServer() *server {
	return &server{providers: make(map[providerId]*pb.Provider)}
}

func (s *server) SaveProvider(_ context.Context, req *pb.SaveProviderRequest) (*pb.SaveProviderResponse, error) {
	provider := new(pb.Provider)
	id := req.GetId()
	name := req.GetName()
	provider.Id = id
	provider.Name = name

	if err := validateSaveProviderRequest(req); err != nil {
		return nil, err
	}

	log.Printf("SaveProvider: received: %s", req.GetId())

	s.mu.Lock()
	s.providers[providerId(provider.Id)] = provider
	s.mu.Unlock()

	return &pb.SaveProviderResponse{
		Id: id,
	}, nil
}

func validateSaveProviderRequest(req *pb.SaveProviderRequest) error {
	id := req.GetId()
	name := req.GetName()
	var violations []*errdetails.BadRequest_FieldViolation
	if len(id) == 0 {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field: "id",
			Description: "empty",
		})
	}
	if len(name) == 0 {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field: "name",
			Description: "empty",
		})
	}
	if len(violations) > 0 {
		st, err := status.New(codes.InvalidArgument, codes.InvalidArgument.String()).WithDetails(
			&errdetails.BadRequest{
				FieldViolations: violations,
			})
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		return st.Err()
	}
	return nil
}

func (s *server) ListProviders(_ context.Context, req *pb.ListProvidersRequest) (*pb.ListProvidersResponse, error) {
	log.Printf("ListProviders: received")

	s.mu.RLock()
	defer s.mu.RUnlock()

	providers := make([]*pb.Provider, 0, len(s.providers))
	for _, provider := range s.providers {
		providers = append(providers, provider)
	}

	return &pb.ListProvidersResponse{
		Providers: providers,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProvidersServiceServer(s, NewServer())

	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
	
}