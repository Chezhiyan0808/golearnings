package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"learnings/banking/service"

	"github.com/gorilla/mux"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

type CustomerHandlers struct {
	service service.CustomerService
}

func (ch *CustomerHandlers) getCustomers(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query().Get("status")
	status := "ALL"
	switch strings.ToLower(vars) {
	case "active":
		status = "1"
	case "inactive":
		status = "0"
	default:
		status = "ALL"
	}

	customers, err := ch.service.GetAllCustomers(status)

	if err != nil {
		writeResponse(w, err.Code, err.AssMessage())
	} else {
		writeResponse(w, http.StatusOK, customers)
	}

	// if err != nil {
	// 	return nil, err
	// }
	// if r.Header.Get("Content-Type") == "application/json" {
	// 	w.Header().Add("Content-Type", "applcation/json")
	// 	json.NewEncoder(w).Encode(customers)
	// }
	// w.Header().Add("Content-Type", "applcation/xml")

	// xml.NewEncoder(w).Encode(customers)
}

func (ch *CustomerHandlers) getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customer, err := ch.service.GetCustomer(vars["customer_id"])
	if err != nil {
		writeResponse(w, err.Code, err.AssMessage())
	} else {
		writeResponse(w, http.StatusOK, customer)
	}
}

func getTime(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query().Get("tz")
	response := make(map[string]string)
	if vars == "" {
		response["current_time"] = time.Now().String()
	} else {
		tzs := strings.Split(vars, ",")
		for _, tz := range tzs {
			loc, err := time.LoadLocation(strings.Trim(tz, " "))
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(fmt.Sprintf("invalid timezone %s in input", tz)))
				return
			}
			response[tz] = time.Now().In(loc).String()
		}
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}

}
