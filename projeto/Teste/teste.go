package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// Estruturas definidas para a compra
type Trecho struct {
	Origem  string `json:"Origem"`
	Destino string `json:"Destino"`
	Vagas   int    `json:"Vagas"`
	Peso    int    `json:"Peso"`
	Comp    string `json:"Comp"`
	ID      string `json:"ID"`
}

type Pessoa struct {
	Nome      string `json:"Nome"`
	Sobrenome string `json:"Sobrenome"`
	Cpf       string `json:"Cpf"`
}

type Compra struct {
	Pessoa        Pessoa   `json:"Pessoa"`
	Trechos       []Trecho `json:"Trechos"`
	Participantes []string `json:"Participantes"`
}

type RetornoCompra struct {
	Resultado bool   `json:"Resultado"`
	Server    string `json:"Server"`
	Compra    Compra `json:"Compra"`
}

var mutex sync.Mutex
var InfoCompras = []RetornoCompra{}
var filePath = "compras.json"

func SalvarInfo(dados RetornoCompra) {
	mutex.Lock()         // Adquire o bloqueio
	defer mutex.Unlock() // Garante que o bloqueio será liberado
	InfoCompras = append(InfoCompras, dados)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Erro ao escrever:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(InfoCompras)
}

// Função para enviar requisições POST para um servidor específico
func enviarRequisicao(wg *sync.WaitGroup, url string, startSignal <-chan struct{}, compra Compra) {
	defer wg.Done()

	// Aguarda o sinal para começar
	<-startSignal

	// Converte a estrutura Compra para JSON
	payload, err := json.Marshal(compra)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}

	// Envia a requisição POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Erro ao enviar requisição:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return
	}

	var dadosResposta RetornoCompra //resposta da preparação. Se conseguiu praparar ou não
	err = json.Unmarshal(body, &dadosResposta)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return
	}

	fmt.Printf("\nResposta do servidor %s:\n", url)
	fmt.Printf("Resultado: %v\n", dadosResposta.Resultado)
	fmt.Printf("Servidor de Resposta: %s\n", dadosResposta.Server)
	fmt.Printf("Pessoa: %s %s (CPF: %s)\n", 
			dadosResposta.Compra.Pessoa.Nome, 
			dadosResposta.Compra.Pessoa.Sobrenome, 
			dadosResposta.Compra.Pessoa.Cpf)

	fmt.Println("Trechos:")
	for i, trecho := range dadosResposta.Compra.Trechos {
		fmt.Printf("  Trecho %d: %s -> %s | Peso: %d | Companhia: %s\n",
				i+1, trecho.Origem, trecho.Destino, trecho.Peso, trecho.Comp)
	}

	fmt.Printf("Participantes: %v\n\n", dadosResposta.Compra.Participantes)
	SalvarInfo(dadosResposta)
	fmt.Println("Status da respota do servidor:", resp.Status)
}

func main() {
	// URLs dos servidores
	servidores := []string{
		"http://server1:8000/compras",
		"http://server2:8001/compras",
		"http://server3:8002/compras",
	}

	// // Número total de requisições desejadas
	// numeroTotalDeRequisicoes := 30
	// // Calcula o número de requisições por servidor
	// // numeroDeRequisicoesPorServidor := numeroTotalDeRequisicoes / len(servidores)
	numero := 1

	// Dados de exemplo para a compra
	compra := Compra{
		Pessoa: Pessoa{
			Nome:      "João",
			Sobrenome: "Silva",
			Cpf:       "12345678900",
		},
		Trechos: []Trecho{
			{
				Origem:  "Brasília",
				Destino: "Salvador",
				Vagas:   12,
				Peso:    5,
				Comp:    "A",
				ID:      "7",
			},
			{
				Origem:  "Belo Horizonte",
				Destino: "Rio de Janeiro",
				Vagas:   2,
				Peso:    20,
				Comp:    "B",
				ID:      "0",
			},
			{
				Origem:  "Fortaleza",
				Destino: "São Paulo",
				Vagas:   10,
				Peso:    5,
				Comp:    "C",
				ID:      "5",
			},
		},
		Participantes: []string{"A", "B", "C"},
	}

	// Usa sync.WaitGroup para aguardar todas as goroutines finalizarem
	var wg sync.WaitGroup

	// Canal para sincronizar o início das requisições
	startSignal := make(chan struct{})

	// Cria uma goroutine para cada requisição simultânea em cada servidor
	for _, url := range servidores {
		for i := 0; i < numero; i++ {
			wg.Add(1)
			go enviarRequisicao(&wg, url, startSignal, compra)
		}
	}

	// Espera para preparar todas as goroutines e então libera o startSignal
	time.Sleep(2 * time.Second) // Aguarda para garantir que todas as goroutines estão prontas
	close(startSignal)

	// Espera todas as goroutines completarem
	wg.Wait()
	Teste2()

	fmt.Println("Teste finalizado")
}
