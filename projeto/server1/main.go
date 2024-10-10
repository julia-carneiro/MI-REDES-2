package main

import (
	"log"
	"net/http"
	"projeto/funcoes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/rota", funcoes.GetRotas).Methods("GET")
	router.HandleFunc("/compras", funcoes.Comprar).Methods("POST")        //Comprar
	router.HandleFunc("/compras/{id}", funcoes.VerCompras).Methods("GET") //Ver compras
	router.HandleFunc("/rota", funcoes.GetRotas).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
