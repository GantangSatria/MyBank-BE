package request

import "time"

type UpdateProfileRequest struct {
	Name        *string `json:"name"          validate:"omitempty,min=2,max=100" example:"John Doe"`
	Occupation  *string `json:"occupation"    validate:"omitempty,max=100" example:"Software Engineer"`
	DateOfBirth *string `json:"date_of_birth" validate:"omitempty,datetime=2006-01-02" example:"1990-01-01"`
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
	Phone string `json:"phone" validate:"required,e164" example:"+6281234567890"`
}

type UpdatePersonalizationRequest struct {
	Enabled bool `json:"enabled" example:"true"`
}