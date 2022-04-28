package domain

type CustomerRepositoryStub struct {
	customers []Customer
}

func (s CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return s.customers, nil
}

func NewCustomerRepositoryStub() CustomerRepositoryStub {
	customers := []Customer{
		{Id: "0001", Name: "Muthu", City: "Chennai", Zipcode: "600100", DateOfBirth: "08/08/1990", Status: "1"},
		{Id: "0002", Name: "Muthu2", City: "Chennai", Zipcode: "600100", DateOfBirth: "08/08/1990", Status: "1"},
		{Id: "0003", Name: "Muthu3", City: "Chennai", Zipcode: "600100", DateOfBirth: "08/08/1990", Status: "1"},
		{Id: "0004", Name: "Muthu4", City: "Chennai", Zipcode: "600100", DateOfBirth: "08/08/1990", Status: "1"},
	}
	return CustomerRepositoryStub{customers: customers}
}
