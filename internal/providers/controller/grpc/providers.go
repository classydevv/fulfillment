package grpc

import (
	"context"
	"log"
	"sync"

	pb "github.com/classydevv/fulfillment/pkg/api/providers"
	"github.com/classydevv/fulfillment/pkg/grpc"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	grpc.Server

	pb.UnimplementedProvidersServiceServer

	mu        sync.RWMutex
	providers map[string]*pb.Provider
}

func (s *Server) SaveProvider(_ context.Context, req *pb.SaveProviderRequest) (*pb.SaveProviderResponse, error) {
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
	s.providers[provider.Id] = provider
	s.mu.Unlock()

	return &pb.SaveProviderResponse{
		Id: id,
	}, nil
}

func validateSaveProviderRequest(req *pb.SaveProviderRequest) error {
	id := req.GetId()
	name := req.GetName()
	var violations []*errdetails.BadRequest_FieldViolation
	if id == "" {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "id",
			Description: "empty",
		})
	}
	if name == "" {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "name",
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

func (s *Server) ListProviders(_ context.Context, _ *pb.ListProvidersRequest) (*pb.ListProvidersResponse, error) {
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
