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
| 6 (fewer allocs) |  13.30ms   | 72.70ms |  531.96ms   | 4.71s      | 46.71s |
| 7 (no map) |  11.19ms   | 54.85ms |  343.34ms   | 3.11s      | 30.92s |
| 8 (single pass) |  10.94ms   | 58.59ms |  337.77ms   | 3.10s      | 30.38s |
| 9 (multithreading) |  20.77ms   | 39.37ms |  119.70ms   | 773.57ms      |  6.81s |
| 10 (parseFloat optimized) |  18.94ms   | 40.13ms |  111.61ms   | 687.14ms      |  5.75s |
| 11 (simplify and remove mod op) |  18.14ms   | 35.22ms |  90.43ms   | 597.45ms      |  5.10s |

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

#### Iteration 6 (fewer allocs)
```
flat  flat%   sum%        cum   cum%
10260ms 27.24% 27.24%    29280ms 77.73%  main.main
4790ms 12.72% 39.95%     4790ms 12.72%  runtime.pthread_cond_signal
4560ms 12.11% 52.06%     4560ms 12.11%  syscall.rawsyscalln
2310ms  6.13% 58.19%     3540ms  9.40%  runtime.mallocgcTiny
2000ms  5.31% 63.50%     2000ms  5.31%  runtime.pthread_cond_wait
1910ms  5.07% 68.57%     5560ms 14.76%  runtime.mapaccess2_faststr
1580ms  4.19% 72.76%     1580ms  4.19%  main.parseFloat (inline)
1410ms  3.74% 76.51%     5040ms 13.38%  runtime.mallocgc
1090ms  2.89% 79.40%     1090ms  2.89%  aeshashbody
910ms  2.42% 81.82%      910ms  2.42%  runtime.memequal
```

#### Iteration 7 (no map)
```
flat  flat%   sum%        cum   cum%
8.64s 31.92% 31.92%      8.64s 31.92%  syscall.rawsyscalln
8.52s 31.47% 63.39%     22.31s 82.42%  main.main
1.80s  6.65% 70.04%      1.80s  6.65%  runtime.pthread_cond_signal
1.78s  6.58% 76.62%      1.78s  6.58%  runtime.pthread_cond_wait
1.63s  6.02% 82.64%      1.64s  6.06%  main.HashBytes (inline)
1.28s  4.73% 87.37%      4.92s 18.18%  main.processRow
1.20s  4.43% 91.80%      1.20s  4.43%  main.parseFloat (inline)
0.80s  2.96% 94.75%      0.80s  2.96%  runtime.memequal
0.37s  1.37% 96.12%      0.67s  2.48%  runtime.scanObject
0.22s  0.81% 96.93%      0.22s  0.81%  runtime.usleep
```

#### Iteration 8 (single pass hash and float parser)

```
flat  flat%   sum%        cum   cum%
12.72s 46.53% 46.53%     23.42s 85.66%  main.main
8.69s 31.78% 78.31%      8.69s 31.78%  syscall.rawsyscalln
1.52s  5.56% 83.87%      1.52s  5.56%  runtime.pthread_cond_wait
1.16s  4.24% 88.11%      1.93s  7.06%  main.processRow
1.15s  4.21% 92.32%      1.15s  4.21%  runtime.pthread_cond_signal
0.77s  2.82% 95.14%      0.77s  2.82%  runtime.memequal
0.29s  1.06% 96.20%      0.61s  2.23%  runtime.scanObject
0.23s  0.84% 97.04%      0.32s  1.17%  runtime.typePointers.next
0.23s  0.84% 97.88%      0.23s  0.84%  runtime.usleep
0.14s  0.51% 98.39%      0.14s  0.51%  runtime.pthread_kill
```

#### Iteration 9 (multithreading)

```
flat  flat%   sum%        cum   cum%
25.62s 66.42% 66.42%     35.37s 91.70%  main.processChunk
5.60s 14.52% 80.94%      7.44s 19.29%  main.processRow
1.78s  4.61% 85.56%      1.78s  4.61%  syscall.rawsyscalln
1.74s  4.51% 90.07%      1.74s  4.51%  runtime.memequal
0.68s  1.76% 91.83%      0.68s  1.76%  runtime.pthread_cond_wait
0.61s  1.58% 93.41%      0.61s  1.58%  runtime.usleep
0.46s  1.19% 94.61%      1.09s  2.83%  runtime.scanObject
0.39s  1.01% 95.62%      0.58s  1.50%  runtime.typePointers.next
0.30s  0.78% 96.40%      0.30s  0.78%  runtime.pthread_cond_signal
0.30s  0.78% 97.17%      0.30s  0.78%  runtime.pthread_kill
```

#### Iteration 10 (parserFloat optimized)

```
flat  flat%   sum%        cum   cum%
15440ms 48.16% 48.16%    28760ms 89.71%  main.processChunk
4800ms 14.97% 63.13%     4850ms 15.13%  main.parseTemperature
4730ms 14.75% 77.89%     6220ms 19.40%  main.processRow
1620ms  5.05% 82.94%     1620ms  5.05%  syscall.rawsyscalln
1470ms  4.59% 87.52%     1470ms  4.59%  runtime.memequal
630ms  1.97% 89.49%     1180ms  3.68%  runtime.scanObject
630ms  1.97% 91.45%      630ms  1.97%  runtime.usleep
610ms  1.90% 93.36%      610ms  1.90%  runtime.pthread_cond_wait
450ms  1.40% 94.76%      450ms  1.40%  runtime.pthread_kill
400ms  1.25% 96.01%      400ms  1.25%  runtime.memmove
```

#### Iteration 11 (remove leftover copy and mod op)

```
flat  flat%   sum%        cum   cum%
14.58s 53.33% 53.33%     26.58s 97.22%  main.processChunk
5.01s 18.32% 71.65%      5.04s 18.43%  main.parseTemperature
3.65s 13.35% 85.00%      5.18s 18.95%  main.processRow
1.72s  6.29% 91.29%      1.72s  6.29%  syscall.rawsyscalln
1.49s  5.45% 96.74%      1.49s  5.45%  runtime.memequal
0.44s  1.61% 98.35%      0.44s  1.61%  runtime.pthread_cond_wait
0.16s  0.59% 98.94%      0.16s  0.59%  runtime.pthread_cond_signal
0.03s  0.11% 99.05%      1.52s  5.56%  bytes.Equal (inline)
0.01s 0.037% 99.09%      0.55s  2.01%  runtime.findRunnable
0.01s 0.037% 99.12%      0.50s  1.83%  runtime.goschedImpl
```