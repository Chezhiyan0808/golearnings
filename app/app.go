package app

import (
	"learnings/banking/domain"
	authdomain "learnings/banking/domain/auth"
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
	authRepository := authdomain.NewAuthRepository(dbClient)
	ch := CustomerHandlers{service: service.NewCustomerService(customerRepositoryDb)}
	ah := AccountHandler{service: service.NewAccountService(accountRepositoryDb)}
	auh := AuthHandler{service.NewLoginService(authRepository, authdomain.GetRolePermissions())}

	am := AuthMiddleware{service.NewLoginService(authRepository, authdomain.GetRolePermissions())}
	/*PRIVATE ROUTES*/
	privateRouter := router.PathPrefix("/api/v1").Subrouter()
	privateRouter.Use(am.authorizationHandler())
	privateRouter.HandleFunc("/customers", ch.getCustomers).Methods(http.MethodGet).Name("GetAllCustomers")
	privateRouter.HandleFunc("/customers/{customer_id}", ch.getCustomer).Methods(http.MethodGet).Name("GetCustomer")
	privateRouter.HandleFunc("/customers/{customer_id}/account", ah.CreateAccount).Methods(http.MethodPost).Name("NewAccount")
	privateRouter.HandleFunc("/customers/{customer_id}/account/{account_id}", ah.MakeTransaction).Methods(http.MethodPost).Name("NewTransaction")

	/* PUBLIC ROUTES */
	router.HandleFunc("/auth/login", auh.Login).Methods(http.MethodPost)
	router.HandleFunc("/auth/register", auh.NotImplementedHandler).Methods(http.MethodPost)

	router.HandleFunc("/refresh", auh.Refresh).Methods(http.MethodPost)
	router.HandleFunc("/verify", auh.Verify).Methods(http.MethodGet)

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
