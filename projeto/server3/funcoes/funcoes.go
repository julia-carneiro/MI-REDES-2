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

type Compra struct { //Estrura de dados de compra
	Cpf     string   `json:"Cpf"`
	Caminho []string `json:"Caminho"`
}
type User struct { //Estrutura de dados do cliente
	Cpf string `json:"Cpf"`
}

type Trecho struct {
	Origem  string `json:"Origem"`
	Destino string `json:"Destino"`
	Vagas   int    `json:"Vagas"`
	Peso    int    `json:"Peso"`
	Comp    string `json:"Comp"`
	ID      string `json:"ID"`
}

// type Dados struct { //Estrutura de dados da mensagem recebida no cliente
// 	Request      Request `json:"Request"`
// 	DadosCompra  *Compra `json:"DadosCompra"`
// 	DadosUsuario *User   `json:"DadosUsuario"`
// }
type PrepareRequest struct {
    TransactionID string     // ID único da transação
    compra       Compra   // Rotas apenas para este servidor
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
	var dados PrepareRequest

    err := json.NewDecoder(r.Body).Decode(&dados)
    if err != nil {
        http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
        return
    }
	fmt.Print("Servidor 2: ", dados.compra)
	
}

//faz a mesma coisa da função ConfirmarTransacao
func ConfirmarCommit(w http.ResponseWriter, r *http.Request){

}

//Faz a mesma coisa da função CancelarTransacao
func CancelarCommit(w http.ResponseWriter, r *http.Request){

}
