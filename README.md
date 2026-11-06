# Trabalho 02 — Listas Sequenciais e Listas Encadeadas

Implementação, em Go, de uma lista sequencial (baseada em vetor) e de uma lista
simplesmente encadeada, ambas genéricas, e experimentos que comparam consulta,
inserção (início e meio) e remoção com 1.000, 10.000 e 100.000 elementos.

**Repositório:** https://github.com/mmcconservador-svg/Trabalho-02---Listas-Sequenciais-e-Listas-Encadeadas

**Grupo:** Ana Beatriz, Anthony Lucas e Miguel Monteiro

## Estrutura

```
lista/                 pacote com a interface Lista[T] e as duas implementações
  lista.go             interface comum e erro ErrPosicaoInvalida
  sequencial.go        lista sequencial (slice)
  encadeada.go         lista simplesmente encadeada
  lista_test.go        testes unitários (int, string e struct)
cmd/experimentos/      programa que cronometra as operações
analise/grafico.py     gera o gráfico a partir dos CSV
resultados/            saída dos experimentos (tempos.csv, memoria.csv, tabela.md, grafico_tempos.png)
```

## Como executar

Requer Go 1.21 ou superior (o código foi testado com Go 1.24).

```bash
go test ./...                    # testes unitários
go run ./cmd/experimentos        # experimentos (N = 1.000, 10.000 e 100.000)
python3 analise/grafico.py       # gráfico (requer matplotlib)
```

Opções do programa de experimentos:

```bash
go run ./cmd/experimentos -ns 1000,5000 -tentativas 3 -saida resultados
```

## Operações

| Operação | Método |
|---|---|
| Inserir no início | `InserirInicio(v)` |
| Inserir no meio (posição `Tamanho()/2`) | `InserirMeio(v)` |
| Consultar pela posição | `Consultar(pos)` |
| Remover de uma posição | `Remover(pos)` |
| Tamanho | `Tamanho()` |

Há ainda `InserirEm(pos, v)` e `Adicionar(v)` (insere no fim), usados para montar
as listas dos experimentos.

## Metodologia (resumo)

- Cada medição é a **mediana de 7 tentativas** do tempo médio por operação.
- Consulta: 10.000 consultas em posições aleatórias (mesma semente para as duas estruturas).
- Inserção e remoção: 2.000 operações em lotes de até N/10; o tamanho da lista é restaurado
  fora do cronômetro entre os lotes.
- Memória: bytes de heap por elemento depois de montar a lista e rodar o coletor de lixo.
- Os tempos absolutos dependem da máquina; refaça os experimentos no computador do grupo
  e registre o processador e a memória no relatório.
