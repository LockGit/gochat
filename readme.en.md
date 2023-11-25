# 🚀 gochat
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/LockGit/gochat/issues)
<img src="https://img.shields.io/badge/gochat-im-green">
<img src="https://img.shields.io/badge/documentation-yes-brightgreen.svg">
<img src="https://img.shields.io/github/license/LockGit/gochat">
<img src="https://img.shields.io/github/contributors/LockGit/gochat">
<img src="https://img.shields.io/github/last-commit/LockGit/gochat">
<img src="https://img.shields.io/github/issues/LockGit/gochat"> 
<img src="https://img.shields.io/github/forks/LockGit/gochat"> 
<img src="https://img.shields.io/github/stars/LockGit/gochat">
<img src="https://img.shields.io/docker/pulls/lockgit/gochat">
<img src="https://img.shields.io/github/repo-size/LockGit/gochat">
<img src="https://img.shields.io/github/followers/LockGit">

### [中文版本(Chinese version)](readme.md)

### gochat is a lightweight im server implemented using pure go
* gochat is an instant messaging system based on go. It supports private message messages and room broadcast messages. It communicates with each other through RPC and supports horizontal expansion.
* Support websocket and TCP access, and support websocket and TCP message interworking in the latest version.
* Based on etcd service discovery among different layers, it will be much more convenient to expand and deploy.
* Using redis as the carrier of message storage and delivery is very lightweight, and it can be replaced by heavier Kafka and rabbitmq in actual scenarios.
* Because of the clear structure of chago, it can run quickly on various platforms.
* This project provides docker with one click to build all environment dependencies, which is very convenient to install. (if it's an experience, we strongly recommend using docker to build)

### Renew
* August 6, 2022
  * Although UI is not very important in this system, it does look awkward. Typescript + react is used to paste a slightly normal front-end UI interface. UI project address: https://github.com/LockGit/gochat-ui
* May 8, 2022
  * At present, the version of golang has been upgraded to 1.18. If you do not use docker, choose to compile and install it yourself. Please ensure that your go version meets the requirements
  * The vendor is also packaged in the warehouse. In addition, the entire project size of the vendor is about 66m, which can be used after git clone. Due to irresistible network factors, it may take a long time
  * Some package dependent versions have been upgraded. It is not recommended to try to compile this project on a lower version of golang. Try to upgrade to 1.18+
  * Optimization: dynamically update the service IP address according to the corresponding kV change in the watch etcd to ensure that other layers can perceive the added / removed layers

### About Websocket && Tcp Support
```
The latest version supports websocket and TCP message interworking,(the test TCP port uses 70017002. If you want to use TCP, the firewall also needs to release the two ports.)
TCP message delivery and receiving test code in the pkg/stickpackage directory of this project: stickpackage_test.go Test in file_Tcpclient method
among stickpackage.go The file is mainly used for TCP unpacking and unpacking. You can trace where the pack and unpack methods are called.
The main reason is that TCP is based on the streaming protocol of layer 4 rather than the application layer protocol, so there is this process.
If it is Android, IOS client to link, then the corresponding is to use a familiar language to implement the TCP unpacking and unpacking process. The code in the example is demo implemented by golang.
go test -v -count=1 *.go -test.run Test_TcpClient

For test TCP message delivery:
Just change the @ todo part of this method to the correct content when you tested it.
AuthToken is the authentication token for TCP link. This token is the user ID. you can see the JSON return in the API return after login on the web, or you can directly find the token of a user in redis.
When the server receives the authToken, it will add the corresponding user's TCP link session to the session bidirectional list of the corresponding room.

After adding TCP support, if there are both users linking through websocket and users linking through TCP in a room, the general process of message delivery is as follows:
1. First, the user logs in, carries the token and room number to connect to the websocket and TCP server, and establishes a long link session in the link layer,
For example, user a is in room 1 and linked through websocket, and user B is also in room 1 through TCP.

2. Users send messages to the room. At present, they directly put the messages into the corresponding message queue through HTTP,
In fact, this place can also be sent through websocket and TCP, but the final messages are to be queued.

3. The message in the queue is consumed by the task layer, and the RPC is broadcast to the corresponding rooms of the websocket link layer and the TCP link layer,
After getting the message, the link layer delivers the message to the corresponding remote user (equivalent to traversing the user link session bidirectional list maintained in the room)
```
![](./architecture/gochat_tcp.gif)

### Architecture design
 ![](./architecture/gochat.png)

### Service discovery
![](./architecture/gochat_discovery.png)


### Message delivery
![](./architecture/single_send.png)
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

6, if it is a private message, 
then the task layer will locate the serverId corresponding to the connect layer according to the userId, 
directly rpc call the connect layer of the serverId, and complete the private message delivery.

Learning other im systems, in order to reduce the lock competition, will be divided in the connect layer bucket:
It will roughly look like the following structure:
Connect layer:
     Bucket: (bucket, reduce lock competition)
         Room: (room)
             Channel: (user session)
```
### timing
![](./architecture/timing.png)

### Chat room preview
Use typescript + react to paste a slightly normal front-end UI. UI project address: https://github.com/LockGit/gochat-ui
![](./architecture/gochat-new.png)

### Internal structure of user long connect session
![](./architecture/session.png)

### Dividing buckets to reduce lock competition
![](./architecture/lock_competition.png)

### Making hash distribution more uniform based on cityhash
![](./architecture/hash.png)

### Preview of the internal structure of a single server in the connect layer
![](./architecture/gochat_room.png)

### Simple pressure measurement
![](./architecture/pressure_test.png)

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
├── readme.md       # Chinese documentation
├── readme.en.md    # English documentation
├── reload.sh       # Compile gochat and execute supervisorctl restart to restart related processes
├── run.sh          # Quickly build a docker container and start the demo
├── site            # site layer, this layer is a pure static page, will http request api layer, can be extended horizontally
├── task            # task layer, the data in the consumption queue of this layer, and the rpc call connect layer performs message transmission, which can be extended horizontally
├── tools           # tools function
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
Please make sure that the 7000, 7070, 8080 ports are not occupied. (the test TCP port uses 70017002. If you want to use TCP, the firewall also needs to release the two ports.)
make sure that the etcd and redis services and the above database tables have been started, 
and then start the layers in the following order. 

If you want to expand the connect layer, 
make sure that the serverId in the connect layer configuration is different!

0,Compile
go build -o gochat.bin -tags=etcd main.go

1,Start the logic layer
./gochat.bin -module logic

2,Start the connect layer
./gochat.bin -module connect_tcp
or  
./gochat.bin -module connect_websocket

3,Start the task layer
./gochat.bin -module task

4,Start the api layer
./gochat.bin -module api 

5,Start a site and start a chat room
./gochat.bin -module site
```

### Start with a docker

#### Use existing image
```
If you feel that the above steps are too cumbersome, 
you can use the following docker image to build all the dependencies and quickly start a chat room.

You can use the image I pushed to the docker hub 
there have been several test users created in the default image.
username  password
demo        111111
test        111111
admin       111111
1,docker pull lockgit/gochat:1.18 (version 1.18 is currently used)
2,git clone git@github.com:LockGit/gochat.git
3,cd gochat && sh run.sh dev 127.0.0.1 (This step requires a certain time to compile each module and wait patiently. Some systems may not have sh. if an error is reported during execution, change sh run.sh dev to ./ run.sh dev 127.0.0.1 for execution)
4,visit http://127.0.0.1:8080/login; enter username/password to access the chat room
5,visit http://127.0.0.1:8080 to open the chat room
```


#### Building a new image
For other platforms, such as arm or m1 chips
```
If you want to build an image yourself, 
you only need to build the Dockerfile under the docker file.
Docker build -f docker/Dockerfile . -t lockgit/gochat
then execute:
1,git clone git@github.com:LockGit/gochat.git
2,cd gochat && sh run.sh dev 127.0.0.1 (This step requires a certain time to compile each module and wait patiently. Some systems may not have sh. if an error is reported during execution, change sh run.sh dev to ./ run.sh dev 127.0.0.1 for execution)
3,visit http://127.0.0.1:8080/login; enter username/password to access the chat room
4,visit http://127.0.0.1:8080 to open the chat room

If you want to deploy on personal vps, you can execute: sh run.sh dev vps_ip
And make sure there are no firewall restrictions on the relevant ports on the vps.
```

### Here is an online chat demo site:
<a href="http://45.77.108.245:8080" target="_blank">http://45.77.108.245:8080</a>
```
Log in with the above user name and password, or register one by yourself.
if the provide url can't visit one day，please use the docker above to visit it.
You can log in to different browsers with different accounts. 
If it is Chrome, you can use several browsers in incognito mode to log in with different accounts.
Then, experience the chat between different accounts and the effect of message delivery.
(Note: TCP channel is not enabled in this demo)
```
If you want to deploy a similar private IM chat service to the public network experience, 
you can also register an account on vultr and purchase a VPS host to deploy the IM service experience:
<a href="https://www.vultr.com/?ref=8981750-8H" target="_blank">https://www.vultr.com/?ref=8981750-8H</a>

### Follow-up
```
gochat implements a simple chat room function. Due to limited energy, 
you can use your own business logic to customize some requirements and optimize some code in gochat.
Please contact the relevant issue for any questions in use, 
and follow up and optimize the relevant code according to the actual situation.
```

### Backer and Sponsor
>jetbrains

<a href="https://www.jetbrains.com/?from=LockGit/gochat" target="_blank">
<img src="https://github.com/LockGit/gochat/blob/master/architecture/jetbrains-gold-reseller.svg" width="100px" height="100px">
</a>

```
JetBrains has an open source sponsorship program that allows you to apply for the JetBrains family bucket license for free through the open source project
Jetbranins officials will ask for the proposal to add their brand logo to the warehouse when giving the license,
But all this is the principle of users' willingness
```

### Support gochat
```
Your support is also a kind of affirmation, if you find it helpful, 
you can also scan the WeChat QR code below to support the author to continue to optimize and improve gochat
```
![](./architecture/support.png)

### License
gochat is licensed under the MIT License. 
