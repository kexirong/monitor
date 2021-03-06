package packetparse

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
)

/*PDU 数据包格式:
 加入version字段可以大大降低处理粘包错误概率
 length字段为type、payload、check、\r\n长度之和
 +------+--------+--------+--------+---------+-------+------+
 | head | version| length |  type  | payload | check | \r\n |
 +------+--------+--------+--------+---------+-------+------+

head    :          --2byte  define  :uint8(0x07)uint8(0x02)
version :		   --3byte  current : 0.1
length  : uint16   --2byte  sum     :1 + len(payload) + 4
type    : uint8    --1byte
payload : []byte   --bytes
check   : crc32    --4byte
\r\n    : end mark --2byte


0x01 normal --default  common string
0x02 json
0x03 heartbeat
0x04 reply //接收端 接收到数据包之后 需要reply, check:填收到的载核的crc32；payload为 ”ok“
*/
type PDU struct {
	Type    _type
	Length  uint16
	Payload []byte //最长65000个字节，超过encode报错
	Check   uint32
}

//pduVersion 必须3个字节长度
var (
	pduVersion = "0.1"

	//PDUhead 供校验头部信息使用
	PDUhead = append([]byte{0x07, 0x02}, pduVersion...)

	//PDUeof 供校验尾部信息使用
	PDUeof = []byte{'\r', '\n'}
)

//PDUType 供调用者封包解包使用
type _type uint8

// 后续可能将 类型跟类别分开定义
const (
	// PDUReply  应答报文
	PDUReply = _type(0x01)
	// PDUHeartBeat  心跳包
	PDUHeartBeat = _type(0x02)
	// PDUJson  TargetPacket的json报文
	PDUJson = _type(0x03)
	// PDUTargetPacket  TargetPacket报文的msgp包
	PDUTargetPacket = _type(0x04)
	// PDUTargetPackets  TargetPacket 批量报文的msgp包
	PDUTargetPackets = _type(0x05)
	// PDUOther  未定义
	PDUOther = _type(0x06)
)

func GenPduWithPayload(tp _type, payload []byte) (PDU, error) {
	var pdu PDU

	pdu.Type = tp
	if len(payload) == 0 || len(payload) > 65000 {
		return pdu, errors.New("GenPduWithPayload error: payload length is 0 or too long ")
	}
	pdu.Length = uint16(len(payload) + 1 + 4)
	pdu.Payload = payload
	pdu.Check = crc32.ChecksumIEEE(payload)
	return pdu, nil
}

//PDUEncode 这里不做Checksum校验，仅判断不为0
func PDUEncode(pdu PDU) ([]byte, error) {
	buf := new(bytes.Buffer)
	if pdu.Type == 0 || pdu.Check == 0 {
		return nil, errors.New("PdUEncode error:  Type or Check  is 0")
	}

	leng := len(pdu.Payload)
	if leng > 65000 {
		return nil, errors.New("PdUEncode error: length is too long")
	}
	if pdu.Length != uint16(leng+1+4) {
		return nil, errors.New("PdUEncode error: payload Length is not equal")
	}
	_, err := buf.Write(PDUhead)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(Network.Uint16ToBytes(pdu.Length))
	if err != nil {
		return nil, err
	}

	err = buf.WriteByte(byte(pdu.Type))
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(pdu.Payload)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(Network.Uint32ToBytes(pdu.Check))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(PDUeof)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/*ReadPDU 读取 +--------+---------+-------+ 部分
              |  type  | payload | check |
  			  +--------+---------+-------+
  校验type、check任务在PDUDecode函数中进行
*/
func ReadPDU(conn io.Reader) ([]byte, error) {

	reader := bufio.NewReader(conn)
	var count = 0
	var tmp = make([]byte, 2)

	for i := 0; i < 5; i++ {
		bit, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if bit == PDUhead[i] {
			count++
			continue
		}
		if bit == 0x07 {
			i = 0
			count = 1
		}
	}

	if count != 5 {
		//	fmt.Println("break and return")
		return nil, errors.New("ReadPDU error: read head failed")
	}

	n, err := reader.Read(tmp)
	if err != nil {
		return nil, err
	}

	if n < 2 {
		return nil, errors.New("ReadPDU error: read length failed")
	}

	leng := int(Network.BytesToUint16(tmp))
	rst := make([]byte, leng)
	for {
		n = 0
		i, err := reader.Read(rst[n:])

		if err != nil {
			return nil, err
		}
		n += i
		if n >= leng {
			break
		}
	}

	if n != leng {
		return nil, fmt.Errorf("ReadPDU error: data leng(%d) not equal length(%d)", n, leng)
	}

	tmp = make([]byte, 2)
	if _, err = reader.Read(tmp); err != nil {
		return nil, err
	}

	if !bytes.Equal(tmp, PDUeof) {
		return nil, errors.New("ReadPDU error: read PDUeof failed")
	}

	return rst, nil
}

/*PDUDecode 参数 arg 格式内容:
  +--------+---------+-------+
  |  type  | payload | check |
  +--------+---------+-------+
*/
func PDUDecode(bits []byte) (PDU, error) {
	var pdu PDU
	pdu.Type = _type(bits[0])

	leng := len(bits)
	pdu.Payload = bits[1 : leng-4]
	pdu.Check = Network.BytesToUint32(bits[leng-4:])
	if pdu.Check != crc32.ChecksumIEEE(pdu.Payload) && pdu.Type != 0x01 {
		return pdu, errors.New("PDUDecode error: payload check failed")
	}
	pdu.Length = uint16(len(bits))
	return pdu, nil
}

//PDUEncodeReply .
func PDUEncodeReply(check uint32, msg []byte) ([]byte, error) {
	if len(msg) > 65000 {
		return nil, errors.New("msg too lang")
	}
	var pdu = PDU{
		Type:    0x01,
		Payload: msg,
		Check:   check,
		Length:  uint16(len(msg) + 1 + 4),
	}
	return PDUEncode(pdu)
}
