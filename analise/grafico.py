"""Gera o gráfico comparativo a partir de resultados/tempos.csv.

Uso (na raiz do repositório):
    python3 analise/grafico.py

Saída: resultados/grafico_tempos.png — um painel por operação, tempo médio por
operação (escala logarítmica) em função do número de elementos.
"""
import csv
import math
import os
from collections import OrderedDict

import matplotlib

matplotlib.use("Agg")
import matplotlib.pyplot as plt  # noqa: E402

AZUL, LARANJA = "#2a6fd6", "#e8683a"
COR = {"Sequencial": AZUL, "Encadeada": LARANJA}
RAIZ = os.path.join(os.path.dirname(__file__), "..")
ENTRADA = os.path.join(RAIZ, "resultados", "tempos.csv")
SAIDA = os.path.join(RAIZ, "resultados", "grafico_tempos.png")


def formatar(ns):
    texto = f"{ns:.1f} ns" if ns < 1000 else f"{ns / 1000:.1f} µs"
    return texto.replace(".", ",")


def rotulo_eixo(valor, _pos):
    return f"{valor:g} ns" if valor < 1000 else f"{valor / 1000:g} µs"


def milhar(n):
    return f"{int(n):,}".replace(",", ".")


def main():
    dados = OrderedDict()  # operacao -> estrutura -> [(n, ns)]
    with open(ENTRADA, newline="", encoding="utf-8") as f:
        for r in csv.DictReader(f):
            dados.setdefault(r["operacao"], OrderedDict()).setdefault(
                r["estrutura"], []
            ).append((int(r["n"]), float(r["mediana_ns"])))

    fig, eixos = plt.subplots(2, 3, figsize=(13, 7.2), sharey=True)
    eixos = eixos.flatten()
    nmax = max(n for ops in dados.values() for pts in ops.values() for n, _ in pts)

    for ax, (operacao, series) in zip(eixos, dados.items()):
        for estrutura, pts in series.items():
            pts.sort()
            ax.plot(*zip(*pts), marker="o", lw=2, color=COR[estrutura], label=estrutura)
            n, ns = pts[-1]
            ax.annotate(formatar(ns), (n, ns), textcoords="offset points",
                        xytext=(0, 8), ha="center", fontsize=8)
        finais = {e: [v for n, v in p if n == nmax][0] for e, p in series.items()}
        rapida = min(finais, key=finais.get)
        fator = max(finais.values()) / min(finais.values())
        ax.set_title(f"{operacao}\n{rapida}: {fator:,.1f}× mais rápida".replace(",", "."),
                     fontsize=10, loc="left")
        ax.set_xscale("log")
        ax.set_yscale("log")
        ax.set_xticks([1000, 10000, 100000])
        ax.set_xticklabels([milhar(n) for n in (1000, 10000, 100000)])
        ax.yaxis.set_major_formatter(matplotlib.ticker.FuncFormatter(rotulo_eixo))
        ax.grid(axis="y", color="#e5e5e5")
        ax.spines[["top", "right"]].set_visible(False)

    for ax in eixos[len(dados):]:
        ax.axis("off")
    eixos[0].legend(frameon=False, loc="upper left", fontsize=9)
    fig.supxlabel("Número de elementos (N)")
    fig.supylabel("Tempo médio por operação (escala logarítmica)")
    fig.suptitle("Lista sequencial × lista encadeada: o vetor vence a consulta e o meio; "
                 "a lista vence o início", fontsize=12)
    fig.tight_layout()
    fig.savefig(SAIDA, dpi=150)
    print("gerado:", os.path.normpath(SAIDA))


if __name__ == "__main__":
    main()
