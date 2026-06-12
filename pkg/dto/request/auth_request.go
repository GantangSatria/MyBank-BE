package request

// RegisterRequest
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=100" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"secret123"`

	Name        string `json:"name" validate:"omitempty,min=2,max=100" example:"John Doe"`
	Phone       string `json:"phone" validate:"omitempty,min=10,max=15" example:"081234567890"`
	DateOfBirth string `json:"date_of_birth" validate:"omitempty" example:"1990-01-01"`
	Occupation  string `json:"occupation" validate:"omitempty,max=100" example:"Software Engineer"`
}

// LoginRequest
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"secret123"`
}

// SetPINRequest (BARU → setelah register)
type SetPINRequest struct {
	PIN string `json:"pin" validate:"required,len=6,numeric" example:"123456"`
}

// ChangePINRequest
type ChangePINRequest struct {
	OldPIN string `json:"old_pin" validate:"required,len=6,numeric" example:"123456"`
	NewPIN string `json:"new_pin" validate:"required,len=6,numeric" example:"654321"`
}

// ChangePasswordRequest
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required" example:"secret123"`
	NewPassword string `json:"new_password" validate:"required,min=6" example:"newsecret123"`
}

// RefreshTokenRequest
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}