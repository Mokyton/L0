package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errHandler(w, "Not Found ", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		http.Redirect(w, r, fmt.Sprintf("/order?uid=%s", r.FormValue("uid")), http.StatusSeeOther)
	}
	tmpl, err := template.ParseFiles("./ui/html/home.gohtml")
	if err != nil {
		errHandler(w, "Can't parse html", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		errHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func errHandler(w http.ResponseWriter, message string, errStatus int) {
	jsn := struct {
		Error string `json:"error"`
	}{message}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(jsn)
	w.WriteHeader(errStatus)
}

func uidHandler(w http.ResponseWriter, r *http.Request) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		errHandler(w, "this order doesn't exist ", http.StatusNotFound)
		return
	}
	data, ok := cache[uid]
	if !ok {
		errHandler(w, "this order doesn't exist ", http.StatusNotFound)
		return
	}

	Model := struct {
		OrderUID    string `json:"order_uid"`
		TrackNumber string `json:"track_number"`
		Entry       string `json:"entry"`
		Delivery    struct {
			Name    string `json:"name"`
			Phone   string `json:"phone"`
			Zip     string `json:"zip"`
			City    string `json:"city"`
			Address string `json:"address"`
			Region  string `json:"region"`
			Email   string `json:"email"`
		} `json:"delivery"`
		Payment struct {
			Transaction  string `json:"transaction"`
			RequestID    string `json:"request_id"`
			Currency     string `json:"currency"`
			Provider     string `json:"provider"`
			Amount       int    `json:"amount"`
			PaymentDt    int    `json:"payment_dt"`
			Bank         string `json:"bank"`
			DeliveryCost int    `json:"delivery_cost"`
			GoodsTotal   int    `json:"goods_total"`
			CustomFee    int    `json:"custom_fee"`
		} `json:"payment"`
		Items []struct {
			ChrtID      int    `json:"chrt_id"`
			TrackNumber string `json:"track_number"`
			Price       int    `json:"price"`
			Rid         string `json:"rid"`
			Name        string `json:"name"`
			Sale        int    `json:"sale"`
			Size        string `json:"size"`
			TotalPrice  int    `json:"total_price"`
			NmID        int    `json:"nm_id"`
			Brand       string `json:"brand"`
			Status      int    `json:"status"`
		} `json:"items"`
		Locale            string    `json:"locale"`
		InternalSignature string    `json:"internal_signature"`
		CustomerID        string    `json:"customer_id"`
		DeliveryService   string    `json:"delivery_service"`
		Shardkey          string    `json:"shardkey"`
		SmID              int       `json:"sm_id"`
		DateCreated       time.Time `json:"date_created"`
		OofShard          string    `json:"oof_shard"`
	}{}

	json.NewDecoder(strings.NewReader(string(data))).Decode(&Model)

	tmpl, err := template.ParseFiles("./ui/html/data.gohtml")
	if err != nil {
		errHandler(w, "Can't parse html", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, Model)

	w.WriteHeader(http.StatusOK)
}
