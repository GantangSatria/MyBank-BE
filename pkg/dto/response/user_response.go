package response

type UserResponse struct {
	ID                       uint64  `json:"id"`
	Name                     *string `json:"name"`
	Email                    string  `json:"email"`
	Phone                    *string `json:"phone"`
	Gender                   *string `json:"gender,omitempty"`
	DateOfBirth              *string `json:"date_of_birth,omitempty"`
	Occupation               string  `json:"occupation,omitempty"`
	MaritalStatus            *string `json:"marital_status,omitempty"`
	Segment                  *string `json:"segment,omitempty"`
	MonthlyIncomeRange       *string `json:"monthly_income_range,omitempty"`
	IsPersonalizationEnabled bool    `json:"is_personalization_enabled"`
	IsActive                 bool    `json:"is_active"`
	LastLoginAt              *string `json:"last_login_at,omitempty"`
	CreatedAt                string  `json:"created_at"`
}

 
type UserWithAccountsResponse struct {
	UserResponse
	Accounts []AccountResponse `json:"accounts"`
}

type AccountResponse struct {
	ID            uint64  `json:"id"`
	AccountNumber string  `json:"account_number"`
	AccountType   string  `json:"account_type"`
	Balance       float64 `json:"balance"`
	Currency      string  `json:"currency"`
	Branch        string  `json:"branch,omitempty"`
	IsActive      bool    `json:"is_active"`
}

//tes