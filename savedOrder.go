package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/z-Wind/gotd"
	"github.com/z-Wind/stock/stocker"
)

func savedOrder(w http.ResponseWriter, req *http.Request) {
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
	savedOrderIDQ := req.URL.Query().Get("savedOrderID")
	if savedOrderIDQ == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	savedOrderID, err := strconv.ParseInt(savedOrderIDQ, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("savedOrderID %s is not Int", savedOrderIDQ), http.StatusBadRequest)
		return
	}

	td := stockers[0].(*stocker.TDAmeritrade).Service
	if _, err := td.SavedOrders.DeleteSavedOrder(accountID, savedOrderID).Do(); err != nil {
		http.Error(w, errors.WithMessage(err, "td.SavedOrders.DeleteSavedOrder").Error(), http.StatusBadRequest)
		return
	}
}

func getSavedOrders(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	td := stockers[0].(*stocker.TDAmeritrade).Service
	orders, err := td.SavedOrders.GetSavedOrdersByPath(accountID).Do()
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

	var data struct {
		Symbol      string  `json:"Symbol"`
		AssetType   string  `json:"AssetType"`
		Instruction string  `json:"Instruction"`
		Price       float64 `json:"Price"`
		Quantity    float64 `json:"Quantity"`
	}

	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	savedOrder := &gotd.SavedOrder{
		Order: &gotd.Order{
			Session:    "NORMAL",
			Duration:   "GOOD_TILL_CANCEL",
			OrderType:  "LIMIT",
			CancelTime: time.Now().AddDate(0, 4, 0).UTC().Format("2006-01-02"),
			Price:      data.Price,
			OrderLegCollections: []*gotd.OrderLegCollection{
				&gotd.OrderLegCollection{
					Instrument: &gotd.Instrument{
						Symbol:    data.Symbol,
						AssetType: data.AssetType,
					},
					Instruction: data.Instruction,
					Quantity:    data.Quantity,
				},
			},
		},
	}

	td := stockers[0].(*stocker.TDAmeritrade).Service
	_, err = td.SavedOrders.CreateSavedOrder(accountID, savedOrder).Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
