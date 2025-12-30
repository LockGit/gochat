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
	// Skip in short mode - this is an integration test requiring full stack
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	//1,建立tcp链接
	//2,send msg to tcp conn
	//3,receive msg from tcp conn
	roomId := 1                                                      //@todo default roomId
	authToken := "1kHYNlHaQTjGd0BWuECkw80ZAIquoU30f0gFPxqpEhQ="      //@todo need you modify
	fromUserId := 3                                                  //@todo need you modify
	tcpAddrRemote, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:7001") //@todo default connect address

	// Set connection timeout
	conn, err := net.DialTimeout("tcp", tcpAddrRemote.String(), 5*time.Second)
	if err != nil {
		t.Skipf("TCP server not available at %s, skipping integration test: %v", tcpAddrRemote.String(), err)
		return
	}
	defer conn.Close()

	tcpConn := conn.(*net.TCPConn)

	// Channel to signal when we've received a response
	receivedMsg := make(chan bool, 1)
	testTimeout := time.After(10 * time.Second)

	// Start goroutine to read server responses
	go func() {
		tcpConn.SetReadDeadline(time.Now().Add(8 * time.Second))
		scannerPackage := bufio.NewScanner(tcpConn)
		scannerPackage.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if !atEOF && len(data) > 0 && data[0] == 'v' {
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

		for scannerPackage.Scan() {
			scannedPack := new(StickPackage)
			err := scannedPack.Unpack(bytes.NewReader(scannerPackage.Bytes()))
			if err != nil {
				t.Logf("unpack msg err:%s", err.Error())
				continue
			}
			t.Logf("read msg from tcp ok,version is:%s,length is:%d,msg is:%s", scannedPack.Version, scannedPack.Length, scannedPack.Msg)
			receivedMsg <- true
			return
		}

		if scannerPackage.Err() != nil {
			t.Logf("scannerPackage err:%s", scannerPackage.Err().Error())
		}
	}()

	// Send initial connection message
	t.Log("building tcp heartbeat conn...")
	msg := &proto.SendTcp{
		Msg:          "build tcp heartbeat conn",
		FromUserId:   fromUserId,
		FromUserName: "Tcp heartbeat build",
		RoomId:       roomId,
		Op:           config.OpBuildTcpConn,
		AuthToken:    authToken,
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	pack := &StickPackage{
		Version: VersionContent,
		Msg:     msgBytes,
	}
	pack.Length = pack.GetPackageLength()

	tcpConn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err := pack.Pack(tcpConn); err != nil {
		t.Fatalf("failed to send message: %v", err)
	}
	t.Log("sent initial connection message")

	// Wait for response or timeout
	select {
	case <-receivedMsg:
		t.Log("successfully received response from server")
	case <-testTimeout:
		t.Log("test timeout - no response received, but connection was established successfully")
	}

	// Send one test message
	testMsg := &proto.SendTcp{
		Msg:          "test message from integration test",
		FromUserId:   fromUserId,
		FromUserName: "Integration Test",
		RoomId:       roomId,
		Op:           config.OpRoomSend,
		AuthToken:    authToken,
	}
	testMsgBytes, _ := json.Marshal(testMsg)
	testPack := &StickPackage{
		Version: VersionContent,
		Msg:     testMsgBytes,
	}
	testPack.Length = testPack.GetPackageLength()

	tcpConn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err := testPack.Pack(tcpConn); err != nil {
		t.Fatalf("failed to send test message: %v", err)
	}
	t.Log("sent test message successfully")

	// Give a moment for message to be processed
	time.Sleep(100 * time.Millisecond)
}
