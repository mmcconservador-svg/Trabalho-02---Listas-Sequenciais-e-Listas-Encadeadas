package lista

// no é um nó da lista encadeada: guarda o valor e a referência (ponteiro)
// para o próximo nó.
type no[T any] struct {
	valor T
	prox  *no[T]
}

// Encadeada é uma lista simplesmente encadeada. Cada elemento vive em um
// nó alocado separadamente; os nós são ligados por ponteiros.
//
// Inserir no início é O(1). Consultar, inserir no meio e remover em uma
// posição qualquer são O(n), pois é preciso percorrer a lista desde a
// cabeça até chegar à posição desejada.
type Encadeada[T any] struct {
	cabeca *no[T]
	cauda  *no[T] // usada só por Adicionar, para montar listas em O(1)
	tam    int
}

// NovaEncadeada cria uma lista encadeada vazia.
func NovaEncadeada[T any]() *Encadeada[T] {
	return &Encadeada[T]{}
}

// noEm devolve o nó da posição pos. Pré-condição: 0 <= pos < tam.
func (l *Encadeada[T]) noEm(pos int) *no[T] {
	atual := l.cabeca
	for i := 0; i < pos; i++ {
		atual = atual.prox
	}
	return atual
}

// InserirInicio insere v na posição 0 em tempo constante.
func (l *Encadeada[T]) InserirInicio(v T) {
	_ = l.InserirEm(0, v)
}

// InserirMeio insere v na posição Tamanho()/2.
func (l *Encadeada[T]) InserirMeio(v T) {
	_ = l.InserirEm(l.tam/2, v)
}

// InserirEm insere v na posição pos.
func (l *Encadeada[T]) InserirEm(pos int, v T) error {
	if pos < 0 || pos > l.tam {
		return ErrPosicaoInvalida
	}
	novo := &no[T]{valor: v}
	switch {
	case pos == 0:
		novo.prox = l.cabeca
		l.cabeca = novo
		if l.tam == 0 {
			l.cauda = novo
		}
	case pos == l.tam:
		l.cauda.prox = novo
		l.cauda = novo
	default:
		anterior := l.noEm(pos - 1)
		novo.prox = anterior.prox
		anterior.prox = novo
	}
	l.tam++
	return nil
}

// Adicionar insere v ao final da lista em tempo constante.
func (l *Encadeada[T]) Adicionar(v T) {
	_ = l.InserirEm(l.tam, v)
}

// Consultar percorre a lista até a posição pos e devolve o elemento.
func (l *Encadeada[T]) Consultar(pos int) (T, error) {
	if pos < 0 || pos >= l.tam {
		var zero T
		return zero, ErrPosicaoInvalida
	}
	return l.noEm(pos).valor, nil
}

// Remover retira o elemento da posição pos.
func (l *Encadeada[T]) Remover(pos int) (T, error) {
	var zero T
	if pos < 0 || pos >= l.tam {
		return zero, ErrPosicaoInvalida
	}
	var removido *no[T]
	if pos == 0 {
		removido = l.cabeca
		l.cabeca = removido.prox
		if l.cabeca == nil {
			l.cauda = nil
		}
	} else {
		anterior := l.noEm(pos - 1)
		removido = anterior.prox
		anterior.prox = removido.prox
		if removido == l.cauda {
			l.cauda = anterior
		}
	}
	l.tam--
	return removido.valor, nil
}

// Tamanho devolve a quantidade de elementos.
func (l *Encadeada[T]) Tamanho() int {
	return l.tam
}
