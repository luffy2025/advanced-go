```
shushu at ssdeMbp in ~
$ redis-benchmark -h localhost -p 6379 -q -t get,set -d 10
SET: 10397.17 requests per second, p50=4.543 msec
GET: 16444.66 requests per second, p50=2.743 msec

(base)
shushu at ssdeMbp in ~
$ redis-benchmark -h localhost -p 6379 -q -t get,set -d 20
SET: 10758.47 requests per second, p50=4.367 msec
GET: 14594.28 requests per second, p50=2.655 msec

(base)
shushu at ssdeMbp in ~
$ redis-benchmark -h localhost -p 6379 -q -t get,set -d 50
SET: 11154.49 requests per second, p50=4.191 msec
GET: 16963.53 requests per second, p50=2.687 msec

(base)
shushu at ssdeMbp in ~
$ redis-benchmark -h localhost -p 6379 -q -t get,set -d 100
SET: 11625.20 requests per second, p50=4.031 msec
GET: 17170.33 requests per second, p50=2.655 msec

(base)
shushu at ssdeMbp in ~
$ redis-benchmark -h localhost -p 6379 -q -t get,set -d 200
SET: 11171.94 requests per second, p50=4.199 msec
GET: 15928.64 requests per second, p50=2.887 msec

(base)
shushu at ssdeMbp in ~
$ redis-benchmark -h localhost -p 6379 -q -t get,set -d 1000
SET: 10241.70 requests per second, p50=4.559 msec
GET: 14863.26 requests per second, p50=3.023 msec

(base)
shushu at ssdeMbp in ~
$ redis-benchmark -h localhost -p 6379 -q -t get,set -d 5000
SET: 6318.72 requests per second, p50=7.111 msec
GET: 10748.07 requests per second, p50=4.207 msec

(base)
```
