| Operação | N | Sequencial | Encadeada | Encadeada / Sequencial |
|---|---:|---:|---:|---:|
| Consulta por posição | 1000 | 2.9 ns | 753.4 ns | 263.78x |
| Consulta por posição | 10000 | 2.5 ns | 7.71 µs | 3131.25x |
| Consulta por posição | 100000 | 13.9 ns | 82.74 µs | 5953.24x |
| Inserção no início | 1000 | 82.3 ns | 26.8 ns | 0.33x |
| Inserção no início | 10000 | 1.62 µs | 23.0 ns | 0.01x |
| Inserção no início | 100000 | 15.93 µs | 37.2 ns | 0.00x |
| Inserção no meio | 1000 | 45.1 ns | 892.9 ns | 19.78x |
| Inserção no meio | 10000 | 320.7 ns | 7.93 µs | 24.74x |
| Inserção no meio | 100000 | 8.09 µs | 76.91 µs | 9.51x |
| Remoção no meio | 1000 | 147.1 ns | 702.0 ns | 4.77x |
| Remoção no meio | 10000 | 944.1 ns | 7.62 µs | 8.07x |
| Remoção no meio | 100000 | 13.94 µs | 75.80 µs | 5.44x |
| Remoção no início | 1000 | 334.8 ns | 3.5 ns | 0.01x |
| Remoção no início | 10000 | 2.86 µs | 3.8 ns | 0.00x |
| Remoção no início | 100000 | 31.39 µs | 5.5 ns | 0.00x |

| Estrutura | N | Bytes por elemento (heap) |
|---|---:|---:|
| Sequencial | 1000 | 10.3 |
| Encadeada | 1000 | 16.0 |
| Sequencial | 10000 | 9.8 |
| Encadeada | 10000 | 16.0 |
| Sequencial | 100000 | 8.8 |
| Encadeada | 100000 | 16.0 |
