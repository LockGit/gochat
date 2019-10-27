# gochat是一个使用纯go实现的im系统 

# architecture design 


# 组件
```
语言：golang
数据库：sqlite3 (可以根据实际业务场景替换成mysql或者其他数据库,在本项目中为方便演示，使用sqlite替代大型关系型数据库，仅存储了简单用户信息）
数据库ORM：gorm 
服务发现：etcd
rpc通讯：rpcx
队列:redis (方便演示使用redis，可以根据实际情况替换为kafka或者rabbitmq)
缓存:redis 
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

# 使用docker一键启动 
```
如果你觉得以上步骤过于繁琐，你可以使用以下docker镜像构建所有依赖环境并快速启动一个聊天室
1，docker pull 
2，sh run.sh 
3，访问
```


# 后续
```
gochat实现了简单聊天室功能，由于精力有限，你可以在此基础上使用自己的业务逻辑定制一些需求，并优化一些gochat中的代码
关于tcp粘包的处理，报文头部指定消息size,读取指定size为一包，不足size继续读取直到满足完整一包
```

```
cd db && sqlite3 gochat.sqlite3
.tables
消息id使用snowflakeId算法，此部分可以单独拆分成微服务，使其成为基础服务的一部分
该算法实现，qps理论最高409.6w，互联网上没有哪一家公司可以达到这么高的并发，除非遭受DDOS攻击 
```


```sqlite
create table user(
  `id` INTEGER PRIMARY KEY autoincrement , -- '用户id'
  `user_name` varchar(20) not null UNIQUE default '', -- '用户名'
  `password` char(40) not null default '', -- '密码'
  `create_time` timestamp NOT NULL DEFAULT current_timestamp -- '创建时间'
);
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



