package funcoes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Rota struct {
	Destino string `json:"Destino"`
	Vagas   int    `json:"Vagas"`
	Peso    int    `json:"Peso"`
	Comp    string `json:"Comp"`
}
type Pessoa struct {
	Nome      string
	Sobrenome string
	Cpf       string
}
type Compra struct {
	Origem  string
	Destino string
	Comp    string
	Pessoa  Pessoa
}

var rotas map[string][]Rota
var filePathRotas = "/app/servidor/dados/rotas.json" //caminho para arquivo de rotas

func GetRotas(w http.ResponseWriter, r *http.Request) {
	rotas = LerRotas()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rotas)
}

func Comprar(w http.ResponseWriter, r *http.Request) {

}

func VerCompras(w http.ResponseWriter, r *http.Request) {

}

func LerRotas() map[string][]Rota {
	// Abra o arquivo
	file, err := os.Open(filePathRotas)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return nil
	}
	defer file.Close()

	// Criar um mapa para armazenar as rotas
	var rotas map[string][]Rota

	// Decodificar o arquivo JSON para o mapa
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&rotas); err != nil {
		fmt.Println("Erro ao decodificar o JSON:", err)
		return nil
	}
	return rotas
}

func LerCompras(Cpf string) []Compra {
	file, err := os.Open(filePathRotas)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return nil
	}
	defer file.Close()

	// Criar um mapa para armazenar as rotas
	var todas_compras map[string][]Compra

	// Decodificar o arquivo JSON para o mapa
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&todas_compras); err != nil {
		fmt.Println("Erro ao decodificar o JSON:", err)
		return nil
	}

	var compras = todas_compras[Cpf]
	return compras
}
