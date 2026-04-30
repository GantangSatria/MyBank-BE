package request

// RegisterRequest
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=100"`

	Name        string `json:"name" validate:"omitempty,min=2,max=100"`
	Phone       string `json:"phone" validate:"omitempty,min=10,max=15"`
	DateOfBirth string `json:"date_of_birth" validate:"omitempty"`
	Occupation  string `json:"occupation" validate:"omitempty,max=100"`
}

// LoginRequest
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// SetPINRequest (BARU → setelah register)
type SetPINRequest struct {
	PIN string `json:"pin" validate:"required,len=6,numeric"`
}

// ChangePINRequest
type ChangePINRequest struct {
	OldPIN string `json:"old_pin" validate:"required,len=6,numeric"`
	NewPIN string `json:"new_pin" validate:"required,len=6,numeric"`
}

// ChangePasswordRequest
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// RefreshTokenRequest
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}