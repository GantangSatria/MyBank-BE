package response

type UserResponse struct {
	ID                       uint64  `json:"id" example:"1"`
	Name                     *string `json:"name" example:"John Doe"`
	Email                    string  `json:"email" example:"user@example.com"`
	Phone                    *string `json:"phone" example:"081234567890"`
	Gender                   *string `json:"gender,omitempty" example:"Male"`
	DateOfBirth              *string `json:"date_of_birth,omitempty" example:"1990-01-01"`
	Occupation               string  `json:"occupation,omitempty" example:"Software Engineer"`
	MaritalStatus            *string `json:"marital_status,omitempty" example:"Single"`
	Segment                  *string `json:"segment,omitempty" example:"Premium"`
	MonthlyIncomeRange       *string `json:"monthly_income_range,omitempty" example:"10M-20M"`
	IsPersonalizationEnabled bool    `json:"is_personalization_enabled" example:"true"`
	IsActive                 bool    `json:"is_active" example:"true"`
	LastLoginAt              *string `json:"last_login_at,omitempty" example:"2023-10-01 12:00:00"`
	CreatedAt                string  `json:"created_at" example:"2023-10-01 12:00:00"`
}

 
type UserWithAccountsResponse struct {
	UserResponse
	Accounts []AccountResponse `json:"accounts"`
}

//tes