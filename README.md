# spinx v2

A golang fastcgi  proxy client.

Contributors:
* LI WEN LONG
* Emre

# quick start

    go get github.com/lwl1989/spinx
    
    cd $gopath/src/github.com/lwl1989/spinx
    
    go build -o spinx main.go
    
    ./spinx

# change log 

2018-10-17

    1. v2 init,rewrite core
    2. self parse protocol


# config.json demo

```
{
    "server": {
        "port": "8081",
        "log": "/tmp/spinx.log",
        "keep_alive_timeout": 3,
        "gzip_level": 1,
        "cache": {
            "len": 10240,
            "expire": "24h"
        }
    },

    "vhosts": [
        {
            "name": "www.test.com wwww.aaa.com",
            "port": "18000",
            "proxy":"127.0.0.1:9000",
            "documentRoot": "/www/web/wordpress",
            "tryFiles": "/index.php?$uri",
            "index":  "index.php index.html"
        },
        {
            "name": "www.phpmyadmin.com",
            "port": "18000",
            "proxy":"127.0.0.1:9000",
            "documentRoot": "/www/web/phpmyadmin",
            "tryFiles": "/index.php?$uri"

        }
    ]
}
```

# The next Step

1. change to other  faster http core
2. keepalive timer  set
3. proxy add

# ab test

```
new version(hello world! PHP 7.2.3):
Server Software:        spinx
Server Hostname:        www.test.com
Server Port:            18000

Document Path:          /
Document Length:        12 bytes

Concurrency Level:      100
Time taken for tests:   7.151 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      15300000 bytes
HTML transferred:       1200000 bytes
Requests per second:    13983.36 [#/sec] (mean)
Time per request:       7.151 [ms] (mean)
Time per request:       0.072 [ms] (mean, across all concurrent requests)
Transfer rate:          2089.31 [Kbytes/sec] received

==========================================================
Server Software:        nginx/1.13.9
Server Hostname:        www.word.com
Server Port:            80

Document Path:          /
Document Length:        12 bytes

Concurrency Level:      100
Time taken for tests:   6.773 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      17400000 bytes
HTML transferred:       1200000 bytes
Requests per second:    14764.09 [#/sec] (mean)
Time per request:       6.773 [ms] (mean)
Time per request:       0.068 [ms] (mean, across all concurrent requests)
Transfer rate:          2508.74 [Kbytes/sec] received
=============================================================
```


# contact

qq 285753421

email liwenlong0922@163.com
    
wechat
<img src="https://github.com/lwl1989/spinx/blob/master/Wechat.jpeg" alt="contact me with wechat" width="600" />





