package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/GantangSatria/MyBank-BE/internal/domain"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
)

type RecommendationRepository interface {
	Create(ctx context.Context, rec *domain.Recommendation) error
	FindByID(ctx context.Context, id uint64) (*domain.Recommendation, error)
	FindActiveByUserID(ctx context.Context, userID uint64) ([]domain.Recommendation, error)
	LogClick(ctx context.Context, click *domain.RecommendationClick) error
	GetClickCount(ctx context.Context, recommendationID uint64) (int64, error)
	GetTotalShownAndClicked(ctx context.Context, userID uint64) (int64, int64, error)
}

type recommendationRepository struct {
	db *sql.DB
}

func NewRecommendationRepository(db *sql.DB) RecommendationRepository {
	return &recommendationRepository{db: db}
}

func (r *recommendationRepository) Create(ctx context.Context, rec *domain.Recommendation) error {
	q := `
		INSERT INTO recommendations (
			user_id, type, title, description, image_url, reason, priority,
			is_active, expires_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, 1, ?, NOW(), NOW())`

	var expiresAt interface{}
	if rec.ExpiresAt != nil {
		expiresAt = *rec.ExpiresAt
	}

	result, err := r.db.ExecContext(ctx, q,
		rec.UserID, rec.Type, rec.Title, rec.Description, rec.ImageURL,
		rec.Reason, rec.Priority, expiresAt,
	)
	if err != nil {
		return apperrors.InternalServerError("gagal membuat rekomendasi")
	}

	id, _ := result.LastInsertId()
	rec.ID = uint64(id)
	return nil
}

func (r *recommendationRepository) FindByID(ctx context.Context, id uint64) (*domain.Recommendation, error) {
	q := `
		SELECT id, user_id, type, title, description, image_url, reason,
			priority, is_active, expires_at, created_at, updated_at
		FROM recommendations WHERE id = ?`

	var rec domain.Recommendation
	var imageURL sql.NullString
	var expiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&rec.ID, &rec.UserID, &rec.Type, &rec.Title, &rec.Description,
		&imageURL, &rec.Reason, &rec.Priority, &rec.IsActive,
		&expiresAt, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, apperrors.NotFound("rekomendasi tidak ditemukan")
	}
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil rekomendasi")
	}

	if imageURL.Valid {
		rec.ImageURL = imageURL.String
	}
	if expiresAt.Valid {
		rec.ExpiresAt = &expiresAt.Time
	}

	return &rec, nil
}

func (r *recommendationRepository) FindActiveByUserID(ctx context.Context, userID uint64) ([]domain.Recommendation, error) {
	q := `
		SELECT id, user_id, type, title, description, image_url, reason,
			priority, is_active, expires_at, created_at, updated_at
		FROM recommendations
		WHERE user_id = ? AND is_active = 1
			AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY priority ASC, created_at DESC`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil rekomendasi")
	}
	defer rows.Close()

	var result []domain.Recommendation
	for rows.Next() {
		var rec domain.Recommendation
		var imageURL sql.NullString
		var expiresAt sql.NullTime

		if err := rows.Scan(
			&rec.ID, &rec.UserID, &rec.Type, &rec.Title, &rec.Description,
			&imageURL, &rec.Reason, &rec.Priority, &rec.IsActive,
			&expiresAt, &rec.CreatedAt, &rec.UpdatedAt,
		); err != nil {
			return nil, apperrors.InternalServerError("gagal scan rekomendasi")
		}

		if imageURL.Valid {
			rec.ImageURL = imageURL.String
		}
		if expiresAt.Valid {
			rec.ExpiresAt = &expiresAt.Time
		}

		result = append(result, rec)
	}
	return result, nil
}

func (r *recommendationRepository) LogClick(ctx context.Context, click *domain.RecommendationClick) error {
	q := `INSERT INTO recommendation_clicks (user_id, recommendation_id, clicked_at) VALUES (?, ?, NOW())`
	result, err := r.db.ExecContext(ctx, q, click.UserID, click.RecommendationID)
	if err != nil {
		return apperrors.InternalServerError("gagal mencatat klik rekomendasi")
	}
	id, _ := result.LastInsertId()
	click.ID = uint64(id)
	click.ClickedAt = time.Now()
	return nil
}

func (r *recommendationRepository) GetClickCount(ctx context.Context, recommendationID uint64) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM recommendation_clicks WHERE recommendation_id = ?`,
		recommendationID,
	).Scan(&count)
	if err != nil {
		return 0, apperrors.InternalServerError("gagal menghitung klik")
	}
	return count, nil
}

func (r *recommendationRepository) GetTotalShownAndClicked(ctx context.Context, userID uint64) (int64, int64, error) {
	var totalShown, totalClicked int64

	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM recommendations WHERE user_id = ? AND is_active = 1`,
		userID,
	).Scan(&totalShown)
	if err != nil {
		return 0, 0, apperrors.InternalServerError("gagal menghitung rekomendasi")
	}

	err = r.db.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT rc.recommendation_id)
		 FROM recommendation_clicks rc
		 JOIN recommendations r ON rc.recommendation_id = r.id
		 WHERE rc.user_id = ?`,
		userID,
	).Scan(&totalClicked)
	if err != nil {
		return 0, 0, apperrors.InternalServerError("gagal menghitung klik rekomendasi")
	}

	return totalShown, totalClicked, nil
}
