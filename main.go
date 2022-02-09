package main

import (
	"log"
	"net/http"
	"sqlorders"
)

func main() {

	go Stan()
	orderRepo := sqlorders.New()
	orderRepo.FromDb()
	reg := Reg{orderRepo}
	TakeMessage("test", "test", "test-1", &reg)
	mux := http.NewServeMux()
	mux.HandleFunc("/Order", reg.OrderReg)

	log.Println("Запуск веб-сервера на http://127.0.0.1:7777")
	err := http.ListenAndServe(":7777", mux)
	log.Fatal(err)
}
