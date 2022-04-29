package domain

import (
	"database/sql"
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

func (d AccountRespositoryDb) ById(id string) (*Account, *errs.AppError) {
	accountSql := "select customer_id, opening_date, account_type, amount, status from  accounts where account_id = ?"

	var a Account
	err := d.client.Get(&a, accountSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("customer not found")
		} else {
			logger.Error("Error while  querying " + err.Error())

			return nil, errs.NewUnexpectedError("unexpected database error")
		}

	}
	return &a, nil
}

func (d AccountRespositoryDb) MakeTransaction(t Transaction) (*Transaction, *errs.AppError) {
	tx, err := d.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for bank account transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// inserting bank account transaction
	result, _ := tx.Exec(`INSERT INTO transactions (account_id, amount, transaction_type, transaction_date) 
											values (?, ?, ?, ?)`, t.AccountId, t.Amount, t.TransactionType, t.TransactionDate)

	// updating account balance
	if t.IsWithdrawal() {
		_, err = tx.Exec(`UPDATE accounts SET amount = amount - ? where account_id = ?`, t.Amount, t.AccountId)
	} else {
		_, err = tx.Exec(`UPDATE accounts SET amount = amount + ? where account_id = ?`, t.Amount, t.AccountId)
	}

	// in case of error Rollback, and changes from both the tables will be reverted
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	// commit the transaction when all is good
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error while commiting transaction for bank account: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	// getting the last transaction ID from the transaction table
	transactionId, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting the last transaction id: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// Getting the latest account information from the accounts table
	account, appErr := d.ById(t.AccountId)
	if appErr != nil {
		return nil, appErr
	}
	t.TransactionId = strconv.FormatInt(transactionId, 10)

	// updating the transaction struct with the latest balance
	t.Amount = account.Amount
	return &t, nil
}
