package app

import (
	"encoding/json"
	"learnings/banking/dto"
	"learnings/banking/service"
	"net/http"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	service service.AccountService
}

func (ah *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	cId := mux.Vars(r)["customer_id"]
	var request dto.NewAccountRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		request.CustomerId = cId
		account, err := ah.service.NewAccount(request)
		if err != nil {
			writeResponse(w, err.Code, err.AssMessage())
		} else {
			writeResponse(w, http.StatusOK, account)
		}
	}

}

func (h AccountHandler) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	// get the account_id and customer_id from the URL
	vars := mux.Vars(r)
	accountId := vars["account_id"]
	customerId := vars["customer_id"]

	// decode incoming request
	var request dto.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {

		//build the request object
		request.AccountId = accountId
		request.CustomerId = customerId

		// make transaction
		account, appError := h.service.MakeTransaction(request)

		if appError != nil {
			writeResponse(w, appError.Code, appError.AssMessage())
		} else {
			writeResponse(w, http.StatusOK, account)
		}
	}

}
