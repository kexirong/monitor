package packetparse

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"unsafe"
)

var Network network

type network struct{}

//network(=big-Endian)
func (network) Float64ToBytes(f float64) []byte {
	bits := *(*uint64)(unsafe.Pointer(&f))
	b := make([]byte, 8)
	b[0] = byte(bits >> 56)
	b[1] = byte(bits >> 48)
	b[2] = byte(bits >> 40)
	b[3] = byte(bits >> 32)
	b[4] = byte(bits >> 24)
	b[5] = byte(bits >> 16)
	b[6] = byte(bits >> 8)
	b[7] = byte(bits)

	return b
}

//native(=little-Endian)
func (network) BytesToFloat64(b []byte) float64 {

	bits := uint64(b[0])<<56 | uint64(b[1])<<48 |
		uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 |
		uint64(b[6])<<8 | uint64(b[7])

	return *(*float64)(unsafe.Pointer(&bits))
}

//network(=big-Endian)
func (network) Uint16ToBytes(u uint16) []byte {
	b := make([]byte, 2)

	b[0] = byte(u >> 8)
	b[1] = byte(u)

	return b

}

func (network) BytesToUint16(b []byte) uint16 {
	return uint16(b[1]) | uint16(b[0])<<8

}

func bufwrite(buf io.Writer, data interface{}) error {
	var err error
	var leng uint16

	switch data.(type) {

	case string:
		leng = uint16(len(data.(string)))
		_, err = buf.Write(Network.Uint16ToBytes(leng))

		if err != nil {
			return err
		}

		_, err = buf.Write([]byte(data.(string)))

		if err != nil {
			return err
		}

	case float64:
		leng = 8
		_, err = buf.Write(Network.Uint16ToBytes(leng))

		if err != nil {
			return err
		}

		_, err = buf.Write(Network.Float64ToBytes(data.(float64)))

		if err != nil {
			return err
		}

	case []float64:
		leng = uint16(8 * len(data.([]float64)))
		_, err = buf.Write(Network.Uint16ToBytes(leng))

		if err != nil {
			return err
		}

		for _, value := range data.([]float64) {
			_, err = buf.Write(Network.Float64ToBytes((value)))
		}
	}

	if err != nil {
		return err
	}

	return nil

}

func Package(pk Packet) ([]byte, error) {
	/*
		if len(pk.Value) < 1 || pk.VlTags == "" {
			return nil, errors.New(" vltags  or   value is nil ")
		}

			if (len(pk.Value) > 1 || pk.Instance == "") && pk.VlTags == "" {
				return nil, errors.New(" VlTags is ''  but value.len ne 1 or instance is also '' ")
			}
	*/
	var err error
	buf := new(bytes.Buffer)

	value := reflect.ValueOf(pk)
	num := value.NumField()

	/*if value.Field(num-1) == "" {
	    num = num -1
	}*/

	for i := 0; i < num; i++ {
		v := value.Field(i)

		//content

		switch v.Interface() {

		case "", float64(0), nil:

			if !(uint16(i) == packMap["message"] || uint16(i) == packMap["instance"]) {
				return nil, fmt.Errorf("Field: %s is empty", parseMap[uint16(i)])
			}

		default:
			_, err = buf.Write(Network.Uint16ToBytes(uint16(i)))

			if err != nil {
				return nil, err
			}

			err = bufwrite(buf, v.Interface())

			if err != nil {
				return nil, err
			}
		}
	}
	return buf.Bytes(), nil

}

func Parse(b []byte) (Packet, error) {
	var err error
	var packet Packet
	var t, l uint16
	var str string
	var f float64
	var sf []float64
	n := 0

	head := make([]byte, 4)

	bufs := bytes.NewReader(b)
	for ; ; n++ {
		_, err = bufs.Read(head)

		if err == io.EOF {
			break
		}
		t = Network.BytesToUint16(head[0:2])
		l = Network.BytesToUint16(head[2:4])

		buf := make([]byte, l)

		_, err = bufs.Read(buf)
		if err != nil {
			return packet, err
		}

		field := parseMap[t]
		kind := typesMap[field]

		switch kind {

		case "string":
			str = string(buf)

		case "float64":
			f = Network.BytesToFloat64(buf)

		case "slice":
			e := 8
			for i := 0; i < int(l); i += 8 {
				f = Network.BytesToFloat64(buf[i : i+e])
				sf = append(sf, f)
			}

		default:
			return packet, fmt.Errorf("packet parse type error : n=%d, t=%d", n, t)

		}

		switch field {
		case "hostname":
			packet.HostName = str

		case "timestamp":
			packet.TimeStamp = f

		case "plugin":
			packet.Plugin = str

		case "instance":
			packet.Instance = str
		case "type":
			packet.Type = str

		case "value":
			packet.Value = sf
		case "vltags":
			packet.VlTags = str
		case "message":
			packet.Message = str
		default:
			return packet, fmt.Errorf("packet parse field error : n=%d, t=%d", n, t)

		}

	}

	return packet, nil

}
