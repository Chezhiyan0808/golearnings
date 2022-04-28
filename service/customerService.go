package service

import (
	"learnings/banking/domain"
	"learnings/banking/dto"
	"learnings/banking/errs"
)

type CustomerService interface {
	GetAllCustomers(status string) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func (s DefaultCustomerService) GetAllCustomers(status string) ([]dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.FindAll(status)
	if err != nil {
		return nil, err
	}
	r := make([]dto.CustomerResponse, len(c))
	for i, v := range c {
		r[i] = v.ToDto()
	}
	return r, nil
}

func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}
	response := c.ToDto()
	return &response, nil
}

func NewCustomerService(respository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{respository}
}
