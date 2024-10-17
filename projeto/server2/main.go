package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Trecho struct {
	Origem  string `json:"Origem"`
	Destino string `json:"Destino"`
	Vagas   int    `json:"Vagas"`
	Peso    int    `json:"Peso"`
	Comp    string `json:"Comp"`
	ID      int    `json: "ID"`
}

func main() {

	response, err := http.Get("http://localhost:8000/rota")
	if err != nil {
		fmt.Println("Error calling Server 1:", err)
		return
	}
	defer response.Body.Close()

	var rota Trecho
	if err := json.NewDecoder(response.Body).Decode(&rota); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Printf("Recebido: %+v\n", rota)
}
