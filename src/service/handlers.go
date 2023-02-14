package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, fmt.Sprintf("/order?uid=%s", r.FormValue("uid")), http.StatusSeeOther)
	}
	tmpl, err := template.ParseFiles("./ui/html/home.html")
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
		errHandler(w, "this order dosen't exist ", http.StatusNotFound)
		return
	}
	data, ok := cache[uid]
	if !ok {
		errHandler(w, "this order dosen't exist ", http.StatusNotFound)
		return
	}

	//tmpl, err := template.ParseFiles("./ui/html/data.html")
	//if err != nil {
	//	errHandler(w, "Can't parse html", http.StatusInternalServerError)
	//	return
	//}
	//err = tmpl.Execute(w, prettyJSON)

	w.WriteHeader(http.StatusOK)
}
