package funcoes

import (
	"fmt"
	"net"
)

type Pessoa struct {
	Nome      string
	Sobrenome string
	Cpf       string
}
type Compra struct {
	Origem  string
	Destino string
	Comp    string
	Pessoa  Pessoa
}

// Conecta com o servidor
func ConectarServidor(ADRESS string) net.Conn {
	// Conectando ao servidor na porta 8080
	conn, err := net.Dial("tcp", ADRESS)
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return nil
	}

	return conn
}
