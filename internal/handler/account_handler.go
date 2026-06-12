package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/middleware"
	"github.com/GantangSatria/MyBank-BE/internal/service"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/request"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
	"github.com/GantangSatria/MyBank-BE/pkg/utils"
)

type AccountHandler struct {
	svc service.AccountService
}

func NewAccountHandler(svc service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

// CreateAccount godoc
// @Summary Create a new bank account
// @Description Create a new bank account for the logged-in user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateAccountRequest true "Create Account Request"
// @Success 201 {object} response.Base{data=response.AccountResponse}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req request.CreateAccountRequest
	if err := c.Bind().Body(&req); err != nil {
		return apperrors.BadRequest("invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return err
	}

	res, err := h.svc.CreateAccount(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return c.JSON(response.Success("Account created successfully", res))
}

// GetAccounts godoc
// @Summary Get user accounts
// @Description Get a list of all bank accounts belonging to the logged-in user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base{data=[]response.AccountResponse}
// @Failure 401 {object} response.Base
// @Router /accounts [get]
func (h *AccountHandler) GetAccounts(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	res, err := h.svc.GetAccountsByUserID(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(response.Success("Accounts retrieved successfully", res))
}
