# Server Design Benchmarks
These servers are tested on a minimal DigitalOcean Droplet (1 CPU, 1GB RAM).
All of the servers are using a TCP listen backlog of 128.
Printing each request to stdout is disabled during the benchmark.

## Command
The tool [hey](https://github.com/rakyll/hey) is used to test each server's performance.
It performs 2000 requests spread across 50 concurrent connections.
```
hey -n 2000 -c 50 http://bloggulus.com
```

## Results
The results for each server design are represented here.
RPS stands for "Requests per Second" and latency values are in seconds.

| Server | RPS | Mean Latency | Mode Latency | Worst Latency |
| --- | --- | --- | --- | --- |
| NGINX | 672.36 | 0.07 | 0.11 | 0.64 |
| simple1.py | 54.23 | 0.84 | 0.92 | 8.10 |
| simple2.py | 52.77 | 0.84 | 0.91 | 7.85 |
| simple3.py | 45.71 | 0.78 | 1.64 | 15.88 |
| forking1.py | 8.70 | 5.60 | 5.14 | 10.14 |
| forking2.py | 81.83 | 0.54 | 0.87 | 7.60 |
| threading1.py | 78.71 | 0.52 | 0.87 | 7.52 |
| threading2.py | 79.06 | 0.51 | 0.88 | 7.73 |
| nonblock.py | 79.33 | 0.53 | 0.63 | 5.23 |
| async.py | 82.97 | 0.53 | 0.84 | 7.35 |
