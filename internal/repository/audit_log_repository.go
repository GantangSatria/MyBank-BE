package repository

import (
	"context"
	"database/sql"

	"github.com/GantangSatria/MyBank-BE/internal/domain"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.AuditLog) error
	FindByUserID(ctx context.Context, userID uint64, page, limit int) ([]domain.AuditLog, int64, error)
}

type auditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	q := `
		INSERT INTO audit_logs (user_id, action, detail, ip_address, created_at)
		VALUES (?, ?, ?, ?, NOW())`

	result, err := r.db.ExecContext(ctx, q, log.UserID, log.Action, log.Detail, log.IPAddress)
	if err != nil {
		return apperrors.InternalServerError("gagal mencatat audit log")
	}
	id, _ := result.LastInsertId()
	log.ID = uint64(id)
	return nil
}

func (r *auditLogRepository) FindByUserID(ctx context.Context, userID uint64, page, limit int) ([]domain.AuditLog, int64, error) {
	var total int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(1) FROM audit_logs WHERE user_id = ?`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, apperrors.InternalServerError("gagal menghitung audit log")
	}

	offset := (page - 1) * limit
	q := `SELECT id, user_id, action, detail, ip_address, created_at
		  FROM audit_logs WHERE user_id = ?
		  ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, 0, apperrors.InternalServerError("gagal mengambil audit log")
	}
	defer rows.Close()

	var result []domain.AuditLog
	for rows.Next() {
		var al domain.AuditLog
		var ipAddr sql.NullString
		if err := rows.Scan(&al.ID, &al.UserID, &al.Action, &al.Detail, &ipAddr, &al.CreatedAt); err != nil {
			return nil, 0, apperrors.InternalServerError("gagal scan audit log")
		}
		if ipAddr.Valid {
			al.IPAddress = ipAddr.String
		}
		result = append(result, al)
	}
	return result, total, nil
}
