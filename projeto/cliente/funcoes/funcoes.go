package funcoes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Pessoa struct {
	Nome      string
	Sobrenome string
	Cpf       string
}
type Compra struct {
	Pessoa        Pessoa
	Trechos       []Trecho
	Participantes []string
}
type Trecho struct {
	Origem  string `json:"Origem"`
	Destino string `json:"Destino"`
	Vagas   int    `json:"Vagas"`
	Peso    int    `json:"Peso"`
	Comp    string `json:"Comp"`
	ID      string `json:"ID"`
}

var server1 = "http://server1:"
var server2 = "http://server2:"
var server3 = "http://server3:"

func BuscarRotaServidor(servidor string) map[string][]Trecho {

	// Inicializa o mapa
	trechos := make(map[string][]Trecho)

	var resp *http.Response
	var err error

	// Condicional para verificar o servidor
	if servidor == "A" {
		// BUSCA NO SERVIDOR 1
		//resp, err = http.Get("http://server1:8000/rota")
		resp, err = http.Get(server1 + "8000/rota")
	} else if servidor == "B" {
		// BUSCA NO SERVIDOR 2
		//resp, err = http.Get("http://server2:8001/rota")
		resp, err = http.Get(server2 + "8001/rota")
	} else if servidor == "C" {
		// BUSCA NO SERVIDOR 3
		//resp, err = http.Get("http://server3:8002/rota")
		resp, err = http.Get(server3 + "8002/rota")
	} else {
		fmt.Println("Servidor desconhecido:", servidor)
		return nil
	}

	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
		return nil
	}
	defer resp.Body.Close()

	// Lendo o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return nil
	}
	//fmt.Println(body)
	// Decodificando o JSON para o mapa
	err = json.Unmarshal(body, &trechos)
	if err != nil {
		fmt.Println("Erro ao converter o JSON:", err)
		return nil
	}

	return trechos
}

func GetRotas() map[string][]Trecho {
	// Inicializa o mapa
	rotas := make(map[string][]Trecho)

	trechosA := BuscarRotaServidor("A")
	trechosB := BuscarRotaServidor("B")
	trechosC := BuscarRotaServidor("C")

	// concatena trechosA com rotas
	for chave, valor := range trechosA {
		rotas[chave] = append(rotas[chave], valor...)
	}

	// concatena trechosB com rotas
	for chave, valor := range trechosB {
		rotas[chave] = append(rotas[chave], valor...)
	}

	// concatena trechosC com rotas
	for chave, valor := range trechosC {
		rotas[chave] = append(rotas[chave], valor...)
	}
	//fmt.Println(rotas)
	return rotas
}

func SolicitarCompra(rota []Trecho) {
	// Convertendo a rota para JSON
	jsonData, err := json.Marshal(rota)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return
	}
	var req *http.Request

	// Descobrir qual servidor será enviada a requisição
	servidor := rota[0].Comp
	if servidor == "A" {
		// Criando a requisição POST
		//req, err = http.NewRequest("POST", "http://server1:8000/compras", bytes.NewBuffer(jsonData))
		req, err = http.NewRequest("POST", "http://localhost:8000/compras", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Erro ao criar a requisição:", err)
			return
		}
	} else if servidor == "B" {
		// Criando a requisição POST
		//req, err = http.NewRequest("POST", "http://server2:8001/compras", bytes.NewBuffer(jsonData))
		req, err = http.NewRequest("POST", "http://localhost:8001/compras", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Erro ao criar a requisição:", err)
			return
		}

	} else if servidor == "C" {
		// Criando a requisição POST
		//		req, err = http.NewRequest("POST", "http://server3:8002/compras", bytes.NewBuffer(jsonData))
		req, err = http.NewRequest("POST", "http://localhost:8002/compras", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Erro ao criar a requisição:", err)
			return
		}
	} else {
		fmt.Println("Servidor inválido!")
	}
	// Definindo o cabeçalho Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Criando o cliente HTTP
	client := &http.Client{}

	// Enviando a requisição
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao enviar a requisição:", err)
		return
	}
	defer resp.Body.Close()

	// Verificando a resposta
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Compra enviada com sucesso!")
	} else {
		fmt.Printf("Erro na solicitação: %s\n", resp.Status)
	}
}
