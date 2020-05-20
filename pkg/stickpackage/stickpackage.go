/**
 * Created by lock
 * Date: 2020/5/20
 */
package stickpackage

import (
	"encoding/binary"
	"fmt"
	"io"
)

var VersionLength = 2
var LengthLength = 2
var TcpHeaderLength = VersionLength + LengthLength
var LengthStartIndex = 2 // 数据部分长度起始字节位置
var LengthStopIndex = 4  // 数据部分长度结束字节位置
var VersionContent = [2]byte{'v', '1'}

type StickPackage struct {
	Version [2]byte // 协议版本
	Length  int16   // 数据部分长度
	Msg     []byte  // 数据
}

func (p *StickPackage) Pack(writer io.Writer) error {
	var err error
	err = binary.Write(writer, binary.BigEndian, &p.Version)
	err = binary.Write(writer, binary.BigEndian, &p.Length)
	err = binary.Write(writer, binary.BigEndian, &p.Msg)
	return err
}

func (p *StickPackage) Unpack(reader io.Reader) error {
	var err error
	err = binary.Read(reader, binary.BigEndian, &p.Version)
	err = binary.Read(reader, binary.BigEndian, &p.Length)
	p.Msg = make([]byte, p.Length-4)
	err = binary.Read(reader, binary.BigEndian, &p.Msg)
	return err
}

func (p *StickPackage) String() string {
	return fmt.Sprintf("version:%s length:%d msg:%s",
		p.Version,
		p.Length,
		p.Msg,
	)
}

func (p *StickPackage) GetPackageLength() int16 {
	p.Length = 4 + int16(len(p.Msg))
	return p.Length
}
