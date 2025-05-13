package v1

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/classydevv/fulfillment/internal/providers/entity"
	"github.com/classydevv/fulfillment/internal/providers/usecase"
	"github.com/classydevv/fulfillment/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type controllerProvider struct {
	uc usecase.Provider
	l  logger.Interface
	v  *validator.Validate
}

func NewRoutesProvider(apiGroup fiber.Router, uc usecase.Provider, l logger.Interface) {
	r := &controllerProvider{uc, l, validator.New(validator.WithRequiredStructEnabled())}
	providerGroup := apiGroup.Group("/providers")
	{
		providerGroup.Post("", r.providerCreate)
		providerGroup.Get("", r.providerGetAll)
		providerGroup.Put("/:providerId", r.providerUpdate)
		providerGroup.Delete("/:providerId", r.providerDelete)
	}
}

type providerCreateRequest struct {
	ProviderId entity.ProviderId `json:"provider_id" validate:"required" example:"kuper"`
	Name       string            `json:"name" validate:"required" example:"Купер"`
}

type providerCreateResponse struct {
	ProviderId entity.ProviderId `json:"provider_id" example:"kuper"`
}

// @Summary     Create a new provider
// @Description Creates a new delivery provider
// @ID          providerCreate
// @Tags  	    provider
// @Accept      json
// @Produce     json
// @Param       body body providerCreateRequest true "Provider create parameters"
// @Success     201 {object} providerCreateResponse
// @Failure     400 {object} responseError
// @Failure     409 {object} responseError
// @Failure     500 {object} responseError
// @Router      /providers [post]
func (c *controllerProvider) providerCreate(ctx *fiber.Ctx) error {
	var requestBody providerCreateRequest

	if err := ctx.BodyParser(&requestBody); err != nil {
		c.l.Error(fmt.Errorf("http - v1 - providerCreate - bodyParser: %w", err))

		return errorResponse(ctx, http.StatusBadRequest, "bad request")
	}

	if err := c.v.Struct(requestBody); err != nil {
		c.l.Error(fmt.Errorf("http - v1 - providerCreate - validate: %w", err))

		return errorResponse(ctx, http.StatusBadRequest, "bad request")
	}

	providerId, err := c.uc.Create(ctx.UserContext(), &entity.Provider{
		ProviderId: requestBody.ProviderId,
		Name:       requestBody.Name,
	})
	if err != nil {
		c.l.Error(fmt.Errorf("http - v1 - providerCreate - uc.Save: %w", err))

		if errors.Is(err, entity.ErrAlreadyExists) {
			return errorResponse(ctx, http.StatusConflict, fmt.Sprintf("%s: %s", requestBody.ProviderId, entity.ErrAlreadyExists.Error()))
		}

		return errorResponse(ctx, http.StatusInternalServerError, "provider database problems")
	}

	return ctx.Status(http.StatusCreated).JSON(providerCreateResponse{providerId})
}

type providerEntityResponse struct {
	ProviderId entity.ProviderId `json:"provider_id" example:"kuper"`
	Name       string            `json:"name" example:"Купер"`
	CreatedAt  time.Time         `json:"created_at" example:"2025-05-08T06:07:14.810915Z"`
	UpdatedAt  time.Time         `json:"updated_at" example:"2025-05-08T06:07:14.810915Z"`
}

type providerListAllResponse []providerEntityResponse

// @Summary     List all providers
// @Description List all available providers registered in the system
// @ID          providerListAll
// @Tags  	    provider
// @Accept      json
// @Produce     json
// @Success     200 {object} providerListAllResponse
// @Failure     500 {object} responseError
// @Router      /providers [get]
func (c *controllerProvider) providerGetAll(ctx *fiber.Ctx) error {
	providers, err := c.uc.ListAll(ctx.UserContext())
	if err != nil {
		c.l.Error(fmt.Errorf("http - v1 - providerGetAll - uc.ListAll: %w", err))

		return errorResponse(ctx, http.StatusInternalServerError, "provider database problems")
	}

	providersEntityResponse := make([]providerEntityResponse, len(providers))

	for i, p := range providers {
		providersEntityResponse[i] = providerEntityResponse(*p)
	}

	return ctx.Status(http.StatusOK).JSON(providerListAllResponse(providersEntityResponse))
}

type paramProviderId entity.ProviderId

type providerUpdateRequest struct {
	Name string `json:"name" example:"Купер"`
}

type providerUpdateResponse providerEntityResponse

// @Summary     Update a provider
// @Description Updates a delivery provider
// @ID          providerUpdate
// @Tags  	    provider
// @Accept      json
// @Produce     json
// @Param       providerId path string true "Provider ID"
// @Param       body body providerUpdateRequest true "Provider update parameters"
// @Success     200 {object} providerUpdateResponse
// @Failure     400 {object} responseError
// @Failure     404 {object} responseError
// @Failure     500 {object} responseError
// @Router      /providers/{providerId} [put]
func (c *controllerProvider) providerUpdate(ctx *fiber.Ctx) error {
	var providerId = paramProviderId(ctx.Params("providerId"))
	if providerId == "" {
		c.l.Error(fmt.Errorf("http - v1 - providerUpdate - providerId not provided"))

		return errorResponse(ctx, http.StatusBadRequest, "bad request")
	}

	var requestBody providerUpdateRequest
	if err := ctx.BodyParser(&requestBody); err != nil {
		c.l.Error(fmt.Errorf("http - v1 - providerUpdate - bodyParser: %w", err))

		return errorResponse(ctx, http.StatusBadRequest, "bad request")
	}

	if err := c.v.Struct(requestBody); err != nil {
		c.l.Error(fmt.Errorf("http - v1 - providerUpdate - validate: %w", err))

		return errorResponse(ctx, http.StatusBadRequest, "bad request")
	}

	providerUpdated, err := c.uc.Update(ctx.UserContext(),
		entity.ProviderId(providerId),
		&entity.Provider{
			Name: requestBody.Name, // TODO: do not update if field "name" is not passed
		})
	if err != nil {
		c.l.Error(fmt.Errorf("http - v1 - providerUpdate - uc.Update: %w", err))

		if errors.Is(err, entity.ErrNotFound) {
			return errorResponse(ctx, http.StatusNotFound, fmt.Sprintf("%s: %s", providerId, entity.ErrNotFound.Error()))
		}

		return errorResponse(ctx, http.StatusInternalServerError, "provider database problems")
	}

	return ctx.Status(http.StatusOK).JSON(providerUpdateResponse(*providerUpdated))
}

// type providerDeleteResponse struct{}

// @Summary     Delete a provider
// @Description Deletes a delivery provider
// @ID          providerDelete
// @Tags  	    provider
// @Accept      json
// @Produce     json
// @Param       providerId path string true "Provider ID"
// @Success     204
// @Failure     400 {object} responseError
// @Failure     404 {object} responseError
// @Failure     500 {object} responseError
// @Router      /providers/{providerId} [delete]
func (c *controllerProvider) providerDelete(ctx *fiber.Ctx) error {
	var providerId = paramProviderId(ctx.Params("providerId"))
	if providerId == "" {
		c.l.Error(fmt.Errorf("http - v1 - providerDelete - providerId not provided"))

		return errorResponse(ctx, http.StatusBadRequest, "bad request")
	}

	err := c.uc.Delete(ctx.UserContext(), entity.ProviderId(providerId))

	if err != nil {
		c.l.Error(fmt.Errorf("http - v1 - providerDelete - uc.Delete: %w", err))

		if errors.Is(err, entity.ErrNotFound) {
			return errorResponse(ctx, http.StatusNotFound, fmt.Sprintf("%s: %s", providerId, entity.ErrNotFound.Error()))
		}

		return errorResponse(ctx, http.StatusInternalServerError, "provider database problems")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
