package v1

import (
	"fmt"
	"net/http"

	"github.com/classydevv/fulfillment/internal/providers/entity"
	"github.com/classydevv/fulfillment/internal/providers/usecase"
	"github.com/classydevv/fulfillment/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type providerRoutes struct {
	uc usecase.Provider
	l logger.Interface
	v *validator.Validate
}

func NewProviderRoutes(apiGroup fiber.Router, uc usecase.Provider, l logger.Interface) {
	r := &providerRoutes{uc, l, validator.New(validator.WithRequiredStructEnabled())}
	providerGroup := apiGroup.Group("/providers")
	{
		providerGroup.Post("", r.providerCreate)
		providerGroup.Get("", r.providerGetAll)
	}
}

type providerCreateRequest struct {
	Id entity.ProviderId `json:"id" validate:"required" example:"kuper"`
	Name string `json:"name" validate:"required" example:"Купер"`
}

type providerCreateResponse struct {
	Id entity.ProviderId `json:"id"`
}

// @Summary     Create a new provider
// @Description Creates a new delivery provider
// @ID          providerCreate
// @Tags  	    provider
// @Accept      json
// @Produce     json
// @Param       request body providerCreateRequest true "Provider parameters"
// @Success     200 {object} providerCreateResponse
// @Failure     400 {object} responseError
// @Failure     500 {object} responseError
// @Router      /providers [post]
func (r *providerRoutes) providerCreate(ctx *fiber.Ctx) error {
	var request providerCreateRequest

	if err := ctx.BodyParser(&request); err != nil {
		r.l.Error(fmt.Errorf("http - v1 - providerCreate - bodyParser: %w", err))

		return errorResponse(ctx, http.StatusBadRequest, "bad request")
	}

	if err := r.v.Struct(request); err != nil {
		r.l.Error(fmt.Errorf("http - v1 - providerCreate - validate: %w", err))

		return errorResponse(ctx, http.StatusBadRequest, "bad request")
	}

	providerId, err := r.uc.Create(ctx.UserContext(), entity.Provider{
		Id: request.Id,
		Name: request.Name,
	})
	if err != nil {
		r.l.Error(fmt.Errorf("http - v1 - providerCreate - uc.Save: %w", err))

		return errorResponse(ctx, http.StatusInternalServerError, "provider database problems")
	}

	return ctx.Status(http.StatusOK).JSON(providerCreateResponse{providerId})
}

type providerListAllResponse struct {
	Providers []entity.Provider `json:"providers"`
}

// @Summary     List all providers
// @Description List all available providers registered in the system
// @ID          providerListAll
// @Tags  	    provider
// @Accept      json
// @Produce     json
// @Success     200 {object} providerListAllResponse
// @Failure     500 {object} responseError
// @Router      /providers [get]
func (r *providerRoutes) providerGetAll(ctx *fiber.Ctx) error {
	providers, err := r.uc.ListAll(ctx.UserContext())
	if err != nil {
		r.l.Error(fmt.Errorf("http - v1 - providerGetAll - uc.ListAll: %w", err))

		return errorResponse(ctx, http.StatusInternalServerError, "provider database problems")
	}

	return ctx.Status(http.StatusOK).JSON(providerListAllResponse{providers})
}