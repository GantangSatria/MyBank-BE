package repository

import (
	"context"
	"database/sql"

	"github.com/GantangSatria/MyBank-BE/internal/domain"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
)

type FeatureClickRepository interface {
	Upsert(ctx context.Context, userID uint64, featureName string) error
	FindByUserID(ctx context.Context, userID uint64) ([]domain.FeatureClick, error)
}

type featureClickRepository struct {
	db *sql.DB
}

func NewFeatureClickRepository(db *sql.DB) FeatureClickRepository {
	return &featureClickRepository{db: db}
}

func (r *featureClickRepository) Upsert(ctx context.Context, userID uint64, featureName string) error {
	q := `
		INSERT INTO feature_clicks (user_id, feature_name, click_count, last_clicked, created_at, updated_at)
		VALUES (?, ?, 1, NOW(), NOW(), NOW())
		ON DUPLICATE KEY UPDATE click_count = click_count + 1, last_clicked = NOW(), updated_at = NOW()`

	_, err := r.db.ExecContext(ctx, q, userID, featureName)
	if err != nil {
		return apperrors.InternalServerError("gagal mencatat klik fitur")
	}
	return nil
}

func (r *featureClickRepository) FindByUserID(ctx context.Context, userID uint64) ([]domain.FeatureClick, error) {
	q := `
		SELECT id, user_id, feature_name, click_count, last_clicked, created_at, updated_at
		FROM feature_clicks
		WHERE user_id = ?
		ORDER BY click_count DESC`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil data klik fitur")
	}
	defer rows.Close()

	var result []domain.FeatureClick
	for rows.Next() {
		var fc domain.FeatureClick
		if err := rows.Scan(
			&fc.ID, &fc.UserID, &fc.FeatureName, &fc.ClickCount,
			&fc.LastClicked, &fc.CreatedAt, &fc.UpdatedAt,
		); err != nil {
			return nil, apperrors.InternalServerError("gagal scan klik fitur")
		}
		result = append(result, fc)
	}
	return result, nil
}
