[English version(英文版本)](readme.md)

### gochat是一个使用纯go实现的轻量级im系统 
```
gochat为纯go实现,各层之间均支持水平扩展,队列使用redis,相对kafaka操作起来更加方便快捷,所以十分轻量。
消息id使用snowflakeId算法，此部分可以单独拆分成微服务，使其成为基础服务的一部分。
该算法实现，id发号器qps理论最高409.6w/s，互联网上没有哪一家公司可以达到这么高的并发，除非遭受DDOS攻击。

双向链表
```

### 架构设计
 ![](https://github.com/LockGit/gochat/blob/master/architecture/gochat.gng)

#### 服务发现
![](https://github.com/LockGit/gochat/blob/master/architecture/gochat_discovery.gng)

### 目录结构及说明
```
➜  gochat git:(master) ✗ tree -L 1
.
├── api                  # api接口层，提供rest api服务，主要rpc call logic层服务，可水平扩展 
├── architecture         # 架构图资源图片文件夹
├── bin                  # golang编译后的二进制文件，不参提交至git仓库
├── config               # 配置文件
├── connect              # 链接层，该层用于handler住大量用户长连接，实时消息推送，除此之外，主要rpc call logic层服务，可水平扩展
├── db                   # 数据库链接初始化，内部gochat.sqlite3为方便演示时使用的sqlite数据库，实际场景可替换为其他关系型数据库
├── docker               # 用于快速构建docker环境
├── go.mod               # go包管理mod文件
├── go.sum               # go包管理sum文件
├── logic                # 逻辑层，该层主要接收connect层与api层的rpc请求。如果是消息，会push到队列，最终被task层消费，可水平扩展
├── main.go              # gochat程序唯一入口文件
├── proto                # golang rpc proto 文件
├── readme.md            # 英文说明文档
├── readme.zh-cn.md      # 中文说明文档
├── reload.sh            # 编译gochat并执行supervisorctl restart重启相关进程
├── run.sh               # 快速构建一个docker容器，启动demo
├── site                 # 站点层，该层为纯静态页面，会http请求api层，可水平扩展
├── task                 # 任务层，该层消费队列中的数据，并且rpc call connect层执行消息的发送，可水平扩展
├── tools                # 工具方法
└── vendor               # vendor包
```

### 单发消息


### 广播消息


### 相关组件
```
语言：golang
数据库：sqlite3 (可以根据实际业务场景替换成mysql或者其他数据库,在本项目中为方便演示，使用sqlite替代大型关系型数据库，仅存储了简单用户信息）
数据库ORM：gorm 
服务发现：etcd
rpc通讯：rpcx
队列:redis (方便演示使用redis，可以根据实际情况替换为kafka或者rabbitmq)
缓存:redis (用户session,以及相关计数器,聊天室房间信息等)
```

### 数据库表结构
```
cd db && sqlite3 gochat.sqlite3
.tables
```
```sqlite
create table user(
  `id` INTEGER PRIMARY KEY autoincrement , -- '用户id'
  `user_name` varchar(20) not null UNIQUE default '', -- '用户名'
  `password` char(40) not null default '', -- '密码'
  `create_time` timestamp NOT NULL DEFAULT current_timestamp -- '创建时间'
);
```

# 按照以下顺序启动各层
```
在启动各层之前，请确保已经启动了etcd与redis服务
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

### 使用docker一键启动 
```
如果你觉得以上步骤过于繁琐，你可以使用以下docker镜像构建所有依赖环境并快速启动一个聊天室
1，docker pull 
2，sh run.sh 
3，访问
```


### 后续
```
gochat实现了简单聊天室功能，由于精力有限，你可以在此基础上使用自己的业务逻辑定制一些需求，并优化一些gochat中的代码
关于tcp粘包的处理，报文头部指定消息size,读取指定size为一包，不足size继续读取直到满足完整一包
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



