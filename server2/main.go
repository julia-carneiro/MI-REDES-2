package main

import (
	"fmt"
	"log"
	"net/http"
	"projeto/funcoes"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Print("teste")
	for i := range funcoes.TrechoLivre {
		funcoes.TrechoLivre[i] = true
	}
	router := mux.NewRouter()
	router.HandleFunc("/rota", funcoes.GetRotas).Methods("GET")
	router.HandleFunc("/compras", funcoes.SolicitacaoCord).Methods("POST")           //Comprar
	router.HandleFunc("/compras/preparar", funcoes.Commit).Methods("POST")           //Preparar para commit
	router.HandleFunc("/compras/confirmar", funcoes.ConfirmarCommit).Methods("POST") //confirmar commit
	router.HandleFunc("/compras/cancelar", funcoes.CancelarCommit).Methods("POST")   //Cancelar commit
	//router.HandleFunc("/compras/{id}", funcoes.VerCompras).Methods("GET") //Ver compras
	//router.HandleFunc("/rota", funcoes.GetRotas).Methods("GET")
	log.Fatal(http.ListenAndServe(":8001", router))
}
