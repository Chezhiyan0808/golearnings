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
