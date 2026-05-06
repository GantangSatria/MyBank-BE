package response

type UserResponse struct {
	ID                       uint64  `json:"id"`
	Name                     *string  `json:"name"`
	Email                    string  `json:"email"`
	Phone                    *string  `json:"phone"`
	DateOfBirth              *string `json:"date_of_birth,omitempty"`
	Occupation               string  `json:"occupation,omitempty"`
	Segment                  *string  `json:"segment,omitempty"`
	IsPersonalizationEnabled bool    `json:"is_personalization_enabled"`
	IsActive                 bool    `json:"is_active"`
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
	IsActive      bool    `json:"is_active"`
}