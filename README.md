# spinx
A golang fastcgi  proxy client.

# 新计划

1. 自己实现http协议处理
2. 更好的使用缓存的代理功能(参考Nginx)

# 2018-06-24

通过看源码

老的spinx 是 接收到一个完整的http协议 parse生成一个request 然后还有一系列的判断

监听仍然是socket 地址和端口

但是假如是转发的话，我可以在解析完协议的第一步立刻进行转发还不要去操作handler(request.go)

自己处理 go channel，问题可能会产生在锁的位置

old:  client->go server->received->handler->proxy->get response and rebuild->response

new:  client->go server->received and proxy->response(直接返回proxy的结果)


# contact

qq 285753421

email liwenlong0922@163.com
    
wechat
<img src="https://github.com/lwl1989/spinx/blob/master/Wechat.jpeg" alt="contact me with wechat" width="600" />





