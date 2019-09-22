# gochat 是一个使用纯go实现的im系统 

# architecture design 

# 依赖
```
服务发现：etcd
rpc通讯：rpcx
```

# 按照以下顺序启动各层
```
0，编译
go build -o gochat.bin -tags=etcd main.go

1，启动logic层
./gochat.bin -module logic

2，启动connect层
./gochat.bin -module connect

3，启动task层
./gochat.bin -module task

4，启动api层
./gochat.bin -module api 

5，启动一个站点，开始聊天室
./gochat.bin -module site
```

# 后续
```markdown
gochat实现了简单基本聊天室功能，由于作者精力有限，你可以在此基础上使用自己的业务逻辑定制一些需求，并优化一些gochat中的代码
```

```markdown
connect层
    bucket (bucket,减少锁竞争)
        room (房间)
            channel (会话)
    setOnline--->redis

task层
    getMsg--->redis
    sendMsg--->rpc call connect

ui层
    chat:
        用户注册  logic
            用户名,密码   sqllite
        用户登录  logic   sqllite
            login,return sessionId
```



