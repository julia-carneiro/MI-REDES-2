package funcoes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// type Request int

// const ( //Tipos de mensagens que podem ser enviadas ao servidor
// 	ROTAS Request = iota
// 	COMPRA
// 	CADASTRO
// 	LERCOMPRAS
// )

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
	Pessoa  Pessoa `json:"Pessoa"`
	Trechos []Trecho `json:"Trechos"`
	Participantes []string `json:"Participantes"`
}

type PrepareRequest struct {
    TransactionID string `json:"TransationID"`     // ID único da transação
    Compra       Compra `json:"Compra"`  // Rotas apenas para este servidor
}
type CommitRequest struct {
    TransactionID string // ID da transação a ser confirmada
}
type CancelRequest struct {
    TransactionID string // ID da transação a ser cancelada
}

var TrechoLivre = make([]bool, 10)
var FilaRequest = make(map[string]PrepareRequest)
var rotas map[string][]Trecho
var filePathRotas = "dados/rotas.json" //caminho para arquivo de rotas

// Pega todas as rotas do arquivo json
func GetRotas(w http.ResponseWriter, r *http.Request) {
	rotas = LerRotas()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rotas)
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
	rotas := make(map[string][]Trecho)

	// Decodificar o arquivo JSON para o mapa
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&rotas); err != nil {
		fmt.Println("Erro ao decodificar o JSON:", err)
		return nil
	}
	return rotas
}

//Função de preparação, vai verificar se o trecho está em uso e se tem vagas nos trechos
//Caso esteja livre reserva os trechos e retorna true
func Commit(w http.ResponseWriter, r *http.Request){
	fmt.Println("server2")
	var dados PrepareRequest

    err := json.NewDecoder(r.Body).Decode(&dados)
    if err != nil {
        http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
        return
    }
	// Exemplo de lógica que determina o valor do booleano a ser retornado
    var result bool
    // Lógica para definir o valor de result
    result = true // ou false, dependendo da lógica do seu sistema

    // Define o código de status e o tipo de conteúdo como texto simples
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "text/plain")

    // Escreve o valor do booleano como uma string ("true" ou "false")
    fmt.Fprintf(w, "%t", result)
	
}

//faz a mesma coisa da função ConfirmarTransacao
func ConfirmarCommit(w http.ResponseWriter, r *http.Request){
	
}

//Faz a mesma coisa da função CancelarTransacao
func CancelarCommit(w http.ResponseWriter, r *http.Request){

}