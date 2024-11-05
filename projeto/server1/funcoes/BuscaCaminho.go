package funcoes

import (
	"container/heap"
	"sort"
	"fmt"
)

type Caminho struct {
	Trechos []Trecho
	Custo   int
}

// EncontrarTodosCaminhos encontra todos os caminhos possíveis do ponto de origem ao destino,
// retornando uma lista de listas de Trecho, ordenada pelo custo total (do menor para o maior).
func EncontrarTodosCaminhos(rotas map[string][]Trecho, origem, destino string) [][]Trecho {
	fmt.Println("Entrou em EncontrarTodosCaminhos")
	pq := &PriorityQueue{}
	heap.Init(pq)

	// Inicializa o heap com a origem
	heap.Push(pq, &Item{
		value:      Caminho{Trechos: []Trecho{}, Custo: 0},
		prioridade: 0,
		cidade:     origem,
	})

	var caminhosEncontrados []Caminho

	for pq.Len() > 0 {
		atual := heap.Pop(pq).(*Item)
		cidadeAtual := atual.cidade
		caminhoAtual := atual.value.Trechos

		// Verifica se o caminho atual já passou por essa cidade para evitar ciclos
		visitados := make(map[string]bool)
		for _, trecho := range caminhoAtual {
			visitados[trecho.Origem] = true
		}

		// Se alcançar o destino, armazena o caminho encontrado
		if cidadeAtual == destino {
			caminhosEncontrados = append(caminhosEncontrados, Caminho{
				Trechos: caminhoAtual,
				Custo:   atual.value.Custo,
			})
			continue
		}

		// Verifica todas as rotas a partir da cidade atual
		for _, trecho := range rotas[cidadeAtual] {
			cidadeAdjacente := trecho.Destino

			// Evita revisitar cidades já visitadas neste caminho
			if visitados[cidadeAdjacente] {
				continue
			}

			novoCusto := atual.value.Custo + trecho.Peso

			// Cria um novo caminho incluindo o trecho atual
			novoCaminho := append([]Trecho{}, caminhoAtual...)
			novoCaminho = append(novoCaminho, trecho)

			// Adiciona o novo caminho na fila de prioridade
			heap.Push(pq, &Item{
				value:      Caminho{Trechos: novoCaminho, Custo: novoCusto},
				prioridade: novoCusto,
				cidade:     cidadeAdjacente,
			})
		}
	}

	// Ordena os caminhos encontrados pelo custo total
	sort.SliceStable(caminhosEncontrados, func(i, j int) bool {
		return caminhosEncontrados[i].Custo < caminhosEncontrados[j].Custo
	})

	// Extrai apenas as listas de trechos para o retorno
	var resultado [][]Trecho
	for _, caminho := range caminhosEncontrados {
		resultado = append(resultado, caminho.Trechos)
	}
	return resultado
}


// Item representa um item na fila de prioridades
type Item struct {
	value      Caminho
	prioridade int
	cidade     string
	index      int
}

// PriorityQueue implementa a interface heap para ser usada como uma min-heap
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].prioridade < pq[j].prioridade
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}
