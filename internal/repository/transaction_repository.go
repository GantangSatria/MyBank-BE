package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/GantangSatria/MyBank-BE/internal/domain"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	FindByID(ctx context.Context, id uint64) (*domain.Transaction, error)
	FindByUserID(ctx context.Context, userID uint64, page, limit int) ([]domain.Transaction, int64, error)
	GetCategorySpending(ctx context.Context, userID uint64, from, to time.Time) ([]domain.CategorySpending, error)
	GetTopMerchants(ctx context.Context, userID uint64, limit int) ([]domain.MerchantSummary, error)
	GetTotalSpendAndCount(ctx context.Context, userID uint64) (float64, int64, error)
	GetFavCategoryAndMethod(ctx context.Context, userID uint64) (string, string, error)
	GetWeeklySpendingByCategory(ctx context.Context, userID uint64) ([]domain.WeeklySpending, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

const txCols = `
	id, user_id, account_id, reference_number, type, status,
	amount, fee, balance_before, balance_after,
	destination_account_number, destination_bank_code, destination_name,
	merchant_id_ref, merchant_name, merchant_category, merchant_location,
	channel, description, note, fail_reason,
	is_recommended, recommendation_id,
	transacted_at, created_at, updated_at`

func scanTransaction(row interface{ Scan(dest ...interface{}) error }) (*domain.Transaction, error) {
	var t domain.Transaction
	var destAccNum, destBankCode, destName sql.NullString
	var merchName, merchCat, merchLoc, channel sql.NullString
	var desc, note, failReason sql.NullString
	var merchIDRef, recID sql.NullInt64

	err := row.Scan(
		&t.ID, &t.UserID, &t.AccountID, &t.ReferenceNumber, &t.Type, &t.Status,
		&t.Amount, &t.Fee, &t.BalanceBefore, &t.BalanceAfter,
		&destAccNum, &destBankCode, &destName,
		&merchIDRef, &merchName, &merchCat, &merchLoc,
		&channel, &desc, &note, &failReason,
		&t.IsRecommended, &recID,
		&t.TransactedAt, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if destAccNum.Valid {
		t.DestinationAccountNumber = destAccNum.String
	}
	if destBankCode.Valid {
		t.DestinationBankCode = destBankCode.String
	}
	if destName.Valid {
		t.DestinationName = destName.String
	}
	if merchIDRef.Valid {
		id := uint64(merchIDRef.Int64)
		t.MerchantID = &id
	}
	if merchName.Valid {
		t.MerchantName = merchName.String
	}
	if merchCat.Valid {
		t.MerchantCategory = merchCat.String
	}
	if merchLoc.Valid {
		t.MerchantLocation = merchLoc.String
	}
	if channel.Valid {
		t.Channel = channel.String
	}
	if desc.Valid {
		t.Description = desc.String
	}
	if note.Valid {
		t.Note = note.String
	}
	if failReason.Valid {
		t.FailReason = failReason.String
	}
	if recID.Valid {
		id := uint64(recID.Int64)
		t.RecommendationID = &id
	}

	return &t, nil
}


func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO transactions (
			user_id, account_id, reference_number, type, status,
			amount, fee, balance_before, balance_after,
			destination_account_number, destination_bank_code, destination_name,
			merchant_id_ref, merchant_name, merchant_category, merchant_location,
			channel, description, note,
			is_recommended, recommendation_id,
			transacted_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), NOW())`

	var recID interface{}
	if tx.RecommendationID != nil {
		recID = *tx.RecommendationID
	}

	var merchIDRef interface{}
	if tx.MerchantID != nil {
		merchIDRef = *tx.MerchantID
	}

	result, err := r.db.ExecContext(ctx, query,
		tx.UserID, tx.AccountID, tx.ReferenceNumber, tx.Type, tx.Status,
		tx.Amount, tx.Fee, tx.BalanceBefore, tx.BalanceAfter,
		nullStr(tx.DestinationAccountNumber), nullStr(tx.DestinationBankCode), nullStr(tx.DestinationName),
		merchIDRef, nullStr(tx.MerchantName), nullStr(tx.MerchantCategory), nullStr(tx.MerchantLocation),
		nullStr(tx.Channel), nullStr(tx.Description), nullStr(tx.Note),
		tx.IsRecommended, recID,
	)

	if err != nil {
		return apperrors.InternalServerError("gagal membuat transaksi")
	}

	id, _ := result.LastInsertId()
	tx.ID = uint64(id)
	return nil
}

func (r *transactionRepository) FindByID(ctx context.Context, id uint64) (*domain.Transaction, error) {
	q := `SELECT ` + txCols + ` FROM transactions WHERE id = ?`
	row := r.db.QueryRowContext(ctx, q, id)

	t, err := scanTransaction(row)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrTransactionNotFound
	}
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil transaksi")
	}
	return t, nil
}

func (r *transactionRepository) FindByUserID(ctx context.Context, userID uint64, page, limit int) ([]domain.Transaction, int64, error) {
	// Count total
	var total int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(1) FROM transactions WHERE user_id = ?`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.InternalServerError("gagal menghitung transaksi")
	}

	offset := (page - 1) * limit
	q := `SELECT ` + txCols + ` FROM transactions WHERE user_id = ? ORDER BY transacted_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, apperrors.InternalServerError("gagal mengambil transaksi")
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, 0, apperrors.InternalServerError("gagal scan transaksi")
		}
		transactions = append(transactions, *t)
	}

	return transactions, total, nil
}

func (r *transactionRepository) GetCategorySpending(ctx context.Context, userID uint64, from, to time.Time) ([]domain.CategorySpending, error) {
	q := `
		SELECT merchant_category, SUM(amount) as total_amount, COUNT(*) as trx_count
		FROM transactions
		WHERE user_id = ? AND status = 'SUCCESS' AND transacted_at BETWEEN ? AND ?
			AND merchant_category IS NOT NULL AND merchant_category != ''
		GROUP BY merchant_category
		ORDER BY total_amount DESC`

	rows, err := r.db.QueryContext(ctx, q, userID, from, to)
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil spending per kategori")
	}
	defer rows.Close()

	var result []domain.CategorySpending
	for rows.Next() {
		var cs domain.CategorySpending
		if err := rows.Scan(&cs.MerchantCategory, &cs.TotalAmount, &cs.TransactionCount); err != nil {
			return nil, apperrors.InternalServerError("gagal scan spending kategori")
		}
		result = append(result, cs)
	}
	return result, nil
}

func (r *transactionRepository) GetTopMerchants(ctx context.Context, userID uint64, limit int) ([]domain.MerchantSummary, error) {
	q := `
		SELECT merchant_name, merchant_category, SUM(amount) as total_amount, COUNT(*) as trx_count
		FROM transactions
		WHERE user_id = ? AND status = 'SUCCESS'
			AND merchant_name IS NOT NULL AND merchant_name != ''
		GROUP BY merchant_name, merchant_category
		ORDER BY trx_count DESC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, q, userID, limit)
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil top merchants")
	}
	defer rows.Close()

	var result []domain.MerchantSummary
	for rows.Next() {
		var ms domain.MerchantSummary
		if err := rows.Scan(&ms.MerchantName, &ms.MerchantCategory, &ms.TotalAmount, &ms.TransactionCount); err != nil {
			return nil, apperrors.InternalServerError("gagal scan merchant")
		}
		result = append(result, ms)
	}
	return result, nil
}

func (r *transactionRepository) GetTotalSpendAndCount(ctx context.Context, userID uint64) (float64, int64, error) {
	var totalSpend float64
	var totalTrx int64

	err := r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(amount), 0), COUNT(*) FROM transactions WHERE user_id = ? AND status = 'SUCCESS'`,
		userID,
	).Scan(&totalSpend, &totalTrx)
	if err != nil {
		return 0, 0, apperrors.InternalServerError("gagal mengambil total spend")
	}
	return totalSpend, totalTrx, nil
}

func (r *transactionRepository) GetFavCategoryAndMethod(ctx context.Context, userID uint64) (string, string, error) {
	// Favourite category (by count)
	var favCat sql.NullString
	_ = r.db.QueryRowContext(ctx,
		`SELECT merchant_category FROM transactions
		 WHERE user_id = ? AND status = 'SUCCESS' AND merchant_category IS NOT NULL AND merchant_category != ''
		 GROUP BY merchant_category ORDER BY COUNT(*) DESC LIMIT 1`,
		userID,
	).Scan(&favCat)

	// Favourite method/type (by count)
	var favMethod sql.NullString
	_ = r.db.QueryRowContext(ctx,
		`SELECT type FROM transactions
		 WHERE user_id = ? AND status = 'SUCCESS'
		 GROUP BY type ORDER BY COUNT(*) DESC LIMIT 1`,
		userID,
	).Scan(&favMethod)

	cat := ""
	if favCat.Valid {
		cat = favCat.String
	}
	method := ""
	if favMethod.Valid {
		method = favMethod.String
	}

	return cat, method, nil
}

func (r *transactionRepository) GetWeeklySpendingByCategory(ctx context.Context, userID uint64) ([]domain.WeeklySpending, error) {
	// Ambil spending minggu ini dan minggu lalu per kategori
	q := `
		SELECT
			cat,
			SUM(CASE WHEN week_offset = 0 THEN amount ELSE 0 END) AS current_week,
			SUM(CASE WHEN week_offset = 1 THEN amount ELSE 0 END) AS previous_week
		FROM (
			SELECT
				merchant_category AS cat,
				amount,
				CASE
					WHEN transacted_at >= DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY)
						THEN 0
					WHEN transacted_at >= DATE_SUB(CURDATE(), INTERVAL (WEEKDAY(CURDATE()) + 7) DAY)
						AND transacted_at < DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY)
						THEN 1
					ELSE 2
				END AS week_offset
			FROM transactions
			WHERE user_id = ? AND status = 'SUCCESS'
				AND merchant_category IS NOT NULL AND merchant_category != ''
				AND transacted_at >= DATE_SUB(CURDATE(), INTERVAL (WEEKDAY(CURDATE()) + 7) DAY)
		) sub
		WHERE week_offset <= 1
		GROUP BY cat
		ORDER BY current_week DESC`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil weekly spending")
	}
	defer rows.Close()

	var result []domain.WeeklySpending
	for rows.Next() {
		var ws domain.WeeklySpending
		if err := rows.Scan(&ws.Category, &ws.CurrentWeekSpend, &ws.PreviousWeekSpend); err != nil {
			return nil, apperrors.InternalServerError("gagal scan weekly spending")
		}
		result = append(result, ws)
	}
	return result, nil
}

// GenerateRefNumber membuat reference number unik
func GenerateRefNumber(txType string) string {
	return fmt.Sprintf("TRX-%s-%d", txType, time.Now().UnixNano())
}
