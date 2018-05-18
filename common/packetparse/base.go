package packetparse

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
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
	Type    uint8
	Length  uint16
	Payload []byte
	Check   uint32
}

//pduVersion 必须3个字节长度
var (
	pduVersion = "0.1"

	//PDUhead 供校验头部信息使用
	PDUhead = append([]byte{07, 02}, pduVersion...)

	//PDUeof 供校验尾部信息使用
	PDUeof = []byte{'\r', '\n'}
)

//PDUTypeMap 供调用者封包解包使用
var PDUTypeMap = map[uint8]string{
	0x01: "reply",
	0x02: "heartbeat",
	0x03: "json",
	0x04: "normal",
	0x05: "targetpackage",
}

func GenPduWithPayload(tp uint8, payload []byte) (PDU, error) {
	var pdu PDU
	if _, ok := PDUTypeMap[tp]; !ok {
		return pdu, errors.New("GenPduWithPayload error: tp not in PDUTypeMap")
	}
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
	n, err = reader.Read(rst)

	if err != nil {
		return nil, err
	}

	if n != leng {
		return nil, errors.New("ReadPDU error: data leng not equal length")
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
	pdu.Type = uint8(bits[0])
	if _, ok := PDUTypeMap[pdu.Type]; !ok {
		return pdu, errors.New("PDUDecode error: type check failed")
	}
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

/*TargetPacket .

+------+--------+---------+
| type | length |  data   |
+------+--------+---------+
|         .......         |
+------+--------+---------+

type    : uint16   --2byte
length  : uint16   --2byte
data    : []byte   --[length]byte
*/
type TargetPacket struct {
	HostName  string    `json:"hostname"`  //ops201
	TimeStamp float64   `json:"timestamp"` //the number of seconds elapsed since January 1, 1970 UTC
	Plugin    string    `json:"plugin"`    // cpu
	Instance  string    `json:"instance"`  // 0,1,2,3 (eth0,eth1)(sda,sdb)
	Type      string    `json:"type"`      //percent(百分比),bool(0|1),gauge(原值),derive(速率,单位v/s)
	Value     []float64 `json:"value"`     //float 对整数兼容，故采用float64而不是interface{}
	VlTags    string    `json:"vltags"`    // "idle|user|system"(rx|tx)(read|write|use|free...)
	Message   string    `json:"message"`   // description ,e.g: the disk is full please clean
}

func (tp *TargetPacket) String() string {
	return fmt.Sprintf(`
		hostname:%s
		timestamp:%v
		plugin:%s
		instance:%s
		type:%s
		value:%v
		vltags:%s
		message:%s
		`,
		tp.HostName,
		tp.TimeStamp,
		tp.Plugin,
		tp.Instance,
		tp.Type,
		tp.Value,
		tp.VlTags,
		tp.Message,
	)
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
