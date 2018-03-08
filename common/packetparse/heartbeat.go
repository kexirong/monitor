package packetparse

func Heartbeat() []byte {
	pdu, err := GenPduWithPayload(0x02, []byte("heart"))
	if err != nil {
		return nil
	}
	bits, err := PDUEncode(pdu)
	if err != nil {
		return nil
	}
	return bits
}
