package request

import "time"

type UpdateProfileRequest struct {
	Name        *string `json:"name"          validate:"omitempty,min=2,max=100"`
	Occupation  *string `json:"occupation"    validate:"omitempty,max=100"`
	DateOfBirth *string `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02"`
}

// ParsedDateOfBirth mengonversi string "YYYY-MM-DD" ke *time.Time
func (r *UpdateProfileRequest) ParsedDateOfBirth() *time.Time {
	if r.DateOfBirth == nil {
		return nil
	}
	t, err := time.Parse("2006-01-02", *r.DateOfBirth)
	if err != nil {
		return nil
	}
	return &t
}

type UpdatePhoneRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

type UpdatePersonalizationRequest struct {
	Enabled bool `json:"enabled"`
}