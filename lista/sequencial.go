package lista

// Sequencial é uma lista sequencial: os elementos ficam lado a lado em um
// vetor (slice do Go) contíguo na memória.
//
// Consultar é O(1), pois o endereço do elemento é calculado a partir da
// posição. Inserir e remover no início ou no meio são O(n), pois exigem
// deslocar os elementos que vêm depois.
type Sequencial[T any] struct {
	dados []T
}

// NovaSequencial cria uma lista sequencial vazia.
func NovaSequencial[T any]() *Sequencial[T] {
	return &Sequencial[T]{}
}

// InserirInicio insere v na posição 0, deslocando todos os elementos.
func (l *Sequencial[T]) InserirInicio(v T) {
	_ = l.InserirEm(0, v)
}

// InserirMeio insere v na posição Tamanho()/2.
func (l *Sequencial[T]) InserirMeio(v T) {
	_ = l.InserirEm(len(l.dados)/2, v)
}

// InserirEm insere v na posição pos, deslocando uma casa para a direita
// os elementos de pos em diante.
func (l *Sequencial[T]) InserirEm(pos int, v T) error {
	if pos < 0 || pos > len(l.dados) {
		return ErrPosicaoInvalida
	}
	var zero T
	l.dados = append(l.dados, zero) // abre espaço (pode realocar o vetor)
	copy(l.dados[pos+1:], l.dados[pos:])
	l.dados[pos] = v
	return nil
}

// Adicionar insere v ao final da lista.
func (l *Sequencial[T]) Adicionar(v T) {
	l.dados = append(l.dados, v)
}

// Consultar devolve o elemento da posição pos em tempo constante.
func (l *Sequencial[T]) Consultar(pos int) (T, error) {
	if pos < 0 || pos >= len(l.dados) {
		var zero T
		return zero, ErrPosicaoInvalida
	}
	return l.dados[pos], nil
}

// Remover retira o elemento da posição pos, deslocando uma casa para a
// esquerda os elementos seguintes.
func (l *Sequencial[T]) Remover(pos int) (T, error) {
	var zero T
	if pos < 0 || pos >= len(l.dados) {
		return zero, ErrPosicaoInvalida
	}
	v := l.dados[pos]
	copy(l.dados[pos:], l.dados[pos+1:])
	l.dados[len(l.dados)-1] = zero // não deixa referência pendurada para o GC
	l.dados = l.dados[:len(l.dados)-1]
	return v, nil
}

// Tamanho devolve a quantidade de elementos.
func (l *Sequencial[T]) Tamanho() int {
	return len(l.dados)
}
