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
// POST /api/v1/transactions
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
// GET /api/v1/transactions
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
// GET /api/v1/transactions/:id
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
// GET /api/v1/transactions/summary
func (h *TransactionHandler) GetSpendingSummary(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	result, err := h.svc.GetSpendingSummary(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success("ringkasan spending berhasil diambil", result))
}
