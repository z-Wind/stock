package main

import (
	"encoding/json"
	"fmt"
	"github.com/z-Wind/stock/api/gotd/api"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

func savedOrder(w http.ResponseWriter, req *http.Request) {
	err := td.RefreshAccessTokenOrNot()
	if err != nil {
		http.Error(w, "Could not RefreshAccessToken", http.StatusInternalServerError)
		return
	}

	if len(accountIDs) == 0 || len(accountIDs) > 1 {
		http.Error(w, "No accountIDs", http.StatusInternalServerError)
		return
	}
	switch req.Method {
	case http.MethodGet:
		getSavedOrders(w, req)
	case http.MethodPost:
		createSavedOrder(w, req)
	case http.MethodDelete:
		deleteSavedOrder(w, req)
	}
}

func deleteSavedOrder(w http.ResponseWriter, req *http.Request) {
	savedOrderIDs, ok := req.URL.Query()["savedOrderID"]
	if !ok || len(savedOrderIDs) > 1 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	savedOrderID, err := strconv.ParseInt(savedOrderIDs[0], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("savedOrderID %s is not Int", savedOrderIDs[0]), http.StatusBadRequest)
		return
	}

	if err := td.DeleteSavedOrder(accountIDs[0], savedOrderID); err != nil {
		http.Error(w, errors.WithMessage(err, "td.DeleteSavedOrder").Error(), http.StatusBadRequest)
		return
	}
}

func getSavedOrders(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	orders, err := td.GetSavedOrdersbyPath(accountIDs[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(orders)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func createSavedOrder(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	var data savedOrderParas
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	o := api.NewOrder(data.Symbol, data.AssetType, data.Instruction, data.Price, data.Qunatity)
	err = td.CreateSavedOrder(accountIDs[0], o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
