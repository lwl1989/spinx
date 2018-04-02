# spinx
A golang fastcgi  proxy client.

# quick start

    go get github.com/lwl1989/spinx
    cd $gopath/src/github.com/lwl1989/spinx
    go build
    ./spinx

    Please input your config path(default /usr/etc/spinx/server.json):


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

# contact

    qq 285753421
    email liwenlong0922@163.com

<img src="https://github.com/lwl1989/spinx/blob/master/Wechat.jpeg" alt="contact me with wechat" width="600" />





