package main

import (
	"html/template"
	"log"
	"net/http"
	"sqlorders"
)

type Reg struct {
	order sqlorders.Repository
}

func (reg *Reg) OrderReg(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	order, _ := reg.order.FindById(id)
	tmpl, _ := template.ParseFiles("ui/static/order.html")
	err := tmpl.Execute(w, order)
	if err != nil {
		log.Fatalf("Error with template: %v", err)
	}

}
