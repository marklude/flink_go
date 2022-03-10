package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/marklude/flink_go/logger"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

var history = make(map[string][]Location)

type jsonHistory struct {
	OrderId string     `json:"order_id"`
	History []Location `json:"history"`
}

func PostLocation(w http.ResponseWriter, r *http.Request) {
	var l Location
	vars := mux.Vars(r)
	// Get orderId
	orderId, ok := vars["order_id"]
	if !ok {
		logger.WarnMessage("missing orderId in request")
	}
	// Unmarshal request
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve existing history for order Id
	oldLoc := history[orderId]

	newLoc := append(oldLoc, l)
	history[orderId] = newLoc

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order placed with ID: %s", orderId)

}

func GetLocation(w http.ResponseWriter, r *http.Request) {
	// Get order id
	vars := mux.Vars(r)
	orderId, ok := vars["order_id"]

	if !ok {
		logger.WarnMessage("missing order id")
		http.Error(w, "missing order id", http.StatusNotFound)
		return
	}

	max := r.URL.Query().Get("max")
	var nHistory []Location
	if max != "" {

		nMax, err := strconv.ParseInt(max, 0, 64)
		if err != nil {
			logger.ErrorMessage("parsing max failed", err)
			http.Error(w, "max should be a number", http.StatusBadRequest)
			return

		}

		if nMax > int64(len(history[orderId])) {
			nHistory = history[orderId]
		} else {
			nHistory = history[orderId][:nMax]
		}

	} else {
		nHistory = history[orderId]
	}

	history, err := json.Marshal(jsonHistory{OrderId: orderId, History: nHistory})

	if err != nil {
		logger.ErrorMessage("cannot marshal history", err)
	}

	// Return http status and message
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(history))
}

func DeleteLocation(w http.ResponseWriter, r *http.Request) {
	// Get the order id
	vars := mux.Vars(r)
	orderId, ok := vars["order_id"]

	// Not okay
	if !ok {
		logger.WarnMessage("missing order id")
		http.Error(w, "missing order id", http.StatusBadRequest)
		return
	}

	// Set history to empty array
	history[orderId] = []Location{}

	// Return http status and message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted order with ID: %s", orderId)
}
