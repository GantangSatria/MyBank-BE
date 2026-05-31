package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/GantangSatria/MyBank-BE/internal/domain"
	apperrors "github.com/GantangSatria/MyBank-BE/pkg/errors"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	FindByID(ctx context.Context, id uint64) (*domain.User, error)
	FindByIDWithAuth(ctx context.Context, id uint64) (*domain.User, error)
	FindByIDWithPIN(ctx context.Context, id uint64) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByPhone(ctx context.Context, phone string) (*domain.User, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)

	UpdateProfile(ctx context.Context, id uint64, u *domain.User) error
	UpdatePhone(ctx context.Context, id uint64, u *domain.User) error
	UpdatePIN(ctx context.Context, id uint64, hashedPIN string) error
	UpdatePassword(ctx context.Context, id uint64, hashedPassword string) error
	UpdatePersonalizationConsent(ctx context.Context, id uint64, enabled bool) error

	SoftDelete(ctx context.Context, id uint64) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

const userCols = `
	id, name, email, phone,
	gender, date_of_birth, occupation, marital_status, segment,
	monthly_income_range,
	is_personalization_enabled, is_active,
	last_login_at, created_at, updated_at`

const userColsWithPIN = `
	id, name, email, phone, pin,
	gender, date_of_birth, occupation, marital_status, segment,
	monthly_income_range,
	is_personalization_enabled, is_active,
	last_login_at, created_at, updated_at`

const userColsWithAuth = `
	id, name, email, phone, password, pin,
	gender, date_of_birth, occupation, marital_status, segment,
	monthly_income_range,
	is_personalization_enabled, is_active,
	last_login_at, created_at, updated_at`

func nullStrPtr(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

func toStrPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func scanUser(row rowScanner) (*domain.User, error) {
	var u domain.User
	var dob, lastLogin, createdAt, updatedAt sql.NullTime
	var name, phone, gender, occupation, maritalStatus, segment, monthlyIncomeRange sql.NullString
	var isPersonalizationEnabled, isActive sql.NullBool

	err := row.Scan(
		&u.ID,
		&name,
		&u.Email,
		&phone,
		&gender,
		&dob,
		&occupation,
		&maritalStatus,
		&segment,
		&monthlyIncomeRange,
		&isPersonalizationEnabled,
		&isActive,
		&lastLogin,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	u.Name = toStrPtr(name)
	u.Phone = toStrPtr(phone)
	u.Gender = toStrPtr(gender)
	u.Occupation = toStrPtr(occupation)
	u.MaritalStatus = toStrPtr(maritalStatus)
	u.Segment = toStrPtr(segment)
	u.MonthlyIncomeRange = toStrPtr(monthlyIncomeRange)
	u.IsPersonalizationEnabled = isPersonalizationEnabled.Valid && isPersonalizationEnabled.Bool
	u.IsActive = isActive.Valid && isActive.Bool

	if dob.Valid {
		u.DateOfBirth = &dob.Time
	}
	if lastLogin.Valid {
		u.LastLoginAt = &lastLogin.Time
	}
	if createdAt.Valid {
		u.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		u.UpdatedAt = updatedAt.Time
	}

	return &u, nil
}

func scanUserWithPIN(row rowScanner) (*domain.User, error) {
	var u domain.User
	var dob, lastLogin, createdAt, updatedAt sql.NullTime
	var name, phone, pin, gender, occupation, maritalStatus, segment, monthlyIncomeRange sql.NullString
	var isPersonalizationEnabled, isActive sql.NullBool

	err := row.Scan(
		&u.ID,
		&name,
		&u.Email,
		&phone,
		&pin,
		&gender,
		&dob,
		&occupation,
		&maritalStatus,
		&segment,
		&monthlyIncomeRange,
		&isPersonalizationEnabled,
		&isActive,
		&lastLogin,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	u.Name = toStrPtr(name)
	u.Phone = toStrPtr(phone)
	u.PIN = toStrPtr(pin)
	u.Gender = toStrPtr(gender)
	u.Occupation = toStrPtr(occupation)
	u.MaritalStatus = toStrPtr(maritalStatus)
	u.Segment = toStrPtr(segment)
	u.MonthlyIncomeRange = toStrPtr(monthlyIncomeRange)
	u.IsPersonalizationEnabled = isPersonalizationEnabled.Valid && isPersonalizationEnabled.Bool
	u.IsActive = isActive.Valid && isActive.Bool

	if dob.Valid {
		u.DateOfBirth = &dob.Time
	}
	if lastLogin.Valid {
		u.LastLoginAt = &lastLogin.Time
	}
	if createdAt.Valid {
		u.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		u.UpdatedAt = updatedAt.Time
	}

	return &u, nil
}

func scanUserWithAuth(row rowScanner) (*domain.User, error) {
	var u domain.User
	var dob, lastLogin, createdAt, updatedAt sql.NullTime
	var name, phone, pin, gender, occupation, maritalStatus, segment, monthlyIncomeRange sql.NullString
	var isPersonalizationEnabled, isActive sql.NullBool

	err := row.Scan(
		&u.ID,
		&name,
		&u.Email,
		&phone,
		&u.Password,
		&pin,
		&gender,
		&dob,
		&occupation,
		&maritalStatus,
		&segment,
		&monthlyIncomeRange,
		&isPersonalizationEnabled,
		&isActive,
		&lastLogin,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	u.Name = toStrPtr(name)
	u.Phone = toStrPtr(phone)
	u.PIN = toStrPtr(pin)
	u.Gender = toStrPtr(gender)
	u.Occupation = toStrPtr(occupation)
	u.MaritalStatus = toStrPtr(maritalStatus)
	u.Segment = toStrPtr(segment)
	u.MonthlyIncomeRange = toStrPtr(monthlyIncomeRange)
	u.IsPersonalizationEnabled = isPersonalizationEnabled.Valid && isPersonalizationEnabled.Bool
	u.IsActive = isActive.Valid && isActive.Bool

	if dob.Valid {
		u.DateOfBirth = &dob.Time
	}
	if lastLogin.Valid {
		u.LastLoginAt = &lastLogin.Time
	}
	if createdAt.Valid {
		u.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		u.UpdatedAt = updatedAt.Time
	}

	return &u, nil
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (
			name,
			email,
			phone,
			password,
			pin,
			date_of_birth,
			occupation,
			is_personalization_enabled,
			is_active,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, NOW(), NOW())`

	var dob interface{}
	if u.DateOfBirth != nil {
		dob = u.DateOfBirth.Format("2006-01-02")
	}

	result, err := r.db.ExecContext(
		ctx,
		query,
		nullStrPtr(u.Name),
		u.Email,
		nullStrPtr(u.Phone),
		u.Password,
		nullStrPtr(u.PIN),
		dob,
		nullStrPtr(u.Occupation),
		u.IsPersonalizationEnabled,
	)

	if err != nil {
		return apperrors.InternalServerError("gagal membuat user")
	}

	id, _ := result.LastInsertId()
	u.ID = uint64(id)

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (*domain.User, error) {
	q := `SELECT ` + userCols + ` FROM users WHERE id = ? AND deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, q, id)

	u, err := scanUser(row)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrUserNotFound
	}
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil user")
	}

	return u, nil
}

func (r *userRepository) FindByIDWithAuth(ctx context.Context, id uint64) (*domain.User, error) {
	q := `SELECT ` + userColsWithAuth + ` FROM users WHERE id = ? AND deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, q, id)

	u, err := scanUserWithAuth(row)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrUserNotFound
	}
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil user")
	}

	return u, nil
}

func (r *userRepository) FindByIDWithPIN(ctx context.Context, id uint64) (*domain.User, error) {
	q := `SELECT ` + userColsWithPIN + ` FROM users WHERE id = ? AND deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, q, id)

	u, err := scanUserWithPIN(row)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrUserNotFound
	}
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil user")
	}

	return u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := `SELECT ` + userColsWithAuth + ` FROM users WHERE email = ? AND deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, q, email)

	u, err := scanUserWithAuth(row)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrUserNotFound
	}
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil user")
	}

	return u, nil
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	q := `SELECT ` + userCols + ` FROM users WHERE phone = ? AND deleted_at IS NULL`

	row := r.db.QueryRowContext(ctx, q, phone)

	u, err := scanUserWithAuth(row)
	if err == sql.ErrNoRows {
		return nil, apperrors.ErrUserNotFound
	}
	if err != nil {
		return nil, apperrors.InternalServerError("gagal mengambil user")
	}

	return u, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int

	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(1) FROM users WHERE email = ? AND deleted_at IS NULL`,
		email,
	).Scan(&count)

	if err != nil {
		return false, apperrors.InternalServerError("gagal cek email")
	}

	return count > 0, nil
}

func (r *userRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var count int

	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(1) FROM users WHERE phone = ? AND deleted_at IS NULL`,
		phone,
	).Scan(&count)

	if err != nil {
		return false, apperrors.InternalServerError("gagal cek nomor HP")
	}

	return count > 0, nil
}

func (r *userRepository) UpdateProfile(ctx context.Context, id uint64, u *domain.User) error {
	var dob interface{}
	if u.DateOfBirth != nil {
		dob = u.DateOfBirth.Format("2006-01-02")
	}

	q := `
		UPDATE users
		SET name = ?, occupation = ?, date_of_birth = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL`

	res, err := r.db.ExecContext(
		ctx,
		q,
		nullStrPtr(u.Name),
		nullStrPtr(u.Occupation),
		dob,
		id,
	)

	if err != nil {
		return apperrors.InternalServerError("gagal update profil")
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) UpdatePIN(ctx context.Context, id uint64, hashedPIN string) error {
	q := `
		UPDATE users
		SET pin = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, q, hashedPIN, id)

	if err != nil {
		return apperrors.InternalServerError("gagal update PIN")
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uint64, hashedPassword string) error {
	q := `
		UPDATE users
		SET password = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, q, hashedPassword, id)

	if err != nil {
		return apperrors.InternalServerError("gagal update password")
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) UpdatePhone(ctx context.Context, id uint64, u *domain.User) error {
	q := `UPDATE users SET phone = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, q, nullStrPtr(u.Phone), id)
	if err != nil {
		return apperrors.InternalServerError("gagal update nomor HP")
	}
	if ra, _ := res.RowsAffected(); ra == 0 {
		return apperrors.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) UpdatePersonalizationConsent(ctx context.Context, id uint64, enabled bool) error {
	q := `
		UPDATE users
		SET is_personalization_enabled = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, q, enabled, id)

	if err != nil {
		return apperrors.InternalServerError("gagal update consent")
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) SoftDelete(ctx context.Context, id uint64) error {
	q := `
		UPDATE users
		SET deleted_at = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, q, time.Now(), id)

	if err != nil {
		return apperrors.InternalServerError("gagal menghapus user")
	}

	if ra, _ := res.RowsAffected(); ra == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
