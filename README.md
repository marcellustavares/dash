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

| Iteration  | 100K      | 1M       | 10M     | 100M   | 1B      |
| ---------  | ----------| -------- | ------- |------- |-------- |
| 1st (naive)| 23.75ms   | 139.78ms | 1.20s   | 11.91s | 1m57.32 |

## Cpu Profile

#### Iteration 1
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