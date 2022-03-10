package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/marklude/flink_go/logger"
	"github.com/marklude/flink_go/redisDB"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

var (
	history = make(map[string][]Location)
	ctx     = context.Background()
)

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

	// Get history ttl from env
	historyTTL := os.Getenv("LOCATION_HISTORY_TTL_SECONDS")

	// Save to redis
	hisJson, _ := json.Marshal(history[orderId])
	if historyTTL != "" {
		ttl, err := strconv.ParseInt(historyTTL, 0, 64)
		if err != nil {
			logger.ErrorMessage("Parsing ttl failed", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = redisDB.Rds.Set(ctx, orderId, hisJson, time.Duration(ttl)*time.Second).Err()
		if err != nil {
			logger.ErrorMessage(fmt.Sprintf("Order:%s not created", orderId), err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err := redisDB.Rds.Set(ctx, orderId, hisJson, 0).Err()
		if err != nil {
			logger.ErrorMessage(fmt.Sprintf("Order:%s not created", orderId), err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order placed with ID: %s", orderId)

}

func GetLocation(w http.ResponseWriter, r *http.Request) {
	var locations []Location
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
	// Retrieve order from redis
	result, _ := redisDB.Rds.Get(ctx, orderId).Result()
	if result != "" {
		json.Unmarshal([]byte(result), &locations)

	}

	if max != "" {

		nMax, err := strconv.ParseInt(max, 0, 64)
		if err != nil {
			logger.ErrorMessage("parsing max failed", err)
			http.Error(w, "max should be a number", http.StatusBadRequest)
			return

		}

		if nMax > int64(len(history[orderId])) {
			nHistory = locations
		} else {
			nHistory = locations[:nMax]
		}

	} else {
		nHistory = locations
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

	// Remove order history from redis
	redisDB.Rds.Del(ctx, orderId)

	// Return http status and message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted order with ID: %s", orderId)
}
