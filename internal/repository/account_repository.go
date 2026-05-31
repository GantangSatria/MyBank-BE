package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GantangSatria/MyBank-BE/internal/domain"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	FindByUserID(ctx context.Context, userID uint64) (*domain.Account, error)
	FindByID(ctx context.Context, accountID uint64) (*domain.Account, error)
	AddBalance(ctx context.Context, accountID uint64, amount float64) error
	SubtractBalance(ctx context.Context, accountID uint64, amount float64) error
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

const accountCols = `id, user_id, account_number, account_type, balance, currency, branch, is_active, created_at, updated_at`

func scanAccount(row interface{ Scan(dest ...interface{}) error }) (*domain.Account, error) {
	var a domain.Account
	var branch sql.NullString
	err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.AccountNumber,
		&a.AccountType,
		&a.Balance,
		&a.Currency,
		&branch,
		&a.IsActive,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *accountRepository) Create(ctx context.Context, a *domain.Account) error {
	query := `
		INSERT INTO accounts (
			user_id, account_number, account_type, balance, currency, branch, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	
	result, err := r.db.ExecContext(
		ctx,
		query,
		a.UserID,
		a.AccountNumber,
		a.AccountType,
		a.Balance,
		a.Currency,
		nil, // branch
		true,
	)
	if err != nil {
		return apperrors.InternalServerError("gagal membuat account")
	}

	id, _ := result.LastInsertId()
	a.ID = uint64(id)
	return nil
}

func (r *accountRepository) FindByUserID(ctx context.Context, userID uint64) (*domain.Account, error) {
	q := `SELECT ` + accountCols + ` FROM accounts WHERE user_id = ? AND is_active = 1 LIMIT 1`
	row := r.db.QueryRowContext(ctx, q, userID)

	a, err := scanAccount(row)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrAccountNotFound
	}
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil data account")
	}
	return a, nil
}

func (r *accountRepository) FindByID(ctx context.Context, accountID uint64) (*domain.Account, error) {
	q := `SELECT ` + accountCols + ` FROM accounts WHERE id = ? AND is_active = 1`
	row := r.db.QueryRowContext(ctx, q, accountID)

	a, err := scanAccount(row)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrAccountNotFound
	}
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil data account")
	}
	return a, nil
}

func (r *accountRepository) AddBalance(ctx context.Context, accountID uint64, amount float64) error {
	q := `UPDATE accounts SET balance = balance + ?, updated_at = NOW() WHERE id = ? AND is_active = 1`
	res, err := r.db.ExecContext(ctx, q, amount, accountID)
	if err != nil {
		return apperrors.InternalServerError("gagal menambah saldo")
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return apperrors.ErrAccountNotFound
	}
	return nil
}

func (r *accountRepository) SubtractBalance(ctx context.Context, accountID uint64, amount float64) error {
	// Memastikan balance tidak minus
	q := `UPDATE accounts SET balance = balance - ?, updated_at = NOW() WHERE id = ? AND balance >= ? AND is_active = 1`
	res, err := r.db.ExecContext(ctx, q, amount, accountID, amount)
	if err != nil {
		return apperrors.InternalServerError("gagal memotong saldo")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		// Cek apakah account ada?
		var exists bool
		_ = r.db.QueryRowContext(ctx, `SELECT 1 FROM accounts WHERE id = ? AND is_active = 1`, accountID).Scan(&exists)
		if !exists {
			return apperrors.ErrAccountNotFound
		}
		// Kalau akun ada tapi row affected = 0, berarti saldo kurang
		return errors.New("insufficient balance")
	}

	return nil
}
