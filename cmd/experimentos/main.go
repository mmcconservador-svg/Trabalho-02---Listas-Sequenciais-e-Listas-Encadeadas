// Programa de experimentos do Trabalho 02: compara a lista sequencial e a
// lista encadeada em consulta, inserção (início e meio) e remoção.
//
// Uso:
//
//	go run ./cmd/experimentos                      # N = 1.000, 10.000 e 100.000
//	go run ./cmd/experimentos -ns 1000,5000 -tentativas 3
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"trabalho02/lista"
)

// sumidouro impede que o compilador elimine as consultas como código morto.
var sumidouro int

type fabrica struct {
	nome string
	nova func() lista.Lista[int]
}

var fabricas = []fabrica{
	{"Sequencial", func() lista.Lista[int] { return lista.NovaSequencial[int]() }},
	{"Encadeada", func() lista.Lista[int] { return lista.NovaEncadeada[int]() }},
}

// experimento descreve uma operação cronometrada. As operações que alteram
// a lista são executadas em lotes; entre um lote e outro, restaura repõe o
// tamanho original (fora do cronômetro) para que o N medido seja o N real.
type experimento struct {
	nome     string
	total    int             // operações cronometradas por tentativa
	lote     func(n int) int // tamanho máximo de cada lote
	executa  func(l lista.Lista[int], k int, posicoes []int)
	restaura func(l lista.Lista[int], k int)
}

var experimentos = []experimento{
	{
		nome:  "Consulta por posição",
		total: 10000,
		lote:  func(int) int { return 10000 },
		executa: func(l lista.Lista[int], k int, pos []int) {
			s := 0
			for i := 0; i < k; i++ {
				v, _ := l.Consultar(pos[i])
				s += v
			}
			sumidouro += s
		},
	},
	{
		nome:  "Inserção no início",
		total: 2000,
		lote:  func(n int) int { return max(100, n/10) },
		executa: func(l lista.Lista[int], k int, _ []int) {
			for i := 0; i < k; i++ {
				l.InserirInicio(i)
			}
		},
		restaura: func(l lista.Lista[int], k int) {
			for i := 0; i < k; i++ {
				l.Remover(0)
			}
		},
	},
	{
		nome:  "Inserção no meio",
		total: 2000,
		lote:  func(n int) int { return max(100, n/10) },
		executa: func(l lista.Lista[int], k int, _ []int) {
			for i := 0; i < k; i++ {
				l.InserirMeio(i)
			}
		},
		restaura: func(l lista.Lista[int], k int) {
			for i := 0; i < k; i++ {
				l.Remover(l.Tamanho() / 2)
			}
		},
	},
	{
		nome:  "Remoção no meio",
		total: 2000,
		lote:  func(n int) int { return max(100, n/10) },
		executa: func(l lista.Lista[int], k int, _ []int) {
			for i := 0; i < k; i++ {
				l.Remover(l.Tamanho() / 2)
			}
		},
		restaura: func(l lista.Lista[int], k int) {
			for i := 0; i < k; i++ {
				l.InserirMeio(i)
			}
		},
	},
	{
		nome:  "Remoção no início",
		total: 2000,
		lote:  func(n int) int { return max(100, n/10) },
		executa: func(l lista.Lista[int], k int, _ []int) {
			for i := 0; i < k; i++ {
				l.Remover(0)
			}
		},
		restaura: func(l lista.Lista[int], k int) {
			for i := 0; i < k; i++ {
				l.InserirInicio(i)
			}
		},
	},
}

type resultado struct {
	estrutura, operacao string
	n                   int
	mediana, min, max   float64 // ns por operação
}

// medirUma monta uma lista nova com n elementos e devolve o tempo médio
// (ns) por operação do experimento e.
func medirUma(f fabrica, n int, e experimento, posicoes []int) float64 {
	l := f.nova()
	for i := 0; i < n; i++ {
		l.Adicionar(i)
	}
	lote := min(e.lote(n), e.total)
	lotes := e.total / lote

	runtime.GC()
	var decorrido time.Duration
	for b := 0; b < lotes; b++ {
		t0 := time.Now()
		e.executa(l, lote, posicoes)
		decorrido += time.Since(t0)
		if e.restaura != nil {
			e.restaura(l, lote)
		}
	}
	return float64(decorrido.Nanoseconds()) / float64(lotes*lote)
}

// medirMemoria devolve os bytes de heap usados por elemento após montar
// uma lista com n inteiros.
func medirMemoria(f fabrica, n int) float64 {
	var antes, depois runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&antes)
	l := f.nova()
	for i := 0; i < n; i++ {
		l.Adicionar(i)
	}
	runtime.GC()
	runtime.ReadMemStats(&depois)
	runtime.KeepAlive(l)
	return float64(int64(depois.HeapAlloc)-int64(antes.HeapAlloc)) / float64(n)
}

func resumir(amostras []float64) (mediana, minimo, maximo float64) {
	s := append([]float64(nil), amostras...)
	sort.Float64s(s)
	return s[len(s)/2], s[0], s[len(s)-1]
}

func formatarNs(ns float64) string {
	switch {
	case ns < 1e3:
		return fmt.Sprintf("%.1f ns", ns)
	case ns < 1e6:
		return fmt.Sprintf("%.2f µs", ns/1e3)
	default:
		return fmt.Sprintf("%.2f ms", ns/1e6)
	}
}

func lerNs(s string) ([]int, error) {
	var ns []int
	for _, p := range strings.Split(s, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil || n < 1 {
			return nil, fmt.Errorf("valor inválido em -ns: %q", p)
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func main() {
	nsFlag := flag.String("ns", "1000,10000,100000", "tamanhos de lista, separados por vírgula")
	tentativas := flag.Int("tentativas", 7, "repetições de cada medição (usa-se a mediana)")
	saida := flag.String("saida", "resultados", "diretório de saída (CSV e tabela)")
	flag.Parse()

	ns, err := lerNs(*nsFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := os.MkdirAll(*saida, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Go %s, %s/%s, %d CPUs, %d tentativas por medição\n\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), *tentativas)

	rng := rand.New(rand.NewSource(42))
	var resultados []resultado

	for _, n := range ns {
		for _, e := range experimentos {
			// Mesmas posições aleatórias para as duas estruturas.
			posicoes := make([]int, e.total)
			for i := range posicoes {
				posicoes[i] = rng.Intn(n)
			}
			for _, f := range fabricas {
				amostras := make([]float64, *tentativas)
				for t := range amostras {
					amostras[t] = medirUma(f, n, e, posicoes)
				}
				med, mn, mx := resumir(amostras)
				resultados = append(resultados, resultado{f.nome, e.nome, n, med, mn, mx})
				fmt.Printf("N=%-7d %-22s %-11s %s\n", n, e.nome, f.nome, formatarNs(med))
			}
		}
	}

	// Memória.
	type memoria struct {
		estrutura string
		n         int
		porElem   float64
	}
	var memorias []memoria
	for _, n := range ns {
		for _, f := range fabricas {
			amostras := make([]float64, *tentativas)
			for t := range amostras {
				amostras[t] = medirMemoria(f, n)
			}
			med, _, _ := resumir(amostras)
			memorias = append(memorias, memoria{f.nome, n, med})
		}
	}

	// CSV de tempos.
	if err := escreverCSV(filepath.Join(*saida, "tempos.csv"),
		[]string{"estrutura", "n", "operacao", "mediana_ns", "min_ns", "max_ns"},
		func(w *csv.Writer) {
			for _, r := range resultados {
				w.Write([]string{r.estrutura, strconv.Itoa(r.n), r.operacao,
					fmt.Sprintf("%.2f", r.mediana), fmt.Sprintf("%.2f", r.min), fmt.Sprintf("%.2f", r.max)})
			}
		}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// CSV de memória.
	if err := escreverCSV(filepath.Join(*saida, "memoria.csv"),
		[]string{"estrutura", "n", "bytes_por_elemento"},
		func(w *csv.Writer) {
			for _, m := range memorias {
				w.Write([]string{m.estrutura, strconv.Itoa(m.n), fmt.Sprintf("%.2f", m.porElem)})
			}
		}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Tabela em Markdown.
	var sb strings.Builder
	sb.WriteString("| Operação | N | Sequencial | Encadeada | Encadeada / Sequencial |\n")
	sb.WriteString("|---|---:|---:|---:|---:|\n")
	for _, e := range experimentos {
		for _, n := range ns {
			var seq, enc float64
			for _, r := range resultados {
				if r.operacao == e.nome && r.n == n {
					if r.estrutura == "Sequencial" {
						seq = r.mediana
					} else {
						enc = r.mediana
					}
				}
			}
			fmt.Fprintf(&sb, "| %s | %d | %s | %s | %.2fx |\n", e.nome, n, formatarNs(seq), formatarNs(enc), enc/seq)
		}
	}
	sb.WriteString("\n| Estrutura | N | Bytes por elemento (heap) |\n|---|---:|---:|\n")
	for _, m := range memorias {
		fmt.Fprintf(&sb, "| %s | %d | %.1f |\n", m.estrutura, m.n, m.porElem)
	}
	if err := os.WriteFile(filepath.Join(*saida, "tabela.md"), []byte(sb.String()), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("\n%s\nArquivos gerados em %s/ (tempos.csv, memoria.csv, tabela.md)\n", sb.String(), *saida)
}

func escreverCSV(caminho string, cabecalho []string, linhas func(w *csv.Writer)) error {
	f, err := os.Create(caminho)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Write(cabecalho)
	linhas(w)
	w.Flush()
	return w.Error()
}
