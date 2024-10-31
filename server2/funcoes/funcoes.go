package funcoes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/google/uuid"
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
	Pessoa        Pessoa   `json:"Pessoa"`
	Trechos       []Trecho `json:"Trechos"`
	Participantes []string `json:"Participantes"`
}

type PrepareRequest struct {
	TransactionID string `json:"TransationID"` // ID único da transação
	Compra        Compra `json:"Compra"`       // Rotas apenas para este servidor
}
type CommitRequest struct {
	TransactionID string // ID da transação a ser confirmada
}
type CancelRequest struct {
	TransactionID string // ID da transação a ser cancelada
}

var TrechoLivre = make([]bool, 100)

var FilaRequest = make(map[string]PrepareRequest)
var Rotas map[string][]Trecho
var filePathRotas = "dados/rotas.json" //caminho para arquivo de Rotas
var mutex sync.Mutex

func ConverteID(idstring string) int {
	id, err := strconv.Atoi(idstring)
	if err != nil {
		fmt.Println("Erro ao converter ID:", err)
		return 0
	}
	return id
}

func SalvarRotas() {
	mutex.Lock()         // Adquire o bloqueio
	defer mutex.Unlock() // Garante que o bloqueio será liberado

	file, err := os.Create(filePathRotas)
	if err != nil {
		fmt.Println("Erro ao escrever:", err)
		return 
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(Rotas)
}


func SubtrairVagas(trechos []Trecho) {
	for _, trecho := range trechos {
		if trecho.Comp == "B" {
			for _, lista := range Rotas[trecho.Origem] {
				for index,valor := range lista{
					if trecho.ID == valor.ID{
						fmt.Println("Server 2 antes:",valor.Vagas)
						Rotas[trecho.Origem][index].Vagas = valor.Vagas -1
						fmt.Println("Server 2 depois:", valor.Vagas)
					}
				}
				
			}
		}
	}
	SalvarRotas()
}

// Pega todas as rotas do arquivo json
func GetRotas(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Rotas)
}

func VerificaVagasTrecho(id string) bool {
	LerRotas()
	for _, valor := range Rotas {
		for _, trecho := range valor {
			if trecho.ID == id && trecho.Vagas > 0 {
				return true
			}
		}
	}
	return false
}

func ReservarTrechos(Request PrepareRequest) {
	for _, trecho := range Request.Compra.Trechos {
		if trecho.Comp == "B" {
			id := ConverteID(trecho.ID)
			TrechoLivre[id] = false

		}
	}
	FilaRequest[Request.TransactionID] = Request
}

// função para enviar aos servidores a mensagem de preparação para commit
// os servidores retornam se conseguiram se prapar ou não
func EnviarRequestPreparacao(server string, Request PrepareRequest) bool {

	var ok = true //variavel para pegar a resposta se o servidor conseguiu se preparar ou não

	var req *http.Request //variavel da requisição

	//transforma dados em json
	jsonData, err := json.Marshal(Request)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return false
	}

	//Verificar o própio servidor
	if server == "B" {
		var id int
		for _, trecho := range Request.Compra.Trechos {
			if trecho.Comp == "B" {

				// Convertendo a string para int
				id, err := strconv.Atoi(trecho.ID)
				if err != nil {
					fmt.Println("Erro ao converter ID:", err)
					return false
				}
				fmt.Println(TrechoLivre[id])
				ok = ok && TrechoLivre[id]                //verifica se não tem outro processo fazendo alteração no trecho no momento
				ok = ok && VerificaVagasTrecho(trecho.ID) //verifica se há vagas no trecho
			}
		}
		if ok { //caso os trechos estiverem livres e tenham vagas, eles são reservados
			TrechoLivre[id] = false //trava o trecho
			ReservarTrechos(Request)
		}
		return ok
		// envia a mensagem para o servidor A se preparar
	} else if server == "A" {
		fmt.Println("\nEnviando compra para servidor A")
		req, err = http.NewRequest("POST", "http://localhost:8000/compras/preparar", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Erro ao criar a requisição:", err)
			return false
		}
		//envia mensagem para o servidor C se preparar
	} else if server == "C" {
		fmt.Println("\nEnviando compra para servidor C")
		req, err = http.NewRequest("POST", "http://localhost:8002/compras/preparar", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Erro ao criar a requisição:", err)
			return false
		}

	}
	// Definindo o cabeçalho Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Criando o cliente HTTP
	client := &http.Client{}

	// Enviando a requisição
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao enviar a requisição:", err)
		return false
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return false
	}

	var dadosResposta bool //resposta da preparação. Se conseguiu praparar ou não
	err = json.Unmarshal(body, &dadosResposta)
	if err != nil {
		fmt.Println("Erro aq ao decodificar JSON:", err)
		return false
	}

	// Verificando a resposta
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Compra enviada com sucesso!")
	} else {
		fmt.Printf("Erro na solicitação: %s\n", resp.Status)
	}
	return dadosResposta //retorna o resultado da preparação
}

// função para cancelar o commit pois algum dos servidores não conseguiu realizar a preparação
func CancelarTransacao(idTransacao string, participantes []string) {
	/*O sevidor participante só precisa ter acesso ao Id da transação, a partir disso
	ele tem acesso a fila de requisiçõese lá tem a compra com todos os trechos
	*/
	//Arquivo de transação
	var transacao = CancelRequest{TransactionID: idTransacao}
	//convertendo pra json
	jsonData, err := json.Marshal(transacao)
	if err != nil {
		fmt.Println("Erro ao converter dados para JSON:", err)
		return
	}
	//envia para toods os servers
	for _, server := range participantes {
		if server == "B" {
			//cancelar o commit no servidor B
			request := FilaRequest[idTransacao]
			for _, trecho := range request.Compra.Trechos {
				if trecho.Comp == "B" {
					//deixa os trechos livres para serem alterados
					_, existe := FilaRequest[idTransacao]
					if existe { //verifica se foi essa requisição que fez o bloqueio do trecho
						id := ConverteID(trecho.ID)
						TrechoLivre[id] = true
					}
				}
			}
			_, existe := FilaRequest[idTransacao]
			if existe {
				//remove o request
				delete(FilaRequest, idTransacao)
			}

		} else if server == "A" {
			// envia a solicitação de cancelar commit para o servidor A
			resp, err := http.Post("http://localhost:8000/compras/cancelar", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Erro ao enviar request:", err)
				return
			}
			defer resp.Body.Close()

		} else if server == "C" {
			// envia a solicitação de cancelar commit para o servidor C
			resp, err := http.Post("http://localhost:8002/compras/cancelar", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Erro ao enviar request:", err)
				return
			}
			defer resp.Body.Close()
		}

	}

}

// Função para confirmar o commit, qunado todos os servidores conseguiram preparar o commit
func ConfirmarTransacao(idTransacao string, participantes []string) {
	//Só é necessário enviar o Id da transação
	var transacao = CommitRequest{TransactionID: idTransacao} //dado da transação
	//convertendo pra json
	jsonData, err := json.Marshal(transacao)
	if err != nil {
		fmt.Println("Erro ao converter dados para JSON:", err)
		return
	}

	//mandar para todos os servidores participantes
	for _, server := range participantes {
		if server == "B" {
			//servidor B(Atual)
			request := FilaRequest[idTransacao]
			//subtrai as vagas
			/*

				precisa de uma função q salva a compra aq
				pode ser junto com a de subtrair vagas

			*/
			SubtrairVagas(request.Compra.Trechos)
			for _, trecho := range request.Compra.Trechos {
				if trecho.Comp == "B" {
					//deixa os trechos livres para serem alterados
					id := ConverteID(trecho.ID)
					TrechoLivre[id] = true
				}
			}
			//remove o request
			delete(FilaRequest, idTransacao)

		} else if server == "A" {
			//envia mensagem de confirmação para o servidor A
			resp, err := http.Post("http://localhost:8000/compras/confirmar", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Erro ao enviar request:", err)
				return
			}
			defer resp.Body.Close()

		} else if server == "C" {
			//envia mensagem de confirmação para o servidor C
			resp, err := http.Post("http://localhost:8002/compras/confirmar", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Erro ao enviar request:", err)
				return
			}
			defer resp.Body.Close()
		}

	}

}

// Essa função é chamada apenas quando o servidor é o coordenador
// função chamada para realizar uma compra
func SolicitacaoCord(w http.ResponseWriter, r *http.Request) {
	// Cria uma instância da struct Compra
	var compra Compra

	// Decodifica o corpo da requisição em JSON
	err := json.NewDecoder(r.Body).Decode(&compra)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	transactionID := uuid.New().String() //cria o id da transação
	var transacao = PrepareRequest{      // determina o dado q será enviado para preparação
		Compra:        compra,
		TransactionID: transactionID,
	}
	//percorre os servidores participantes da compra para poder mandar a compra para eles
	for _, participante := range compra.Participantes {
		//envia a requisição de preparação para os outros servidores
		result := EnviarRequestPreparacao(participante, transacao)
		fmt.Println("Retorno de ", participante, " ", result)
		if !result { //verifica se todos os servidores conseguiram preparar
			//Cancela o commit caso algum servidor não tenha conseguido preparar para o commit
			CancelarTransacao(transactionID, compra.Participantes)
			return
		}
	}
	// caso todos os servidores conseguirem se preparar para o commit, então o commit é realizado
	ConfirmarTransacao(transactionID, compra.Participantes)

}

// Função de preparação, vai verificar se o trecho está em uso e se tem vagas nos trechos
// Caso esteja livre reserva os trechos e retorna true
func Commit(w http.ResponseWriter, r *http.Request) {

	var dados PrepareRequest
	ok := true
	entra_if := false
	err := json.NewDecoder(r.Body).Decode(&dados)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}
	var id int
	for _, trecho := range dados.Compra.Trechos {
		if trecho.Comp == "B" {
			// Convertendo a string para int
			id, err = strconv.Atoi(trecho.ID)
			if err != nil {
				fmt.Println("Erro ao converter ID:", err)
				return
			}

			ok = ok && TrechoLivre[id]              //verifica se não tem outro processo fazendo alteração no trecho no momento
			ok = ok && VerificaVagasTrecho(trecho.ID) //verifica se há vagas no trecho
			if ok { entra_if = true }
		}
		if ok && entra_if{ //caso os trechos estiverem livres e tenham vagas, eles são reservados
			TrechoLivre[id] = false //trava o trecho
			ReservarTrechos(dados)
		}
	}

	var result bool = ok

	// Define o código de status e o tipo de conteúdo como texto simples
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")

	// Escreve o valor do booleano como uma string ("true" ou "false")
	fmt.Fprintf(w, "%t", result)
}


// faz a mesma coisa da função ConfirmarTransacao
func ConfirmarCommit(w http.ResponseWriter, r *http.Request) {
	var dados CommitRequest
	fmt.Println("Commit server 2")

	err := json.NewDecoder(r.Body).Decode(&dados)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("DADOS", dados)

	trechosCompra := FilaRequest[dados.TransactionID].Compra.Trechos
	fmt.Println("Trechos compra:", trechosCompra)

	SubtrairVagas(trechosCompra)
	for _, trecho := range trechosCompra {
		if trecho.Comp == "B" {
			//deixa os trechos livres para serem alterados
			id := ConverteID(trecho.ID)
			TrechoLivre[id] = true
		}
	}
	//remove o request
	delete(FilaRequest, dados.TransactionID)

}

// Faz a mesma coisa da função CancelarTransacao
func CancelarCommit(w http.ResponseWriter, r *http.Request) {
	var dados CancelRequest
	err := json.NewDecoder(r.Body).Decode(&dados)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}
	trechosCompra := FilaRequest[dados.TransactionID].Compra.Trechos
	_, existe := FilaRequest[dados.TransactionID]
	if existe { //verifica se foi essa requisição que fez o bloqueio do trecho
		for _, trecho := range trechosCompra {
			if trecho.Comp == "B" {
				//cancelar o commit no servidor 1
				id := ConverteID(trecho.ID)
				TrechoLivre[id] = true

			}
		}
		delete(FilaRequest, dados.TransactionID)
	}

}

func LerRotas() {
	// Abra o arquivo
	file, err := os.Open(filePathRotas)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Criar um mapa para armazenar as Rotas
	rotas := make(map[string][]Trecho)

	// Decodificar o arquivo JSON para o mapa
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&rotas); err != nil {
		fmt.Println("Erro ao decodificar o JSON:", err)
		return
	}
	Rotas = rotas
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