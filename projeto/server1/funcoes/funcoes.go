package funcoes

import (
	// "bytes"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"

	// "log"
	"net/http"
	"os"
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

type RetornoCompra struct {
	Resultado bool   `json:"Resultado"`
	Server    string `json:"Server"`
	Compra    Compra `json:"Compra"`
}

type ReqRotas struct{
	Origem string `json:"Origem"`
	Destino string `json:"Destino"`
}

var TrechoLivre = make([]bool, 100)
var FilaRequest = make(map[string]PrepareRequest)
var Rotas map[string][]Trecho
var filePathRotas = "dados/rotas.json" //caminho para arquivo de rotas
var mutex sync.Mutex
var mutexVagas sync.Mutex
var mutexCommit sync.Mutex
var server2 = "http://server2:"
var server3 = "http://server3:"

func BuscarRotaServidor(servidor string) map[string][]Trecho {
	fmt.Println("Função de buscar as rotas nos servidores")
	LerRotas()
	fmt.Println("Servidor ",servidor)
	// Inicializa o mapa
	trechos := make(map[string][]Trecho)

	var resp *http.Response
	var err error

	// Condicional para verificar o servidor
	if servidor == "A" {
		// BUSCA NO SERVIDOR 1
		LerRotas() // Certifique-se de que LerRotas popula corretamente o mapa Rotas
		trechos = Rotas
		fmt.Println("Dentro de if A")
	} else if servidor == "B" {
		// BUSCA NO SERVIDOR 2
		fmt.Println("Dentro de if B")
		resp, err = http.Get(server2 + "8001/rota")
	} else if servidor == "C" {
		// BUSCA NO SERVIDOR 3
		fmt.Println("Dentro de if C")
		resp, err = http.Get(server3 + "8002/rota")
	} else {
		fmt.Println("Servidor desconhecido:", servidor)
		return nil
	}

	if(servidor == "B" || servidor == "C"){

		// Tratamento de erro da requisição
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
	
		// Decodificando o JSON para o mapa
		err = json.Unmarshal(body, &trechos)
		if err != nil {
			fmt.Println("Erro ao converter o JSON:", err)
			return nil
		}
	}

	return trechos
}

func BuscaRotas(w http.ResponseWriter, r *http.Request) {
	var reqrotas ReqRotas
	err := json.NewDecoder(r.Body).Decode(&reqrotas)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}
	fmt.Println("Função buscar rotas, rota: ", reqrotas)

	// Inicializa o mapa
	rotas := make(map[string][]Trecho)

	// Chama os servidores e adiciona os trechos, se disponíveis
	if trechosA := BuscarRotaServidor("A"); trechosA != nil {
		for chave, valor := range trechosA {
			rotas[chave] = append(rotas[chave], valor...)
		}
	}
	fmt.Println("Depois de solicitar rotas de A")

	if trechosB := BuscarRotaServidor("B"); trechosB != nil {
		for chave, valor := range trechosB {
			rotas[chave] = append(rotas[chave], valor...)
		}
	}
	fmt.Println("Depois de solicitar rotas de A")
	if trechosC := BuscarRotaServidor("C"); trechosC != nil {
		for chave, valor := range trechosC {
			rotas[chave] = append(rotas[chave], valor...)
		}
	}
	fmt.Println("Depois de solicitar rotas de A")
	// Busca todos os caminhos a partir dos dados combinados de `rotas`
	menoresCaminhos := EncontrarTodosCaminhos(rotas, reqrotas.Origem, reqrotas.Destino)
	fmt.Printf("Retorno do buscar rotas: %v\n", menoresCaminhos)

	// Envia a resposta como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menoresCaminhos)
}

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
	mutexVagas.Lock() // Adquire o bloqueio
	defer mutexVagas.Unlock()
	for _, trecho := range trechos {
		if trecho.Comp == "A" {
			for i, x := range Rotas[trecho.Origem] {
				if trecho.ID == x.ID {
					Rotas[trecho.Origem][i].Vagas = x.Vagas - 1
					fmt.Println("Vagas: ", Rotas[trecho.Origem][i].Vagas)
				}
			}
		}
	}
	SalvarRotas()
}

// Pega todas as rotas do arquivo json
func GetRotas(w http.ResponseWriter, r *http.Request) {
	LerRotas()
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
		if trecho.Comp == "A" {
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
	if server == "A" {
		var id int
		for _, trecho := range Request.Compra.Trechos {
			if trecho.Comp == "A" {

				// Convertendo a string para int
				id, err = strconv.Atoi(trecho.ID)
				if err != nil {
					fmt.Println("Erro ao converter ID:", err)
					return false
				}
				ok = ok && TrechoLivre[id]                //verifica se não tem outro processo fazendo alteração no trecho no momento
				ok = ok && VerificaVagasTrecho(trecho.ID) //verifica se há vagas no trecho
			}
		}
		if ok { //caso os trechos estiverem livres e tenham vagas, eles são reservados
			TrechoLivre[id] = false //trava o trecho
			ReservarTrechos(Request)
		}
		return ok
		// envia a mensagem para o servidor 2 se preparar
	} else if server == "B" {
		fmt.Println("\nEnviando compra para servidor 2")
		//req, err = http.NewRequest("POST", "http://server2:8001/compras/preparar", bytes.NewBuffer(jsonData))
		req, err = http.NewRequest("POST", server2+"8001/compras/preparar", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Erro ao criar a requisição:", err)
			return false
		}
		//envia mensagem para o servidor 3 se preparar
	} else if server == "C" {
		fmt.Println("\nEnviando compra para servidor 3")
		//req, err = http.NewRequest("POST", "http://server3:8002/compras/preparar", bytes.NewBuffer(jsonData))
		req, err = http.NewRequest("POST", server3+"8002/compras/preparar", bytes.NewBuffer(jsonData))
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
		if server == "A" {
			//cancelar o commit no servidor 1
			request := FilaRequest[idTransacao]
			for _, trecho := range request.Compra.Trechos {
				if trecho.Comp == "A" {
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

		} else if server == "B" {
			// envia a solicitação de cancelar commit para o servidor 2
			//resp, err := http.Post("http://server2:8001/compras/cancelar", "application/json", bytes.NewBuffer(jsonData))
			resp, err := http.Post(server2+"8001/compras/cancelar", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Erro ao enviar request:", err)
				return
			}
			defer resp.Body.Close()

		} else if server == "C" {
			// envia a solicitação de cancelar commit para o servidor 3
			//resp, err := http.Post("http://server3:8002/compras/cancelar", "application/json", bytes.NewBuffer(jsonData))
			resp, err := http.Post(server3+"8002/compras/cancelar", "application/json", bytes.NewBuffer(jsonData))
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
		if server == "A" {
			//servidor 1(Atual)
			request := FilaRequest[idTransacao]
			//subtrai as vagas
			/*

				precisa de uma função q salva a compra aq
				pode ser junto com a de subtrair vagas

			*/
			SubtrairVagas(request.Compra.Trechos)
			for _, trecho := range request.Compra.Trechos {
				if trecho.Comp == "A" {
					//deixa os trechos livres para serem alterados
					id := ConverteID(trecho.ID)
					TrechoLivre[id] = true
				}
			}
			//remove o request
			delete(FilaRequest, idTransacao)

		} else if server == "B" {
			//envia mensagem de confirmação para o servidor 2
			resp, err := http.Post(server2+"8001/compras/confirmar", "application/json", bytes.NewBuffer(jsonData))
			//resp, err := http.Post("http://server2:8001/compras/confirmar", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Erro ao enviar request:", err)
				return
			}
			defer resp.Body.Close()

		} else if server == "C" {
			//envia mensagem de confirmação para o servidor 3
			//resp, err := http.Post("http://server3:8002/compras/confirmar", "application/json", bytes.NewBuffer(jsonData))
			resp, err := http.Post(server3+"8002/compras/confirmar", "application/json", bytes.NewBuffer(jsonData))
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

	fmt.Println("\n\nCompra recebida pelo servidor 1: ", compra)

	transactionID := uuid.New().String() //cria o id da transação
	var transacao = PrepareRequest{      // determina o dado q será enviado para preparação
		Compra:        compra,
		TransactionID: transactionID,
	}

	//percorre os servidores participantes da compra para poder mandar a compra para eles
	for _, participante := range compra.Participantes {
		contador := 0
		for contador <= 10 {
			fmt.Print(contador)
			//envia a requisição de preparação para os outros servidores
			result := EnviarRequestPreparacao(participante, transacao)
			fmt.Println("Retorno de ", participante, " ", result)

			if !result && contador == 10 { //verifica se todos os servidores conseguiram preparar
				fmt.Println("\nEntrou no contador")
				//Cancela o commit caso algum servidor não tenha conseguido preparar para o commit
				CancelarTransacao(transactionID, compra.Participantes)
				// retorna que a compra não teve sucesso
				retorno := RetornoCompra{
					Resultado: false,
					Server:    "A",
					Compra:    compra,
				}

				// Serializando a resposta em JSON
				response, err := json.Marshal(retorno)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					// http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				fmt.Println("\nResposta da solicitação de compra 1: ", retorno)
				// Enviando a resposta
				w.WriteHeader(http.StatusOK)
				w.Write(response)
				return //acaba a solicitação
			} else if result {
				break
			}
			contador++
			time.Sleep(500 * time.Millisecond)
		}

	}
	// caso todos os servidores conseguirem se preparar para o commit, então o commit é realizado
	ConfirmarTransacao(transactionID, compra.Participantes)
	//retorna que a compra foi bem sucedida
	retorno := RetornoCompra{
		Resultado: true,
		Server:    "A",
		Compra:    compra,
	}
	fmt.Println("\nResposta da solicitação de compra 2: ", retorno)
	// Serializando a resposta em JSON
	response, err := json.Marshal(retorno)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enviando a resposta
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

//Funções que serão usadas quando esse servidor não for o coordenador e for apenas participante

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
	fmt.Println("\n Mensagem de preparação: ", dados)
	var id int
	for _, trecho := range dados.Compra.Trechos {
		if trecho.Comp == "A" {
			// Convertendo a string para int
			id, err = strconv.Atoi(trecho.ID)
			if err != nil {
				fmt.Println("Erro ao converter ID:", err)
				return
			}

			mutexCommit.Lock()
			ok = ok && TrechoLivre[id]                //verifica se não tem outro processo fazendo alteração no trecho no momento
			ok = ok && VerificaVagasTrecho(trecho.ID) //verifica se há vagas no trecho
			if ok {
				entra_if = true
			}
			mutexCommit.Unlock()
		}
		if ok && entra_if { //caso os trechos estiverem livres e tenham vagas, eles são reservados
			TrechoLivre[id] = false //trava o trecho
			ReservarTrechos(dados)
		}
	}

	var result bool = ok

	// Define o código de status e o tipo de conteúdo como texto simples
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Println("\nResposta retornada: ", result)
	// Escreve o valor do booleano como uma string ("true" ou "false")
	fmt.Fprintf(w, "%t", result)
}

// faz a mesma coisa da função ConfirmarTransacao
func ConfirmarCommit(w http.ResponseWriter, r *http.Request) {
	var dados CommitRequest

	err := json.NewDecoder(r.Body).Decode(&dados)
	if err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	trechosCompra := FilaRequest[dados.TransactionID].Compra.Trechos
	fmt.Println("\n Requisição a ser confirmada: ", dados.TransactionID)

	SubtrairVagas(trechosCompra)
	for _, trecho := range trechosCompra {
		if trecho.Comp == "A" {
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
	fmt.Println("\n Requisição a ser cancelada: ", dados.TransactionID)
	_, existe := FilaRequest[dados.TransactionID]
	if existe { //verifica se foi essa requisição que fez o bloqueio do trecho
		for _, trecho := range trechosCompra {
			if trecho.Comp == "A" {
				//cancelar o commit no servidor 1
				id := ConverteID(trecho.ID)
				TrechoLivre[id] = true

			}
		}
		delete(FilaRequest, dados.TransactionID)
	}

}

func VerCompras(w http.ResponseWriter, r *http.Request) {

}

func LerRotas() {
	// Abra o arquivo
	file, err := os.Open(filePathRotas)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Criar um mapa para armazenar as rotas
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
