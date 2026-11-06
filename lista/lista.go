// Package lista implementa duas estruturas de dados lineares genéricas:
// uma lista sequencial (baseada em vetor) e uma lista simplesmente
// encadeada. Ambas satisfazem a interface Lista[T] e funcionam com
// qualquer tipo de dados.
package lista

import "errors"

// ErrPosicaoInvalida é devolvido quando a posição pedida está fora dos
// limites da lista.
var ErrPosicaoInvalida = errors.New("lista: posição inválida")

// Lista é o contrato comum às duas estruturas. As posições começam em 0.
type Lista[T any] interface {
	// InserirInicio insere v na posição 0.
	InserirInicio(v T)

	// InserirMeio insere v na posição Tamanho()/2.
	InserirMeio(v T)

	// InserirEm insere v na posição pos (0 <= pos <= Tamanho()).
	InserirEm(pos int, v T) error

	// Adicionar insere v ao final da lista. Não faz parte das operações
	// pedidas no enunciado; existe apenas para montar rapidamente as
	// listas usadas nos experimentos.
	Adicionar(v T)

	// Consultar devolve o elemento da posição pos (0 <= pos < Tamanho()).
	Consultar(pos int) (T, error)

	// Remover retira e devolve o elemento da posição pos.
	Remover(pos int) (T, error)

	// Tamanho devolve a quantidade de elementos armazenados.
	Tamanho() int
}
