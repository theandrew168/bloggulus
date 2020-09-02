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
RPS stands for "Requests per Second" and latency is measured in seconds.

| Server | RPS | Mean Latency | Mode Latency | Worst Latency | % NGINX (RPS) |
| --- | --- | --- | --- | --- | --- |
| NGINX | 672.36 | 0.07 | 0.11 | 0.64 | 100.00 |
| Waitress | 449.07 | 0.10 | 0.09 | 0.42 | 66.79 |
| Gunicorn[sync] | 46.31 | 0.99 | 0.95 | 8.23 | 6.89 |
| Gunicorn[gevent] | 338.36 | 0.14 | 0.09 | 0.43 | 50.32 |
| simple1.py | 54.23 | 0.84 | 0.92 | 8.10 | 8.07 |
| simple2.py | 52.77 | 0.84 | 0.91 | 7.85 | 7.85 |
| simple3.py | 45.71 | 0.78 | 1.64 | 15.88 | 6.80 |
| forking1.py | 8.70 | 5.60 | 5.14 | 10.14 | 1.29 |
| forking2.py | 81.83 | 0.54 | 0.87 | 7.60 | 12.17 |
| threading1.py | 78.71 | 0.52 | 0.87 | 7.52 | 11.71 |
| threading2.py | 79.06 | 0.51 | 0.88 | 7.73 | 11.76 |
| processpool.py | 48.52 | 0.71 | 1.01 | 8.81 | 7.22 |
| threadpool.py | 70.04 | 0.62 | 0.87 | 7.49 | 10.42 |
| nonblocking.py | 79.33 | 0.53 | 0.63 | 5.23 | 11.80 |
| async.py | 82.97 | 0.53 | 0.84 | 7.35 | 12.34 |
