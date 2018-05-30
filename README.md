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

1. ssl support
2. channel add
3. keepalive support

# ab test

```
new version test(wordpress):

Server Software:        spinx
Server Hostname:        www.test.com
Server Port:            18000

Document Path:          /
Document Length:        53302 bytes

Concurrency Level:      100
Time taken for tests:   20.158 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      53507000 bytes
HTML transferred:       53302000 bytes
Requests per second:    49.61 [#/sec] (mean)
Time per request:       2015.799 [ms] (mean)
Time per request:       20.158 [ms] (mean, across all concurrent requests)
Transfer rate:          2592.17 [Kbytes/sec] received
==========================================================
Server Software:        nginx/1.13.9
Server Hostname:        www.test.com
Server Port:            18000

Document Path:          /
Document Length:        53302 bytes

Concurrency Level:      100
Time taken for tests:   20.277 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      53548000 bytes
HTML transferred:       53302000 bytes
Requests per second:    49.32 [#/sec] (mean)
Time per request:       2027.685 [ms] (mean)
Time per request:       20.277 [ms] (mean, across all concurrent requests)
Transfer rate:          2578.95 [Kbytes/sec] received
=============================================================
Performance has caught up with nginx.
```


# contact

qq 285753421

email liwenlong0922@163.com
    
wechat
<img src="https://github.com/lwl1989/spinx/blob/master/Wechat.jpeg" alt="contact me with wechat" width="600" />





