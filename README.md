# spinx
A golang fastcgi  proxy client.

# quick start

    go get github.com/lwl1989/spinx
    
    cd $gopath/src/github.com/lwl1989/spinx
    
    go build -o spinx main.go
    
#### install    
    sudo ./spinx  -c=config_path install
#### remove    
    sudo ./spinx  remove
#### start
    sudo ./spinx start
    or
    ./spinx -d=false -c=config_path
#### stop
    sudo ./spinx stop    

# change log 

2018-5-14
  
    1. add handler coroutine
    2. add logger(handler.SetLogger(log))

2018-05-30
    
    1. add daemon
    2. add cmd set config and list help
    3. add to system service

2018-06-03

    1. keep-alive support

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





