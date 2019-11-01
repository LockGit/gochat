[English version(英文版本)](readme.en.md)

### gochat是一个使用纯go实现的轻量级im系统 
```
gochat为纯go实现的即时通讯系统,支持私信消息与房间广播消息,各层之间通过rpc通讯,支持水平扩展。
使用redis作为消息存储与投递的载体,相对kafka操作起来更加方便快捷,所以十分轻量。
各层之间基于etcd服务发现,在扩容部署时将会方便很多。
由于go的交叉编译特性,编译后可以快速在各个平台上运行,gochat架构及目录结构清晰,
并且本项目还贴心的提供了docker一键构建所有环境依赖,安装起来十分便捷。
```

### 架构设计
 ![](https://github.com/LockGit/gochat/blob/master/architecture/gochat.png)

### 服务发现
![](https://github.com/LockGit/gochat/blob/master/architecture/gochat_discovery.png)


### 消息投递 
![](https://github.com/LockGit/gochat/blob/master/architecture/single_send.png)
```
消息发送必须在登录状态下,如上图,用户A向用户B发送了一条消息。那么经历了如下历程：
1,用户A调用api层接口登录系统,登录成功后与connect层认证保持长链接,
并rpc call logic层记录用户A在connect层登录的serverId,默认加入房间号1

2,用户B调用api层接口登录系统,登录成功后与connect层认证保持长链接,
并rpc call logic层记录用户B在connect层登录的serverId,默认加入房间号1

3,用户A调用api层接口发送消息,api层rpc call logic层发送消息方法,
logic层将该消息push到队列中等待task层消费

4,task层订阅logic层发送到队列中的消息后,
根据消息内容(userId,roomId,serverId)可以定位用户B目前在connect层那一个serverId上保持长链接,
进一步rpc call connect层方法

5,connect层被task层rpc call之后,会将该消息投递到相关房间内，
再进一步投递给在该房间内的用户B,完成一次完整的消息会话通讯

6,如果是私信消息,那么task层会根据userId定位connect层对应的serverId,
直接rpc call 该serverId的connect层，完成私信消息投递。

学习其他im系统，为了减少锁的竞争,会在connect层会划分bucket:
大致会像如下结构：
connect层:
    bucket: (bucket,减少锁竞争)
        room: (房间)
            channel: (用户会话)
```

### 聊天室预览
![](https://github.com/LockGit/gochat/blob/master/architecture/gochat.gif)

### 目录结构
```
➜  gochat git:(master) ✗ tree -L 1
.
├── api                  # api接口层,提供rest api服务,主要rpc call logic层服务,可水平扩展 
├── architecture         # 架构图资源图片文件夹
├── bin                  # golang编译后的二进制文件,不参与提交至git仓库
├── config               # 配置文件
├── connect              # 链接层,该层用于handler住大量用户长连接,实时消息推送,除此之外,主要rpc call logic层服务,可水平扩展
├── db                   # 数据库链接初始化,内部gochat.sqlite3为方便演示时使用的sqlite数据库,实际场景可替换为其他关系型数据库
├── docker               # 用于快速构建docker环境
├── go.mod               # go包管理mod文件
├── go.sum               # go包管理sum文件
├── logic                # 逻辑层,该层主要接收connect层与api层的rpc请求。如果是消息,会push到队列,最终被task层消费,可水平扩展
├── main.go              # gochat程序唯一入口文件
├── proto                # golang rpc proto 文件
├── readme.md            # 中文说明文档
├── readme.en.md         # 英文说明文档
├── reload.sh            # 编译gochat并执行supervisorctl restart重启相关进程
├── run.sh               # 快速构建一个docker容器,启动demo
├── site                 # 站点层,该层为纯静态页面,会http请求api层,可水平扩展
├── task                 # 任务层,该层消费队列中的数据,并且rpc call connect层执行消息的发送,可水平扩展
├── tools                # 工具方法
└── vendor               # vendor包
```

### 相关组件
```
语言：golang

数据库：sqlite3 
可以根据实际业务场景替换成mysql或者其他数据库,
在本项目中为方便演示,使用sqlite替代大型关系型数据库,仅存储了简单用户信息

数据库ORM：gorm 

服务发现：etcd

rpc通讯：rpcx

队列:redis 
方便演示使用redis,可以根据实际情况替换为kafka或者rabbitmq

缓存:redis 
存储用户session,以及相关计数器,聊天室房间信息等

消息id:
snowflakeId算法,此部分可以单独拆分成微服务,使其成为基础服务的一部分。
该id发号器qps理论最高409.6w/s,互联网上没有哪一家公司可以达到这么高的并发,除非遭受DDOS攻击
```

### 数据库及表结构
```
在本demo中,为了足够轻量方便演示,数据库使用了sqlite3,基于gorm,所以这个是可以替换的。
如果需要替换其他关系型数据库,仅仅只需要修改相关db驱动即可。

相关表结构：
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

### 安装
```
在启动各层之前,请确保已经启动了etcd与redis服务以及以上数据库表,
并确保7000, 7070, 8080端口没有被占用。
然后按照以下顺序启动各层,如果要扩容connect层,
请确保connect层配置中各个serverId不一样!

0,编译
go build -o gochat.bin -tags=etcd main.go

1,启动logic层
./gochat.bin -module logic

2,启动connect层
./gochat.bin -module connect

3,启动task层
./gochat.bin -module task

4,启动api层
./gochat.bin -module api 

5,启动一个站点,开始聊天室
./gochat.bin -module site
```

### 使用docker一键启动 
```
如果你觉得以上步骤过于繁琐,你可以使用以下docker镜像构建所有依赖环境并快速启动一个聊天室

你可以使用我推到docker hub上的镜像:
默认镜像中已经创建了几个测试用户：
用户名,密码
demo  111111
test  111111
admin 111111
1,docker pull lockgit/gochat:latest
2,sh run.sh dev
3,访问 http://127.0.0.1:8080 开启聊天室


如果你想自己构建一个镜像,那么只需要build docker文件下的Dockerfile
docker build -f docker/Dockerfile . -t lockgit/gochat
然后执行sh run.sh dev即可 

如果你要部署在个人vps上,记得修改site/js/common.js中socketUrl与apiUrl的地址为你的vps的ip地址,
并确保vps上没有针对相关端口的防火墙限制。
```

### 这里提供一个在线聊天demo站点：
<a href="http://155.138.194.113:8080" target="_blank">http://155.138.194.113:8080</a>
```
用以上用户名密码登录即可,也可以自己注册一个。
可用不同的账号在不同的浏览器登录,如果是Chrome,可用以隐身模式多启动几个浏览器分别以不同账号登录,
然后,体验不同账号之间聊天以及消息投递效果。
```


### 后续
```
gochat实现了简单聊天室功能,由于精力有限,
你可以在此基础上使用自己的业务逻辑定制一些需求,并优掉一些gochat中的代码。
看懂跟实现一遍是不一样的,照着例子做完之后应该会理解的更深刻一些,
另关于tcp粘包的处理,通用做法在报文头部指定消息size,
读取指定size为一包,不足size继续读取直到满足完整一包,本demo中暂未处理。
后面会根据情况优化和改进gochat中的代码与设计。
```



