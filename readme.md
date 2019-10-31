[中文版本(Chinese version)](readme.zh-cn.md)

### Gochat is a lightweight im server implemented using pure go
```
gochat is an instant messaging system implemented by pure go. 
It supports private message and room broadcast messages. 
rpc communication is provided between layers to support horizontal expansion.
Using redis as a carrier for message storage and delivery, 
it is more convenient and faster to operate than kafaka.
so it is very lightweight. 
Based on the etcd service discovery between layers, 
it will be much more convenient when expanding and deploying.
Due to the cross-compilation feature of go, 
it can be run on various platforms quickly after compilation. 
The gochat architecture and directory structure are clear.
And the project also provides docker one-click to build all environment dependencies, 
which is very convenient to install.
```

### Architecture design
 ![](https://github.com/LockGit/gochat/blob/master/architecture/gochat.png)

### Service discovery
![](https://github.com/LockGit/gochat/blob/master/architecture/gochat_discovery.png)


### Message delivery
![](https://github.com/LockGit/gochat/blob/master/architecture/signle_send.png)
```
The message must be sent in the login state. As shown above, User A sends a message to User B. 
Then experienced the following process:
1. User A invokes the api layer interface to log in to the system. After the login succeeds, 
the user keeps a long link with the connect layer authentication, 
and the rpc call logic layer records the serverId that user A logs in at the connect layer. 
The default joins the room number 1.

2. User B calls the api layer interface to log in to the system. 
After the login succeeds, the user keeps a long link with the connect layer authentication, 
and the rpc call logic layer records the serverId of the user B login at the connect layer. 
The default joins the room number 1.

3. User A calls the api layer interface to send a message, 
and the api layer rpc call logic layer sends a message method. 
The logic layer pushes the message to the queue and waits for the task layer to consume.

4, After the task layer subscribes to the message sent by the logic layer to the queue, 
according to the message content (userId, roomId, serverId), 
the user B can be positioned to maintain a long link on the serverId of the connect layer, 
and the rpc call connect layer method is further

5, After the connect layer is rpc call by the task layer, 
the message will be delivered to the relevant room, 
and then further delivered to the user B in the room to complete a complete message session communication.

Learning other im systems, in order to reduce the lock competition, will be divided in the connect layer bucket:
It will roughly look like the following structure:
Connect layer:
     Bucket: (bucket, reduce lock competition)
         Room: (room)
             Channel: (user session)
```

### Chat room preview
![](https://github.com/LockGit/gochat/blob/master/architecture/gochat.gif)


### Directory Structure
```
➜ gochat git:(master) ✗ tree -L 1
.
├── api             # api interface layer, providing rest api service, main rpc call logic layer service, horizontal expansion
├──architecture     # architecture resource picture folder
├── bin             # golang compiled binary, do not submit to the git repository
├── config          # configuration file
├── connect         # link layer, this layer is used for handlers to live a large number of user long connections, real-time message push, in addition, the main rpc call logic layer service, can be horizontally extended
├── db              # database link initialization, internal gochat.sqlite3 for the convenience of the use of the sqlite database, the actual scene can be replaced with other relational database
├── docker          # for quickly building a docker environment
├── go.mod          # go package management mod file
├── go.sum          # go package management sum file
├── logic           # Logic layer, this layer mainly receives the rpc request of the connect layer and the api layer. If it is a message, it will be pushed to the queue, and finally consumed by the task layer, which can be horizontally expanded.
├── main.go         # gochat program only entry file
├── proto           # golang rpc proto file
├── readme.md       # English documentation
├── readme.zh-cn.md # Chinese documentation
├── reload.sh       # Compile gochat and execute supervisorctl restart to restart related processes
├── run.sh          # Quickly build a docker container and start the demo
├── site            # site layer, this layer is a pure static page, will http request api layer, can be extended horizontally
├── task            # task layer, the data in the consumption queue of this layer, and the rpc call connect layer performs message transmission, which can be extended horizontally
├── tools           #tools function
└── vendor          # vendor package
```

### Related components
```
Language: golang

Database: sqlite3
can be replaced by mysql or other database according to the actual business scenario. In this project, 
for the convenience of demonstration, use sqlite instead of large relational database, 
only store simple user information

Database ORM: gorm

Service discovery: etcd

Rpc communication: rpcx

Queue: redis 
convenient to use redis, can be replaced with kafka or rabbitmq according to the actual situation

Cache: redis 
store user session, and related counters, chat room information, etc.

Message id: 
The snowflakeId algorithm, 
which can be split into microservices separately, 
making it part of the underlying service. 
The id numbering device qps theory is up to 409.6w/s. 
No company on the Internet can achieve such high concurrency unless it suffers from DDOS attacks.
```

### Database and table structure
```
In this demo, in order to be lightweight and convenient to demonstrate, 
the database uses sqlite3, based on gorm, 
so this can be replaced. 
If you need to replace other relational databases, just modify the relevant db driver.

Related table structure:
cd db && sqlite3 gochat.sqlite3
.tables
```
```sqlite
create table user(
  `id` INTEGER PRIMARY KEY autoincrement , -- 'userId'
  `user_name` varchar(20) not null UNIQUE default '', -- 'userName'
  `password` char(40) not null default '', -- 'passWord'
  `create_time` timestamp NOT NULL DEFAULT current_timestamp -- 'createTime'
);
```

### Installation
```
Before starting each layer, 
make sure that the etcd and redis services and the above database tables have been started, 
and then start the layers in the following order. 

If you want to expand the connect layer, 
make sure that the serverId in the connect layer configuration is different!

0,Compile
go build -o gochat.bin -tags=etcd main.go

1,Start the logic layer
./gochat.bin -module logic

2,Start the connect layer
./gochat.bin -module connect

3,Start the task layer
./gochat.bin -module task

4,Start the api layer
./gochat.bin -module api 

5,Start a site and start a chat room
./gochat.bin -module site
```

### Start with a docker
```
If you feel that the above steps are too cumbersome, 
you can use the following docker image to build all the dependencies and quickly start a chat room.

You can use the image I pushed to the docker hub 
(there have been several test users created in the default image: lock, demo, test password: 111111)
1,docker pull lockgit/gochat:latest
2,sh run.sh dev
3,visit http://127.0.0.1:8080 to open the chat room


If you want to build an image yourself, 
you only need to build the Dockerfile under the docker file.
Docker build -f docker/Dockerfile . -t lockgit/gochat
then execute sh run.sh dev
```


### Follow-up
```
gochat implements a simple chat room function. Due to limited energy, 
you can use your own business logic to customize some requirements and optimize some code in gochat.
It is not the same as implementing it once. It should be understood more deeply after the example is done. 
For the handling of tcp sticky packets, the common practice is to specify the message size in the message header and read the specified size as a package,
Insufficient size continues to read until the complete package is satisfied. This demo has not been processed yet.
The code and design in gochat will be optimized and improved as needed.
```
