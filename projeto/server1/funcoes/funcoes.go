package funcoes

import (
	// "bytes"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
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
	Nome      string
	Sobrenome string
	Cpf       string
}
type Compra struct {
	Pessoa  Pessoa
	Trechos []Trecho
	Participantes []string
}

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
var Rotas map[string][]Trecho
var filePathRotas = "dados/rotas.json" //caminho para arquivo de rotas

func ConverteID(idstring string)int{
	id, err := strconv.Atoi(idstring)
	if err != nil {
		fmt.Println("Erro ao converter ID:", err)
		return 0
	}
	return id
}

func SalvarRotas(){}

func SubtrairVagas( trechos[]Trecho){
	for _,trecho := range trechos{
		if trecho.Comp == "A"{
			for _, x := range Rotas[trecho.Origem]{
				if trecho.ID == x.ID{
					x.Vagas= x.Vagas-1
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

func VerificaVagasTrecho(id string)bool{
	LerRotas()
	for _,valor := range Rotas{
		for _,trecho := range valor{
			if trecho.ID == id && trecho.Vagas > 0{
				return true
			}
		}
	}
	return false
}

func ReservarTrechos(Request PrepareRequest){
	for _ ,trecho := range Request.compra.Trechos{
		if trecho.Comp == "A"{
			id := ConverteID(trecho.ID)
			TrechoLivre[id]= false

		}
	}
	FilaRequest[Request.TransactionID] = Request
}
// função para enviar aos servidores a mensagem de preparação para commit
// os servidores retornam se conseguiram se prapar ou não
func EnviarRequestPreparacao(server string, Request PrepareRequest)bool{
	var ok  = true //variavel para pegar a resposta se o servidor conseguiu se preparar ou não

	var req *http.Request //variavel da requisição

	//transforma dados em json
	jsonData, err := json.Marshal(Request.compra)
	if err != nil {
		fmt.Println("Erro ao converter para JSON:", err)
		return false
	}

	//Verificar o própio servidor
	if server == "A"{
		for _, trecho := range Request.compra.Trechos{
			if trecho.Comp == "A"{
				
				// Convertendo a string para int
				id, err := strconv.Atoi(trecho.ID)
				if err != nil {
					fmt.Println("Erro ao converter ID:", err)
					return false
   				}	
				ok = ok && TrechoLivre[id]//verifica se não tem outro processo fazendo alteração no trecho no momento
				ok = ok && VerificaVagasTrecho(trecho.ID)//verifica se há vagas no trecho
			}
		}
		if ok{//caso os trechos estiverem livres e tenham vagas, eles são reservados
			ReservarTrechos(Request)
		}
		return ok
	// envia a mensagem para o servidor 2 se preparar	
	}else if server == "B"{
		req, err = http.NewRequest("POST", "http://localhost:8001/compras/preparar", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Erro ao criar a requisição:", err)
			return false
		}
	//envia mensagem para o servidor 3 se preparar
	}else if server == "C"{
		req, err = http.NewRequest("POST", "http://localhost:8001/compras/preparar", bytes.NewBuffer(jsonData))
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
        fmt.Println("Erro ao decodificar JSON:", err)
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
//função para cancelar o commit pois algum dos servidores não conseguiu realizar a preparação
func CancelarTransacao(idTransacao string, participantes []string){
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
	for _, server := range participantes{
		if server == "A"{
			//cancelar o commit no servidor 1
			request := FilaRequest[idTransacao]
			for _,trecho :=range request.compra.Trechos{
				if trecho.Comp == "A"{
					//deixa os trechos livres para serem alterados
					id := ConverteID(trecho.ID)
					TrechoLivre[id]= true
				}
			}
			//remove o request 
			delete(FilaRequest, idTransacao)


		}else if server =="B"{
			// envia a solicitação de cancelar commit para o servidor 2
			resp, err := http.Post("http://localhost:8001/compras/cancelar", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Erro ao enviar request:", err)
				return
			}
    		defer resp.Body.Close()

		}else if server == "C"{
			// envia a solicitação de cancelar commit para o servidor 3
			resp, err := http.Post("http://localhost:8002/compras/cancelar", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Erro ao enviar request:", err)
				return
			}
    		defer resp.Body.Close()
		}
		

	}

}
//Função para confirmar o commit, qunado todos os servidores conseguiram preparar o commit
func ConfirmarTransacao(idTransacao string, participantes []string){
	//Só é necessário enviar o Id da transação
	var transacao = CommitRequest{TransactionID: idTransacao}//dado da transação
	//convertendo pra json
	jsonData, err := json.Marshal(transacao)
    if err != nil {
        fmt.Println("Erro ao converter dados para JSON:", err)
        return
    }

	//mandar para todos os servidores participantes
	for _, server := range participantes{
		if server == "A"{
			//servidor 1(Atual)
			request := FilaRequest[idTransacao]
			//subtrai as vagas
			/*
			
			precisa de uma função q salva a compra aq
			pode ser junto com a de subtrair vagas
			
			*/
			SubtrairVagas(request.compra.Trechos)
			for _,trecho :=range request.compra.Trechos{
				if trecho.Comp == "A"{
					//deixa os trechos livres para serem alterados
					id := ConverteID(trecho.ID)
					TrechoLivre[id]= true
				}
			}
			//remove o request 
			delete(FilaRequest, idTransacao)


			
			
		}else if server =="B"{
			//envia mensagem de confirmação para o servidor 2
			resp, err := http.Post("http://localhost:8001/compras/confirmar", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Erro ao enviar request:", err)
				return
			}
    		defer resp.Body.Close()

		}else if server == "C"{
			//envia mensagem de confirmação para o servidor 3
			resp, err := http.Post("http://localhost:8002/compras/confirmar", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Erro ao enviar request:", err)
				return
			}
    		defer resp.Body.Close()
		}
		

	}

}


//Essa função é chamada apenas quando o servidor é o coordenador
//função chamada para realizar uma compra
func SolicitacaoCord(w http.ResponseWriter, r *http.Request) {
	// Cria uma instância da struct Compra
    var compra Compra

    // Decodifica o corpo da requisição em JSON
    err := json.NewDecoder(r.Body).Decode(&compra)
    if err != nil {
        http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
        return
    }

	transactionID := uuid.New().String()//cria o id da transação
	var transacao = PrepareRequest{// determina o dado q será enviado para preparação
		compra: compra,
		TransactionID: transactionID,
	}
	//percorre os servidores participantes da compra para poder mandar a compra para eles
	for _,participante := range compra.Participantes{
		//envia a requisição de preparação para os outros servidores
		result := EnviarRequestPreparacao(participante, transacao) 
		if !result{//verifica se todos os servidores conseguiram preparar
			//Cancela o commit caso algum servidor não tenha conseguido preparar para o commit
			CancelarTransacao(transactionID, compra.Participantes)
			return
		}
	} 
	// caso todos os servidores conseguirem se preparar para o commit, então o commit é realizado
	ConfirmarTransacao(transactionID, compra.Participantes)
	
}

//Funções que serão usadas quando esse servidor não for o coordenador e for apenas participante

//Função de preparação, vai verificar se o trecho está em uso e se tem vagas nos trechos
//Caso esteja livre reserva os trechos e retorna true
func Commit(w http.ResponseWriter, r *http.Request){}

//faz a mesma coisa da função ConfirmarTransacao
func ConfirmarCommit(w http.ResponseWriter, r *http.Request){}

//Faz a mesma coisa da função CancelarTransacao
func CancelarCommit(w http.ResponseWriter, r *http.Request){}

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
