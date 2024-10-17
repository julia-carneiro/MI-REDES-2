package funcoes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type Trecho struct {
	Origem  string `json:"Origem"`
	Destino string `json:"Destino"`
	Vagas   int    `json:"Vagas"`
	Peso    int    `json:"Peso"`
	Comp    string `json:"Comp"`
	ID      string `json:"ID"`
}
type Pessoa struct {
	Nome      string
	Sobrenome string
	Cpf       string
}
type Compra struct {
	Pessoa  Pessoa
	Trechos []Trecho
}

var rotas map[string][]Trecho
var filePathRotas = "dados/rotas.json" //caminho para arquivo de rotas

var locks = make(map[string]*sync.Mutex)
var lockMutex sync.Mutex // Protege o mapa de locks

func obterLock(id string) *sync.Mutex {
	lockMutex.Lock()
	defer lockMutex.Unlock()

	if _, ok := locks[id]; !ok {
		locks[id] = &sync.Mutex{}
	}
	return locks[id]
}

func EnviarTrechosServidor(trechos []Trecho, pessoa Pessoa) {
	var compra = Compra{
		Pessoa:  pessoa,
		Trechos: trechos,
	}
	// Converte o objeto para JSON
	jsonData, err := json.Marshal(compra)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}
	if trechos[0].Comp == "B" {
		//envia requisição para servidor B

		resp, err := http.Post("http://localhost:8001/compras", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(resp.Status)

	} else {
		// envia requisição para servidor C
		resp, err := http.Post("http://localhost:8002/compras", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(resp.Status)
	}

}

func ProcessarCompra(trechos []Trecho) {
	// Travando trechos de A
	var mutexes []*sync.Mutex
	for _, trecho := range trechos {
		lock := obterLock(trecho.ID)
		lock.Lock()
		mutexes = append(mutexes, lock)
	}

	/*
		LOGICA DA COMPRA AQUI
	*/

	// Libera todos os mutexes
	for _, m := range mutexes {
		m.Unlock()
	}

}

// Pega todas as rotas do arquivo json
func GetRotas(w http.ResponseWriter, r *http.Request) {
	rotas = LerRotas()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rotas)
}

// Função para compra de uma passagem
func Comprar(w http.ResponseWriter, r *http.Request) {
	// Lê o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	// Decodifica o JSON no objeto Compra
	var compra Compra
	err = json.Unmarshal(body, &compra)
	if err != nil {
		http.Error(w, "Erro ao decodificar o JSON", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	var trechos []Trecho = compra.Trechos
	var trechosA []Trecho
	var trechosB []Trecho
	var trechosC []Trecho

	for _, trecho := range trechos {

		if trecho.Comp == "A" {
			trechosA = append(trechosA, trecho)
		} else if trecho.Comp == "B" {
			trechosA = append(trechosB, trecho)
		} else {
			trechosA = append(trechosC, trecho)
		}

	}

	// mandando requisições para os outros servidores
	if len(trechosB) != 0 {
		// manda requisição para servidor B

		go EnviarTrechosServidor(trechosB, compra.Pessoa)
	}
	if len(trechosC) != 0 {
		// manda requisição para servidor B
		go EnviarTrechosServidor(trechosC, compra.Pessoa)
	}

	go ProcessarCompra(trechosA)

	wg.Wait()

}

func VerCompras(w http.ResponseWriter, r *http.Request) {

}

func LerRotas() map[string][]Trecho {
	// Abra o arquivo
	file, err := os.Open(filePathRotas)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return nil
	}
	defer file.Close()

	// Criar um mapa para armazenar as rotas
	var rotas map[string][]Trecho

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
