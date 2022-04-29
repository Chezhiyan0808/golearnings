package app

import (
	"learnings/banking/domain"
	"learnings/banking/service"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func Start() {
	// mux := http.NewServeMux()
	router := mux.NewRouter()

	dbClient := getDbClient()
	customerRepositoryDb := domain.NewCustomerRepositoryDb(dbClient)
	accountRepositoryDb := domain.NewAccountRepositoryDb(dbClient)
	ch := CustomerHandlers{service: service.NewCustomerService(customerRepositoryDb)}
	ah := AccountHandler{service: service.NewAccountService(accountRepositoryDb)}

	router.HandleFunc("/customers", ch.getCustomers).Methods(http.MethodGet)

	router.HandleFunc("/customers/{customer_id}", ch.getCustomer).Methods(http.MethodGet)

	router.HandleFunc("/customers/{customer_id}/account", ah.CreateAccount).Methods(http.MethodPost)

	router.HandleFunc("/customers/{customer_id}/account/{account_id}", ah.MakeTransaction).Methods(http.MethodPost)

	router.HandleFunc("/greet", greet)

	router.HandleFunc("/api/time", getTime)

	log.Fatal(http.ListenAndServe("localhost:8000", router))
}

func getDbClient() *sqlx.DB {
	client, err := sqlx.Open("mysql", "muthu:password@tcp(localhost:3306)/banking")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client
}
