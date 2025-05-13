package v1

import (
	"context"
	"fmt"

	"github.com/classydevv/fulfillment/internal/providers/entity"
	"github.com/classydevv/fulfillment/internal/providers/usecase"
	pb "github.com/classydevv/fulfillment/pkg/api/providers"
	"github.com/classydevv/fulfillment/pkg/logger"
	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type controllerProvider struct {
	pb.UnimplementedProvidersServiceServer

	uc usecase.Provider
	l  logger.Interface
	v  *validator.Validate
}

func NewControllerProvider(s *grpc.Server, uc usecase.Provider, l logger.Interface) {
	c := &controllerProvider{uc: uc, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	{
		pb.RegisterProvidersServiceServer(s, c)
	}

}

func (c *controllerProvider) CreateProvider(ctx context.Context, req *pb.CreateProviderRequest) (*pb.CreateProviderResponse, error) {
	provider := entity.Provider{}
	provider.ProviderId = entity.ProviderId(req.GetId())
	provider.Name = req.GetName()

	if err := validateSaveProviderRequest(req); err != nil {
		c.l.Error(fmt.Errorf("grpc - v1 - CreateProvider - validateSaveProviderRequest: %w", err))

		return nil, fmt.Errorf("grpc - v1 - CreateProvider - validateSaveProviderRequest: %w", err)
	}

	providerId, err := c.uc.Create(ctx, provider)

	if err != nil {
		c.l.Error(fmt.Errorf("grpc - v1 - CreateProvider - uc.Create: %w", err))

		return nil, fmt.Errorf("grpc - v1 - CreateProvider - uc.Create: %w", err)
	}

	return &pb.CreateProviderResponse{
		Id: string(providerId),
	}, nil
}

func validateSaveProviderRequest(req *pb.CreateProviderRequest) error {
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

func (c *controllerProvider) ListAllProviders(ctx context.Context, _ *pb.ListAllProvidersRequest) (*pb.ListAllProvidersResponse, error) {
	providersEntity, err := c.uc.ListAll(ctx)

	if err != nil {
		c.l.Error(fmt.Errorf("grpc - v1 - ListAllProviders - uc.ListAll: %w", err))

		return nil, fmt.Errorf("grpc - v1 - ListAllProviders - uc.ListAll: %w", err)
	}

	providers := make([]*pb.Provider, len(providersEntity))

	for i, provider := range providersEntity {
		providers[i] = &pb.Provider{
			Id: string(provider.ProviderId),
			Name: provider.Name,
		}
	}

	return &pb.ListAllProvidersResponse{
		Providers: providers,
	}, nil
}
