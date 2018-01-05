package packetparse

import (
	"bytes"
	"fmt"
	"io"
)

func TargetPackage(pk TargetPacket) ([]byte, error) {

	buf := new(bytes.Buffer)
	if pk.HostName == "" {
		return nil, fmt.Errorf("TargetPacket.HostName is none")
	}

	if err := targetWriteBuf(buf, targetPacketMap["hostname"], pk.HostName); err != nil {
		return nil, err
	}
	if pk.TimeStamp <= 0 {
		return nil, fmt.Errorf("TargetPacket.TimeStamp le 0")
	}
	if err := targetWriteBuf(buf, targetPacketMap["timestamp"], pk.TimeStamp); err != nil {
		return nil, err
	}
	if pk.Plugin == "" {
		return nil, fmt.Errorf("TargetPacket.Plugin is none")
	}
	if err := targetWriteBuf(buf, targetPacketMap["plugin"], pk.Plugin); err != nil {
		return nil, err
	}

	if pk.Value == nil {
		return nil, fmt.Errorf("TargetPacket.Value is none")
	}

	if err := targetWriteBuf(buf, targetPacketMap["value"], pk.Value); err != nil {
		return nil, err
	}
	if pk.VlTags == "" {
		return nil, fmt.Errorf("TargetPacket.VlTags is none")
	}
	if err := targetWriteBuf(buf, targetPacketMap["vltags"], pk.VlTags); err != nil {
		return nil, err
	}

	if pk.Instance != "" {
		if err := targetWriteBuf(buf, targetPacketMap["instance"], pk.Instance); err != nil {
			return nil, err
		}
	}
	if pk.Message != "" {
		if err := targetWriteBuf(buf, targetPacketMap["message"], pk.Message); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil

}

func TargetParse(b []byte) (TargetPacket, error) {
	var err error
	var packet TargetPacket
	var t uint16
	var l uint32
	var str string
	var f float64
	var sf []float64
	n := 0

	head := make([]byte, 6)

	bufs := bytes.NewReader(b)
	for ; ; n++ {
		_, err = bufs.Read(head)

		if err == io.EOF {
			break
		}
		t = Network.BytesToUint16(head[0:2])
		l = Network.BytesToUint32(head[2:6])

		buf := make([]byte, l)

		_, err = bufs.Read(buf)
		if err != nil {
			return packet, err
		}

		field := targetParseMap[t]
		kind := targetTypesMap[field]

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

func targetWriteBuf(buf io.Writer, seq uint16, data interface{}) error {
	var err error
	var leng uint32
	_, err = buf.Write(Network.Uint16ToBytes(seq))
	if err != nil {
		return err
	}
	switch data.(type) {

	case string:
		leng = uint32(len(data.(string)))
		_, err = buf.Write(Network.Uint32ToBytes(leng))

		if err != nil {
			return err
		}

		_, err = buf.Write([]byte(data.(string)))

		if err != nil {
			return err
		}

	case float64:
		leng = 8
		_, err = buf.Write(Network.Uint32ToBytes(leng))

		if err != nil {
			return err
		}

		_, err = buf.Write(Network.Float64ToBytes(data.(float64)))

		if err != nil {
			return err
		}

	case []float64:
		leng = uint32(8 * len(data.([]float64)))
		_, err = buf.Write(Network.Uint32ToBytes(leng))

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
