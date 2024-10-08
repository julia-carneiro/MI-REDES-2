package main

import (
	// "bufio"
	// "fmt"
	// "net"
	// "os"
	// "sessao3/cliente/funcoesCliente"
	// "strings"
	
)
/*
- cliente escolhe a rota normalmente
- a solicitação de compra é enviada a companhia do primeiro trecho
- essa companhia recebe o trecho e trava a função ára ninguem poder comprar
- os trechos que são de outra companhia precisam ser enviados para os servidores das outras companhias
- estrtura de compras precisará ter a companhia especificada
- em trechos de companhias diferentes, caso um não tenha vaga oq fazer?
-uma companhia "chama" a outra, uma companhia só é chamada se o trecho da companhia atual tiver vagas,
	caso contrário não há necessidade de "chamar" outra companhia, pois a compra não é possível.
-caso a companhia atual tenha vaga, a vaga só será decrementada após o retorno da companhia que foi chamada
	pois caso a outra companhia não tinha vagas em um dos trechos a compra já não pode ser realizada.
-caso uma companhia tenha mais de um trecho da rota, verificar as vagas de todos os trechos, pois depois o mutex estará fechado
	e não será possível entrar na função
*/

func main(){
	
}