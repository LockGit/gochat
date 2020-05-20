/**
 * Created by lock
 * Date: 2020/5/20
 */
package stickpackage

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
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

	for scanner.Scan() {
		scannedPack := new(StickPackage)
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
