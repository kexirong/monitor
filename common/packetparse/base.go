package packetparse

import (
	"bytes"
	"errors"
	"fmt"
	"hash/crc32"
	"reflect"
)

/*PDU 数据包格式:
 加入version字段可以大大降低处理粘包错误概率
 length字段为type、payload、check、\r\n长度之和
 +------+--------+--------+--------+---------+-------+------+
 | head | version| length |  type  | payload | check | \r\n |
 +------+--------+--------+--------+---------+-------+------+

head    :          --2byte  define  :uint8(0x07)uint8(0x02)
version :		   --3byte  current : 0.1
type    : uint8    --1byte
length  : uint32   --4byte  sum     :1 + len(payload) + 4 + 2
payload : []byte   --bytes
check   : crc32    --4byte
\r\n    : end mark --2byte


0x01 normal --default  common string
0x02 json
0x03 heartbeat
0x04 reply //接收端 接收到数据包之后 需要reply, check:填收到的载核的crc32；payload为 ”ok“
*/
type PDU struct {
	Type    uint8
	Length  uint32
	Payload []byte
	Check   uint32
}

//pduVersion must be 3 bytes !!!
var pduVersion = []byte("0.1")

//PDUTypeMap 供调用者封包解包使用
var PDUTypeMap = map[uint8]string{
	0x01: "normal",
	0x02: "json",
	0x03: "heartbeat",
	0x04: "reply",
	0x05: "targetpackage",
}

//PdUEncode don't check payload
func PdUEncode(pdu PDU) ([]byte, error) {
	buf := new(bytes.Buffer)
	if pdu.Type == 0 {
		return nil, fmt.Errorf("PDU.Type is 0")
	}

	err := buf.WriteByte(byte(pdu.Type))
	if err != nil {
		return nil, err
	}
	pdu.Length = uint32(len(pdu.Payload))
	if pdu.Length == 0 {
		return nil, fmt.Errorf("PDU.Length is 0")
	}

	_, err = buf.Write(Network.Uint32ToBytes(pdu.Length))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/*PDUDecode 参数 arg 格式内容:
  +--------+---------+-------+
  |  type  | payload | check |
  +--------+---------+-------+
*/
func PDUDecode(bits []byte) (PDU, error) {
	var pdu PDU
	pdu.Type = uint8(bits[0])
	if _, ok := PDUTypeMap[pdu.Type]; !ok {
		return pdu, errors.New("PDUDecode error: type check failed")
	}
	leng := len(bits)
	pdu.Payload = bits[1 : leng-4]
	pdu.Check = Network.BytesToUint32(bits[leng-5:])
	if pdu.Check != crc32.ChecksumIEEE(pdu.Payload) {
		return pdu, errors.New("PDUDecode error: payload check failed")
	}
	return pdu, nil
}

/*TargetPacket .

+------+--------+---------+
| type | length |  data   |
+------+--------+---------+
|         .......         |
+------+--------+---------+

type    : uint16   --2byte
length  : uint32   --8byte
data    : []byte   --(length)byte
*/
type TargetPacket struct {
	HostName  string    `json:"hostname"`  //ops201
	TimeStamp float64   `json:"timestamp"` //the number of seconds elapsed since January 1, 1970 UTC
	Plugin    string    `json:"plugin"`    // cpu
	Instance  string    `json:"instance"`  // 0,1,2,3 (eth0,eth1)(sda,sdb)
	Type      string    `json:"type"`      //percent(百分比),counter(正数速率,主要是趋势),gauge(原值),derive(速率)
	Value     []float64 `json:"value"`
	VlTags    string    `json:"vltags"`  // "idle|user|system"(rx|tx)(read|write|use|free...)
	Message   string    `json:"message"` // description ,e.g: the disk is full please clean
}

//0为初始化值,所以不使用
var targetPacketMap = map[string]uint16{
	"hostname":  0x0001,
	"timestamp": 0x0002,
	"plugin":    0x0003,
	"instance":  0x0004,
	"type":      0x0005,
	"value":     0x0006,
	"vltags":    0x0007,
	"message":   0x0008,
}

var targetTypesMap = make(map[string]string) //[name]type
var targetParseMap = make(map[uint16]string) //[seq]name
//利用init函数使用反射初始化targetTypesMap和targetParseMap
func init() {

	for key, vl := range targetPacketMap {
		targetParseMap[vl] = key
	}

	var packet TargetPacket
	t := reflect.TypeOf(packet)
	v := reflect.ValueOf(packet)

	for k := 0; k < t.NumField(); k++ {
		targetTypesMap[t.Field(k).Tag.Get("json")] = v.Field(k).Kind().String()
	}
}
