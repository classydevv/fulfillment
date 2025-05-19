package v1

import (
	"context"
	"fmt"

	"github.com/classydevv/fulfillment/internal/providers/entity"
	"github.com/classydevv/fulfillment/internal/providers/usecase"

	pb "github.com/classydevv/fulfillment/pkg/api/providers"
	"github.com/classydevv/fulfillment/pkg/grpcserver"
	"github.com/classydevv/fulfillment/pkg/logger"
	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type controllerProvider struct {
	pb.UnimplementedProvidersServiceServer

	uc usecase.Provider
	l  logger.Interface
	v  *validator.Validate
}

func NewControllerProvider(ctx context.Context, s *grpcserver.Server, uc usecase.Provider, l logger.Interface) {
	c := &controllerProvider{uc: uc, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	{
		pb.RegisterProvidersServiceServer(s.GRPC.Server, c)
		pb.RegisterProvidersServiceHandlerServer(ctx, s.Gateway.Mux, c)
	}
}

func (c *controllerProvider) ProviderCreate(ctx context.Context, req *pb.ProviderCreateRequest) (*pb.ProviderCreateResponse, error) {
	provider := new(entity.Provider)
	provider.ProviderId = entity.ProviderId(req.GetProviderId())
	provider.Name = req.GetName()

	if err := validateProviderCreateRequest(req); err != nil {
		c.l.Error(fmt.Errorf("grpc - v1 - CreateProvider - validateProviderCreateRequest: %w", err))

		return nil, fmt.Errorf("grpc - v1 - CreateProvider - validateProviderCreateRequest: %w", err)
	}

	providerId, err := c.uc.Create(ctx, provider)
	if err != nil {
		c.l.Error(fmt.Errorf("grpc - v1 - ProviderCreate - uc.Create: %w", err))

		return nil, fmt.Errorf("grpc - v1 - ProviderCreate - uc.Create: %w", err)
	}

	return &pb.ProviderCreateResponse{
		ProviderId: string(providerId),
	}, nil
}

func validateProviderCreateRequest(req *pb.ProviderCreateRequest) error {
	providerId := req.GetProviderId()
	name := req.GetName()
	var violations []*errdetails.BadRequest_FieldViolation
	if providerId == "" {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "provider_id",
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

func (c *controllerProvider) ProviderListAll(ctx context.Context, _ *pb.ProviderListAllRequest) (*pb.ProviderListAllResponse, error) {
	providersEntity, err := c.uc.ListAll(ctx)
	if err != nil {
		c.l.Error(fmt.Errorf("grpc - v1 - ProviderListAll - uc.ListAll: %w", err))

		return nil, fmt.Errorf("grpc - v1 - ProviderListAll - uc.ListAll: %w", err)
	}

	providers := make([]*pb.Provider, len(providersEntity))

	for i, provider := range providersEntity {
		providers[i] = &pb.Provider{
			ProviderId: string(provider.ProviderId),
			Name:       provider.Name,
		}
	}

	return &pb.ProviderListAllResponse{
		Providers: providers,
	}, nil
}

func (c *controllerProvider) ProviderUpdate(ctx context.Context, req *pb.ProviderUpdateRequest) (*pb.ProviderUpdateResponse, error) {
	provider := new(entity.Provider)
	provider.ProviderId = entity.ProviderId(req.GetProviderId())
	provider.Name = req.GetName()

	if err := validateProviderUpdateRequest(req); err != nil {
		c.l.Error(fmt.Errorf("grpc - v1 - ProviderUpdate - validateProviderUpdateRequest: %w", err))

		return nil, fmt.Errorf("grpc - v1 - ProviderUpdate - validateProviderUpdateRequest: %w", err)
	}

	providerUpdated, err := c.uc.Update(ctx, provider.ProviderId, provider)
	if err != nil {
		c.l.Error(fmt.Errorf("grpc - v1 - ProviderUpdate - uc.Update: %w", err))

		return nil, fmt.Errorf("grpc - v1 - ProviderUpdate - uc.Update: %w", err)
	}

	return &pb.ProviderUpdateResponse{
		Provider: &pb.Provider{
			ProviderId: string(providerUpdated.ProviderId),
			Name:       providerUpdated.Name,
		},
	}, nil
}

func validateProviderUpdateRequest(req *pb.ProviderUpdateRequest) error {
	providerId := req.GetProviderId()
	var violations []*errdetails.BadRequest_FieldViolation
	if providerId == "" {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "provider_id",
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

func (c *controllerProvider) ProviderDelete(ctx context.Context, req *pb.ProviderDeleteRequest) (*pb.ProviderDeleteResponse, error) {
	providerId := entity.ProviderId(req.GetProviderId())

	if err := validateProviderDeleteRequest(req); err != nil {
		c.l.Error(fmt.Errorf("grpc - v1 - ProviderDelete - validateProviderDeleteRequest: %w", err))

		return nil, fmt.Errorf("grpc - v1 - ProviderDelete - validateProviderDeleteRequest: %w", err)
	}

	err := c.uc.Delete(ctx, providerId)
	if err != nil {
		c.l.Error(fmt.Errorf("grpc - v1 - ProviderDelete - uc.Delete: %w", err))

		return nil, fmt.Errorf("grpc - v1 - ProviderDelete - uc.Delete: %w", err)
	}

	return &pb.ProviderDeleteResponse{}, nil
}

func validateProviderDeleteRequest(req *pb.ProviderDeleteRequest) error {
	providerId := req.GetProviderId()
	var violations []*errdetails.BadRequest_FieldViolation
	if providerId == "" {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       "provider_id",
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
