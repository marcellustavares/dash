# Dash
[1BRC](https://github.com/gunnarmorling/1brc) in Golang.

## Hardware

- Chip: Apple M1 chip (8-core CPU).
- Memory: 8GB 
    - High-Performance Cores (4)
        - L1: 192 KB instruction cache and 128 KB data cache
        - L2: 12 MB shared across the four performance cores.
    - Energy-Efficient Cores (4)
        - L1: 128 KB instruction cache and 64 KB data cache.
        - L2: 4 MB shared across the four efficiency cores.

## Iterations

| Iteration     | 100K      | 1M       | 10M     | 100M   | 1B      |
| ---------     | ----------| -------- | ------- |------- |-------- |
| 1st (Naive)   | 23.75ms   | 139.78ms | 1.20s   | 11.91s | 1m57.32 |
| 2 (NewReader) | 22.67ms   | 156.48ms | 1.28s   | 12.64s | 2m7.24 |
| 3 (NewReader + Custom Buffer) |  21.24ms   | 116.00ms | 853.21ms   | 8.22s      | 1m21.80s |
| 4 (rm IndexByteString) |  18.48ms   | 109.47ms | 821.32ms   | 7.94s      | 1m19.88s |
| 5 (custom parseFloat) |  20.22ms   | 95.85ms |  622.61ms   | 6.01s      | 1m1.09s |

## Cpu Profile

#### Iteration 1 (Naive)
```
flat  flat%   sum%        cum   cum%
104.96s 85.26% 85.26%    104.96s 85.26%  syscall.rawsyscalln
6.28s  5.10% 90.37%      6.28s  5.10%  runtime.pthread_cond_wait
2.70s  2.19% 92.56%      2.70s  2.19%  runtime.usleep
2.67s  2.17% 94.73%      2.67s  2.17%  runtime.madvise
2.36s  1.92% 96.65%      2.36s  1.92%  runtime.pthread_cond_timedwait_relative_np
1.23s     1% 97.64%      1.23s     1%  runtime.pthread_kill
0.72s  0.58% 98.23%      0.72s  0.58%  runtime.kevent
0.08s 0.065% 98.29%      1.75s  1.42%  runtime.stealWork
0.05s 0.041% 98.33%    105.95s 86.07%  main.main
0.03s 0.024% 98.36%      8.73s  7.09%  runtime.findRunnable
```

#### Iteration 2 (NewReader)
```
flat  flat%   sum%        cum   cum%
112.88s 85.89% 85.89%    112.88s 85.89%  syscall.rawsyscalln
5.29s  4.03% 89.92%      5.29s  4.03%  runtime.pthread_cond_wait
4.48s  3.41% 93.33%      4.48s  3.41%  runtime.madvise
2.76s  2.10% 95.43%      2.76s  2.10%  runtime.usleep
1.71s  1.30% 96.73%      1.71s  1.30%  runtime.pthread_cond_timedwait_relative_np
1.26s  0.96% 97.69%      1.26s  0.96%  runtime.pthread_kill
0.12s 0.091% 97.78%      1.80s  1.37%  runtime.stealWork
0.04s  0.03% 97.81%    114.08s 86.81%  main.main
0.03s 0.023% 97.83%    112.97s 85.96%  bufio.(*Reader).ReadSlice
0.03s 0.023% 97.85%    113.15s 86.10%  bufio.(*Reader).ReadString
```

#### Iteration 3 (NewReader + 4MB Buffer)
```
flat  flat%   sum%        cum   cum%
7.74s 11.94% 11.94%      7.74s 11.94%  syscall.rawsyscalln
7.37s 11.37% 23.30%      7.37s 11.37%  runtime.pthread_cond_signal
6.15s  9.48% 32.79%      6.15s  9.48%  internal/bytealg.IndexByteString
5.59s  8.62% 41.41%      5.59s  8.62%  internal/strconv.readFloat
3.69s  5.69% 47.10%      6.86s 10.58%  runtime.mallocgcTiny
2.97s  4.58% 51.68%      2.97s  4.58%  aeshashbody
2.84s  4.38% 56.06%      4.93s  7.60%  runtime.mapassign_faststr
2.64s  4.07% 60.13%      2.64s  4.07%  runtime.pthread_cond_wait
2.43s  3.75% 63.88%      2.43s  3.75%  runtime.memmove
2.14s  3.30% 67.18%      2.14s  3.30%  runtime.madvise
```

#### Iteration 4 (rm IndexByteString)
```
flat  flat%   sum%        cum   cum%
10.43s 16.58% 16.58%     51.06s 81.19%  main.main
7.61s 12.10% 28.69%      7.61s 12.10%  syscall.rawsyscalln
6.77s 10.76% 39.45%      6.77s 10.76%  internal/strconv.readFloat
6.31s 10.03% 49.48%      6.32s 10.05%  runtime.pthread_cond_signal
3.17s  5.04% 54.52%      4.38s  6.96%  runtime.mallocgcTiny
3.06s  4.87% 59.39%      5.18s  8.24%  runtime.mapassign_faststr
2.48s  3.94% 63.33%      2.48s  3.94%  runtime.pthread_cond_wait
2.36s  3.75% 67.09%      2.36s  3.75%  aeshashbody
2.06s  3.28% 70.36%      9.16s 14.57%  runtime.slicebytetostring
1.99s  3.16% 73.53%      5.85s  9.30%  runtime.mapaccess2_faststr
```

#### Iteration 5 (custom parseFloat)
```
flat  flat%   sum%        cum   cum%
10.70s 21.18% 21.18%     39.42s 78.01%  main.main
6.74s 13.34% 34.51%      6.74s 13.34%  syscall.rawsyscalln
6.36s 12.59% 47.10%      6.36s 12.59%  runtime.pthread_cond_signal
3.17s  6.27% 53.37%      5.54s 10.96%  runtime.mapassign_faststr
2.50s  4.95% 58.32%      3.83s  7.58%  runtime.mallocgcTiny
2.19s  4.33% 62.66%      6.27s 12.41%  runtime.mapaccess2_faststr
2.14s  4.24% 66.89%      2.14s  4.24%  runtime.madvise
1.99s  3.94% 70.83%      1.99s  3.94%  aeshashbody
1.75s  3.46% 74.29%      1.75s  3.46%  runtime.memequal
1.73s  3.42% 77.72%      1.73s  3.42%  runtime.pthread_cond_wait
```