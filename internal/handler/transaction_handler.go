package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/GantangSatria/MyBank-BE/internal/middleware"
	"github.com/GantangSatria/MyBank-BE/internal/service"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/request"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
	"github.com/GantangSatria/MyBank-BE/pkg/utils"
)

type TransactionHandler struct {
	svc service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction (transfer, payment, topup, withdraw, deposit, qris)
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateTransactionRequest true "Create Transaction Request"
// @Success 201 {object} response.Base{data=response.TransactionResponse}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	var req request.CreateTransactionRequest
	if err := c.Bind().JSON(&req); err != nil {
		return apperrors.BadRequest("body request tidak valid")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return apperrors.BadRequest(err.Error())
	}

	result, err := h.svc.CreateTransaction(c.Context(), userID, &req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success("transaksi berhasil dibuat", result))
}

// GetTransactions godoc
// @Summary Get transactions
// @Description Get a paginated list of transactions for the logged-in user
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} response.Base{data=[]response.TransactionResponse}
// @Failure 401 {object} response.Base
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	txs, total, err := h.svc.GetTransactions(c.Context(), userID, page, limit)
	if err != nil {
		return err
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return c.JSON(response.SuccessWithMeta("daftar transaksi berhasil diambil", txs, &response.Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}))
}

// GetTransactionDetail godoc
// @Summary Get transaction detail
// @Description Get details of a specific transaction by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.Base{data=response.TransactionResponse}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 404 {object} response.Base
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionDetail(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	txID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return apperrors.BadRequest("ID transaksi tidak valid")
	}

	result, err := h.svc.GetTransactionDetail(c.Context(), userID, txID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("detail transaksi berhasil diambil", result))
}

// GetSpendingSummary godoc
// @Summary Get spending summary
// @Description Get a summary of user spending including total spend, favourite category, weekly comparisons, and top merchants
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base{data=response.SpendingSummaryResponse}
// @Failure 401 {object} response.Base
// @Router /transactions/summary [get]
func (h *TransactionHandler) GetSpendingSummary(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	result, err := h.svc.GetSpendingSummary(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("ringkasan spending berhasil diambil", result))
}
