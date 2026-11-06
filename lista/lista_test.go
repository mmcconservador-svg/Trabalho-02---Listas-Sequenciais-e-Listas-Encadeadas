package lista

import (
	"errors"
	"math/rand"
	"testing"
)

type pessoa struct {
	Nome  string
	Idade int
}

// conteudo devolve todos os elementos da lista, em ordem.
func conteudo[T any](t *testing.T, l Lista[T]) []T {
	t.Helper()
	out := make([]T, 0, l.Tamanho())
	for i := 0; i < l.Tamanho(); i++ {
		v, err := l.Consultar(i)
		if err != nil {
			t.Fatalf("Consultar(%d) falhou: %v", i, err)
		}
		out = append(out, v)
	}
	return out
}

func iguais[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// suite roda todos os testes de comportamento para um tipo T e para as
// duas implementações. valor(i) gera o i-ésimo valor de teste.
func suite[T comparable](t *testing.T, valor func(i int) T) {
	fabricas := map[string]func() Lista[T]{
		"Sequencial": func() Lista[T] { return NovaSequencial[T]() },
		"Encadeada":  func() Lista[T] { return NovaEncadeada[T]() },
	}

	for nome, nova := range fabricas {
		t.Run(nome+"/InserirInicio", func(t *testing.T) {
			l := nova()
			for i := 1; i <= 3; i++ {
				l.InserirInicio(valor(i))
			}
			if got, want := conteudo(t, l), []T{valor(3), valor(2), valor(1)}; !iguais(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})

		t.Run(nome+"/InserirMeio", func(t *testing.T) {
			l := nova()
			l.InserirMeio(valor(0)) // lista vazia: posição 0
			l.Remover(0)
			for i := 1; i <= 4; i++ {
				l.Adicionar(valor(i))
			}
			l.InserirMeio(valor(99)) // tamanho 4 -> posição 2
			want := []T{valor(1), valor(2), valor(99), valor(3), valor(4)}
			if got := conteudo(t, l); !iguais(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
			l.InserirMeio(valor(77)) // tamanho 5 -> posição 2
			want = []T{valor(1), valor(2), valor(77), valor(99), valor(3), valor(4)}
			if got := conteudo(t, l); !iguais(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})

		t.Run(nome+"/Remover", func(t *testing.T) {
			l := nova()
			for i := 1; i <= 5; i++ {
				l.Adicionar(valor(i))
			}
			if v, err := l.Remover(0); err != nil || v != valor(1) {
				t.Fatalf("Remover(0) = %v, %v", v, err)
			}
			if v, err := l.Remover(l.Tamanho() - 1); err != nil || v != valor(5) {
				t.Fatalf("Remover(último) = %v, %v", v, err)
			}
			if v, err := l.Remover(1); err != nil || v != valor(3) {
				t.Fatalf("Remover(1) = %v, %v", v, err)
			}
			l.Adicionar(valor(7)) // garante que o "final" da lista continua coerente
			want := []T{valor(2), valor(4), valor(7)}
			if got := conteudo(t, l); !iguais(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})

		t.Run(nome+"/EsvaziarEReusar", func(t *testing.T) {
			l := nova()
			l.Adicionar(valor(1))
			l.Remover(0)
			if l.Tamanho() != 0 {
				t.Fatalf("tamanho = %d, want 0", l.Tamanho())
			}
			l.Adicionar(valor(2))
			l.InserirInicio(valor(1))
			if got, want := conteudo(t, l), []T{valor(1), valor(2)}; !iguais(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		})

		t.Run(nome+"/PosicoesInvalidas", func(t *testing.T) {
			l := nova()
			if _, err := l.Consultar(0); !errors.Is(err, ErrPosicaoInvalida) {
				t.Errorf("Consultar em lista vazia: err = %v", err)
			}
			if _, err := l.Remover(0); !errors.Is(err, ErrPosicaoInvalida) {
				t.Errorf("Remover em lista vazia: err = %v", err)
			}
			l.Adicionar(valor(1))
			for _, pos := range []int{-1, 1, 10} {
				if _, err := l.Consultar(pos); !errors.Is(err, ErrPosicaoInvalida) {
					t.Errorf("Consultar(%d): err = %v", pos, err)
				}
				if _, err := l.Remover(pos); !errors.Is(err, ErrPosicaoInvalida) {
					t.Errorf("Remover(%d): err = %v", pos, err)
				}
			}
			for _, pos := range []int{-1, 2, 10} {
				if err := l.InserirEm(pos, valor(2)); !errors.Is(err, ErrPosicaoInvalida) {
					t.Errorf("InserirEm(%d): err = %v", pos, err)
				}
			}
			if l.Tamanho() != 1 {
				t.Errorf("operações inválidas alteraram o tamanho: %d", l.Tamanho())
			}
		})

		// Compara a lista com um slice simples usado como modelo de
		// referência, após uma sequência aleatória de operações.
		t.Run(nome+"/ModeloDeReferencia", func(t *testing.T) {
			rng := rand.New(rand.NewSource(1))
			l := nova()
			var ref []T
			for passo := 0; passo < 3000; passo++ {
				v := valor(passo)
				switch op := rng.Intn(5); {
				case op == 0:
					l.InserirInicio(v)
					ref = append([]T{v}, ref...)
				case op == 1:
					l.InserirMeio(v)
					m := len(ref) / 2
					ref = append(ref[:m], append([]T{v}, ref[m:]...)...)
				case op == 2:
					l.Adicionar(v)
					ref = append(ref, v)
				case op == 3 && len(ref) > 0:
					p := rng.Intn(len(ref))
					got, err := l.Remover(p)
					if err != nil || got != ref[p] {
						t.Fatalf("passo %d: Remover(%d) = %v, %v; want %v", passo, p, got, err, ref[p])
					}
					ref = append(ref[:p], ref[p+1:]...)
				case op == 4 && len(ref) > 0:
					p := rng.Intn(len(ref))
					got, err := l.Consultar(p)
					if err != nil || got != ref[p] {
						t.Fatalf("passo %d: Consultar(%d) = %v, %v; want %v", passo, p, got, err, ref[p])
					}
				}
				if l.Tamanho() != len(ref) {
					t.Fatalf("passo %d: tamanho = %d, want %d", passo, l.Tamanho(), len(ref))
				}
			}
			if got := conteudo(t, l); !iguais(got, ref) {
				t.Fatalf("conteúdo final difere do modelo de referência")
			}
		})
	}
}

func TestInt(t *testing.T) {
	suite(t, func(i int) int { return i })
}

func TestString(t *testing.T) {
	suite(t, func(i int) string { return "item-" + string(rune('a'+i%26)) + string(rune('A'+i/26%26)) })
}

func TestStruct(t *testing.T) {
	suite(t, func(i int) pessoa { return pessoa{Nome: "p" + string(rune('a'+i%26)), Idade: i} })
}
