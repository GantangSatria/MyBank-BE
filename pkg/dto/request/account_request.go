package request

type CreateAccountRequest struct {
	AccountType string `json:"account_type" validate:"required,oneof=SAVING CHECKING" example:"SAVING"`
}
