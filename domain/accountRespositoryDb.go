package domain

import (
	"learnings/banking/errs"
	"learnings/banking/logger"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type AccountRespositoryDb struct {
	client *sqlx.DB
}

func (d AccountRespositoryDb) Save(a Account) (*Account, *errs.AppError) {
	insertSQL := "INSERT INTO accounts (customer_id, opening_date, account_type, amount, status) values (?, ?, ?, ?, ?)"
	result, err := d.client.Exec(insertSQL, a.CustomerId, a.OpeningDate, a.AccountType, a.Amount, a.Status)
	if err != nil {
		logger.Error("Error while creating account " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpexcted DB error")
	}
	accountId, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting the last created  account id" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpexcted DB error")
	}
	a.AccountId = strconv.FormatInt(accountId, 10)

	return &a, nil
}

func NewAccountRepositoryDb(dbClient *sqlx.DB) AccountRespositoryDb {
	return AccountRespositoryDb{dbClient}
}
