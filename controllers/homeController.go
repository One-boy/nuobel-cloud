package controllers

import (
	"html/template"
	"net/http"
)

func OnHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		//认证通过
		temp, err := template.ParseFiles("views/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		temp.Execute(w, nil)
	}
}
