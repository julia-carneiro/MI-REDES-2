package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ReqRotas struct {
	Origem  string `json:"Origem"`
	Destino string `json:"Destino"`
}

var server1 = "http://server1:"
var server2 = "http://server2:"
var server3 = "http://server3:"

func Teste2() {
	rota := ReqRotas{
		Origem:  "Feira",
		Destino: "São Paulo",
	}

	jsonData, err := json.Marshal(rota)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}

	// Envia para server1
	resp, err := http.Post(server1+"8000/rota/MenorCaminho", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erro ao enviar request para server1:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta do server1:", err)
		return
	}

	var menoresCaminhos1 [][]Trecho
	err = json.Unmarshal(body, &menoresCaminhos1)
	if err != nil {
		fmt.Println("Erro ao converter o JSON de server1:", err)
		return
	}

	// Envia para server2
	resp, err = http.Post(server2+"8001/rota/MenorCaminho", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erro ao enviar request para server2:", err)
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta do server2:", err)
		return
	}

	var menoresCaminhos2 [][]Trecho
	err = json.Unmarshal(body, &menoresCaminhos2)
	if err != nil {
		fmt.Println("Erro ao converter o JSON de server2:", err)
		return
	}

	// Envia para server3
	resp, err = http.Post(server3+"8002/rota/MenorCaminho", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Erro ao enviar request para server3:", err)
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta do server3:", err)
		return
	}

	var menoresCaminhos3 [][]Trecho
	err = json.Unmarshal(body, &menoresCaminhos3)
	if err != nil {
		fmt.Println("Erro ao converter o JSON de server3:", err)
		return
	}

	fmt.Printf("\nServer 1: %v\n\n", menoresCaminhos1[0])
	fmt.Printf("\nServer 2: %v\n\n", menoresCaminhos2[0])
	fmt.Printf("\nServer 3: %v\n\n", menoresCaminhos3[0])
}
