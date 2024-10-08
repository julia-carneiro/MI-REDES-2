package funcoesServer02

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
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

//Conecta com o servidor
func ConectarServidor(ADRESS string) net.Conn {
	// Conectando ao servidor na porta 8080
	conn, err := net.Dial("tcp", ADRESS)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return nil
	}

	return conn
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("\nCliente conectado:", conn.RemoteAddr())

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Erro ao ler a mensagem:", err)
		return
	}

	var dados Dados
	err = json.Unmarshal([]byte(message), &dados)
	if err != nil {
		conn.Write([]byte("Erro no formato dos dados enviados. Esperado JSON.\n"))
		return
	}

	fmt.Println("\nMensagem recebida do cliente:", dados)

	
}