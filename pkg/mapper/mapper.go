package mapper

import (
	"github.com/GantangSatria/MyBank-BE/internal/domain"
	"github.com/GantangSatria/MyBank-BE/pkg/dto/response"
)

func MapUserToResponse(u *domain.User) response.UserResponse {

	return response.UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Phone:    u.Phone,
		Segment:  u.Segment,
		IsActive: u.IsActive,
	}
}