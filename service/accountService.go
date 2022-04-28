package service

import (
	"fmt"
	"learnings/banking/domain"
	"learnings/banking/dto"
	"learnings/banking/errs"
	"time"
)

type AccountService interface {
	NewAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError)
}

type DefaultAccountService struct {
	repo domain.AccountRespository
}

func (s DefaultAccountService) NewAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {
	valErr := req.Validate()
	if valErr != nil {
		return nil, valErr
	}

	requestAccount := domain.Account{
		AccountId:   "",
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format("2006-01-02 15:04:05"),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}
	fmt.Println(requestAccount)
	newAccount, err := s.repo.Save(requestAccount)
	if err != nil {
		return nil, err
	}
	na := newAccount.ToNewAccountResponseDto()
	return &na, nil
}

func NewAccountService(repo domain.AccountRespository) DefaultAccountService {
	return DefaultAccountService{repo: repo}
}
