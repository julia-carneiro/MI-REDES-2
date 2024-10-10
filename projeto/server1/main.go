package main

import (
	"log"
	"net/http"
	"projeto/server1/funcoes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/rota2", funcoes.GetRotas).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
