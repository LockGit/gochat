/**
 * Created by lock
 * Date: 2020/5/20
 */
package stickpackage

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gochat/config"
	"gochat/proto"
	"log"
	"net"
	"testing"
	"time"
)

func Test_TestStick(t *testing.T) {
	pack := &StickPackage{
		Version: VersionContent,
		Msg:     []byte(("now time:" + time.Now().Format("2006-01-02 15:04:05"))),
	}
	pack.Length = pack.GetPackageLength()

	buf := new(bytes.Buffer)
	//test package , BigEndian
	_ = pack.Pack(buf)
	_ = pack.Pack(buf)
	_ = pack.Pack(buf)
	_ = pack.Pack(buf)
	// scanner
	scanner := bufio.NewScanner(buf)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF && data[0] == 'v' {
			if len(data) > 4 {
				packSumLength := int16(0)
				_ = binary.Read(bytes.NewReader(data[2:4]), binary.BigEndian, &packSumLength)
				if int(packSumLength) <= len(data) {
					return int(packSumLength), data[:packSumLength], nil
				}
			}
		}
		return
	})

	scannedPack := new(StickPackage)
	for scanner.Scan() {
		err := scannedPack.Unpack(bytes.NewReader(scanner.Bytes()))
		if err != nil {
			log.Println(err.Error())
		}
		log.Println(scannedPack)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Invalid data")
		t.Fail()
	}
}

func Test_TcpClient(t *testing.T) {
	//1,建立tcp链接
	//2,send msg to tcp conn
	//3,receive msg from tcp conn
	roomId := 1                                                      //@todo default roomId
	authToken := "1kHYNlHaQTjGd0BWuECkw80ZAIquoU30f0gFPxqpEhQ="      //@todo need you modify
	fromUserId := 3                                                  //@todo need you modify
	tcpAddrRemote, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:7001") //@todo default connect address
	conn, err := net.DialTCP("tcp", nil, tcpAddrRemote)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		panic("conn err:" + err.Error())
	}

	go func() {
		//读取服务端广播的信息
		onMessageReceive := func(conn *net.TCPConn) error {
			for {
				scannerPackage := bufio.NewScanner(conn)
				scannerPackage.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
					if !atEOF && data[0] == 'v' {
						if len(data) > TcpHeaderLength {
							packSumLength := int16(0)
							_ = binary.Read(bytes.NewReader(data[LengthStartIndex:LengthStopIndex]), binary.BigEndian, &packSumLength)
							if int(packSumLength) <= len(data) {
								return int(packSumLength), data[:packSumLength], nil
							}
						}
					}
					return
				})
				scanErrorTimes := 0
				for {
					if scanErrorTimes > 3 {
						break
					}
					scanErrorTimes++
					fmt.Println("start read tcp msg from conn...")
					for scannerPackage.Scan() {
						scannedPack := new(StickPackage)
						err := scannedPack.Unpack(bytes.NewReader(scannerPackage.Bytes()))
						if err != nil {
							log.Printf("unpack msg err:%s", err.Error())
							break
						}
						fmt.Println(fmt.Sprintf("read msg from tcp ok,version is:%s,length is:%d,msg is:%s", scannedPack.Version, scannedPack.Length, scannedPack.Msg))
					}
					if scannerPackage.Err() != nil {
						log.Printf("scannerPackage err:%s", err.Error())
						break
					}
				}
				return nil
			}
		}(conn)
		fmt.Println("onMessageReceive err is:", onMessageReceive)
	}()
	var i int
	for {
		if i == 0 {
			fmt.Println("build tcp heartbeat conn...")
			msg := &proto.SendTcp{
				Msg:          "build tcp heartbeat conn",
				FromUserId:   fromUserId,
				FromUserName: "Tcp heartbeat build",
				RoomId:       roomId,
				Op:           config.OpBuildTcpConn,
				AuthToken:    authToken, //todo 增加token验证，用于验证tcp部分
			}
			msgBytes, _ := json.Marshal(msg)
			//生成带房间号的的msg并pack write conn io
			pack := &StickPackage{
				Version: VersionContent,
				//Msg:     []byte(("now time:" + time.Now().Format("2006-01-02 15:04:05"))),
				Msg: msgBytes,
			}
			pack.Length = pack.GetPackageLength()
			//test package, BigEndian
			_ = pack.Pack(conn) //写入要发送的消息
		}
		fmt.Println("time wait , you can remove the code!")
		time.Sleep(10 * time.Second)
		msg := &proto.SendTcp{
			Msg:          "from tcp client,time is:" + time.Now().Format("2006-01-02 15:04:05"),
			FromUserId:   fromUserId,
			FromUserName: "I am Tcp msg",
			RoomId:       roomId,
			Op:           config.OpRoomSend,
			AuthToken:    authToken, //todo 增加token验证，用于验证tcp部分
		}
		msgBytes, _ := json.Marshal(msg)
		//生成带房间号的的msg并pack write conn io
		pack := &StickPackage{
			Version: VersionContent,
			//Msg:     []byte(("now time:" + time.Now().Format("2006-01-02 15:04:05"))),
			Msg: msgBytes,
		}
		pack.Length = pack.GetPackageLength()
		//test package, BigEndian
		_ = pack.Pack(conn) //写入要发送的消息
		i++
		fmt.Println(fmt.Sprintf("第%d次send msg to tcp server", i))
	}
}
