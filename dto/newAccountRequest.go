package dto

import (
	"learnings/banking/errs"
	"strings"
)

type NewAccountRequest struct {
	CustomerId  string  `json:"customer_id"`
	AccountType string  `json:"account_type"`
	Amount      float64 `json:"amount"`
}

func (r NewAccountRequest) Validate() *errs.AppError {
	if r.Amount < 5000 {
		return errs.NewValidationError("Initial Deposit should be a min on 5000")
	}
	if strings.ToLower(r.AccountType) != "checking" && strings.ToLower(r.AccountType) != "saving" {
		return errs.NewValidationError("Currently only checking and saving accounts can be processed")
	}
	return nil
}
