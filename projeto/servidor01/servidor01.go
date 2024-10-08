package main

import (
	"fmt"
	"net"
	"projeto/servidor01/funcoesServer01"
)

var ADRESS string = "0.0.0.0:22356"

func main() {
	// Escutando na porta 22356
	ln, err := net.Listen("tcp", ADRESS)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer ln.Close()
	fmt.Printf("Servidor iniciado em: %s", ADRESS)

	for {
		// Aceitando conexões dos clientes
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		//cria uma gorotine para cada conexão
		go funcoesServer01.HandleConnection(conn)
	}
}
