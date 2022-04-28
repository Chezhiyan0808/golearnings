package domain

import (
	"database/sql"
	"learnings/banking/errs"
	"learnings/banking/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {

	findAllSql := "select customer_id, name, city, zipcode, date_of_birth, status from  customers"
	if status != "ALL" {
		findAllSql = findAllSql + " where status = " + status
	}

	logger.Debug(findAllSql)

	customers := make([]Customer, 0)

	err := d.client.Select(&customers, findAllSql)
	if err != nil {
		logger.Error("Error while  querying " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	// rows, err := d.client.Query(findAllSql)

	// err = sqlx.StructScan(rows, &customers)
	// if err != nil {
	// 	logger.Error("Error while  querying " + err.Error())
	// 	return nil, errs.NewUnexpectedError("unexpected database error")
	// }
	// for rows.Next() {
	// 	var c Customer
	// 	err := rows.Scan(&c.Id, &c.Name, &c.City, &c.Zipcode, &c.DateOfBirth, &c.Status)
	// 	if err != nil {
	// 		logger.Error("Error while  querying " + err.Error())
	// 		return nil, errs.NewUnexpectedError("unexpected database error")
	// 	}
	// 	customers = append(customers, c)

	// }
	return customers, nil
}

func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	customerSql := "select customer_id, name, city, zipcode, date_of_birth, status from  customers where customer_id = ?"

	var c Customer
	err := d.client.Get(&c, customerSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("customer not found")
		} else {
			logger.Error("Error while  querying " + err.Error())

			return nil, errs.NewUnexpectedError("unexpected database error")
		}

	}
	return &c, nil
}

func NewCustomerRepositoryDb(dbClient *sqlx.DB) CustomerRepositoryDb {

	return CustomerRepositoryDb{client: dbClient}
}
