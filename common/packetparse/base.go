package packetparse

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

/*PDU .
+------+--------+---------+-------+------+
| type | length | payload | check | \r\n |
+-------+-------+---------+-------+------+

type    : uint8    --1byte
length  : uint64   --8byte
payload : []byte   --bytes
check   : crc32    --4byte
\r\n    : end mark --2byte


0x01 normal #default  common string
0x02 json
0x03 heartbeat
0x04 reply //接收端 接收到数据包之后 需要reply check 用收到的载核的 crc32；payload 中为 ”ok“
*/
type PDU struct {
	Type    uint8
	Length  uint64
	Payload []byte
	Check   uint32
}

var PDUTypeMap = map[uint8]string{
	0x01: "normal",
	0x02: "json",
	0x03: "heartbeat",
	0x04: "reply",
}

func PdUEncode(pdu PDU) ([]byte, error) {
	buf := new(bytes.Buffer)
	if pdu.Type == 0 {
		return nil, fmt.Errorf("PDU.Type is 0")
	}

	err := buf.WriteByte(byte(pdu.Type))
	if err != nil {
		return nil, err
	}
	pdu.Length = uint64(len(pdu.Payload))
	if pdu.Length == 0 {
		return nil, fmt.Errorf("PDU.Length is 0")
	}

	_, err = buf.Write(Network.Uint64ToBytes(pdu.Length))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func PDUDecode(rd io.Reader) (PDU, error) {

	return PDU{}, nil
}

/*TargetPacket .

+------+--------+---------+
| type | length |  data   |
+------+--------+---------+
|         .......         |
+------+--------+---------+

type    : uint16    --2byte
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
var targetParseMap = make(map[uint16]string) //[id]name

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
