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
