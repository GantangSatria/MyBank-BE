package repository

import (
	"context"
	"database/sql"

	"github.com/GantangSatria/MyBank-BE/internal/domain"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
)

type MerchantRepository interface {
	Create(ctx context.Context, m *domain.Merchant) error
	FindAll(ctx context.Context) ([]domain.Merchant, error)
	FindByMerchantID(ctx context.Context, merchantID string) (*domain.Merchant, error)
	FindByCategory(ctx context.Context, category string) ([]domain.Merchant, error)
	Update(ctx context.Context, m *domain.Merchant) error
	Delete(ctx context.Context, id uint64) error
}

type merchantRepository struct {
	db *sql.DB
}

func NewMerchantRepository(db *sql.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

const merchantCols = `
	id, merchant_id, merchant_name, merchant_category,
	merchant_city, merchant_type, merchant_status,
	created_at, updated_at`

func scanMerchant(row interface{ Scan(dest ...interface{}) error }) (*domain.Merchant, error) {
	var m domain.Merchant
	var city sql.NullString

	err := row.Scan(
		&m.ID, &m.MerchantID, &m.MerchantName, &m.MerchantCategory,
		&city, &m.MerchantType, &m.MerchantStatus,
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if city.Valid {
		m.MerchantCity = city.String
	}

	return &m, nil
}

func (r *merchantRepository) Create(ctx context.Context, m *domain.Merchant) error {
	query := `
		INSERT INTO merchants (
			merchant_id, merchant_name, merchant_category,
			merchant_city, merchant_type, merchant_status,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`

	result, err := r.db.ExecContext(ctx, query,
		m.MerchantID, m.MerchantName, m.MerchantCategory,
		nullStr(m.MerchantCity), m.MerchantType, m.MerchantStatus,
	)
	if err != nil {
		return apperrors.InternalServerError("gagal membuat merchant")
	}

	id, _ := result.LastInsertId()
	m.ID = uint64(id)
	return nil
}

func (r *merchantRepository) FindAll(ctx context.Context) ([]domain.Merchant, error) {
	q := `SELECT ` + merchantCols + ` FROM merchants ORDER BY merchant_id ASC`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil daftar merchant")
	}
	defer rows.Close()

	var merchants []domain.Merchant
	for rows.Next() {
		m, err := scanMerchant(rows)
		if err != nil {
			return nil, apperrors.InternalServerError("gagal scan merchant")
		}
		merchants = append(merchants, *m)
	}
	return merchants, nil
}

func (r *merchantRepository) FindByMerchantID(ctx context.Context, merchantID string) (*domain.Merchant, error) {
	q := `SELECT ` + merchantCols + ` FROM merchants WHERE merchant_id = ?`
	row := r.db.QueryRowContext(ctx, q, merchantID)

	m, err := scanMerchant(row)
	if err == sql.ErrNoRows {
		return nil, apperrors.NotFound("merchant tidak ditemukan")
	}
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil merchant")
	}
	return m, nil
}

func (r *merchantRepository) FindByCategory(ctx context.Context, category string) ([]domain.Merchant, error) {
	q := `SELECT ` + merchantCols + ` FROM merchants WHERE merchant_category = ? AND merchant_status = 'aktif' ORDER BY merchant_name ASC`

	rows, err := r.db.QueryContext(ctx, q, category)
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil merchant")
	}
	defer rows.Close()

	var merchants []domain.Merchant
	for rows.Next() {
		m, err := scanMerchant(rows)
		if err != nil {
			return nil, apperrors.InternalServerError("gagal scan merchant")
		}
		merchants = append(merchants, *m)
	}
	return merchants, nil
}

func (r *merchantRepository) Update(ctx context.Context, m *domain.Merchant) error {
	q := `
		UPDATE merchants
		SET merchant_name = ?, merchant_category = ?,
			merchant_city = ?, merchant_type = ?, merchant_status = ?,
			updated_at = NOW()
		WHERE id = ?`

	res, err := r.db.ExecContext(ctx, q,
		m.MerchantName, m.MerchantCategory,
		nullStr(m.MerchantCity), m.MerchantType, m.MerchantStatus,
		m.ID,
	)
	if err != nil {
		return apperrors.InternalServerError("gagal update merchant")
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return apperrors.NotFound("merchant tidak ditemukan")
	}
	return nil
}

func (r *merchantRepository) Delete(ctx context.Context, id uint64) error {
	q := `DELETE FROM merchants WHERE id = ?`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return apperrors.InternalServerError("gagal menghapus merchant")
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return apperrors.NotFound("merchant tidak ditemukan")
	}
	return nil
}
