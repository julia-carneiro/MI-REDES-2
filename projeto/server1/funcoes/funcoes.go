package funcoes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"mux"
	"net/http"
	"github.com/gorilla/mux"
	//  "os"
	// "sync"
	// "path/filepath"
)

type Request int

const (//Tipos de mensagens que podem ser enviadas ao servidor
	ROTAS Request = iota
	COMPRA
	CADASTRO
	LERCOMPRAS
)

type Compra struct {//Estrura de dados de compra
	Cpf     string   `json:"Cpf"`
	Caminho []string `json:"Caminho"`
}

type User struct {//Estrutura de dados do cliente
	Cpf  string `json:"Cpf"`
}

type Rota struct {//Estrura de dados de uma rota
	Destino string `json:"Destino"`
	Vagas   int    `json:"Vagas"`
	Peso    int    `json:"Peso"`
}

type Dados struct {//Estrutura de dados da mensagem recebida no cliente
	Request      Request `json:"Request"`
	DadosCompra  *Compra `json:"DadosCompra"`
	DadosUsuario *User   `json:"DadosUsuario"`
}


func GetRotas(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Rota acessada!")
}