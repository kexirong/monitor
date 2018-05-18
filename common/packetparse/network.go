package packetparse

import (
	"time"
	"unsafe"
)

var Network network

type network struct{}

//network(=big-Endian)
func (network) Uint64ToBytes(ui64 uint64) []byte {
	b := make([]byte, 8)
	b[0] = byte(ui64 >> 56)
	b[1] = byte(ui64 >> 48)
	b[2] = byte(ui64 >> 40)
	b[3] = byte(ui64 >> 32)
	b[4] = byte(ui64 >> 24)
	b[5] = byte(ui64 >> 16)
	b[6] = byte(ui64 >> 8)
	b[7] = byte(ui64)
	return b
}

func (network) Uint32ToBytes(ui32 uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(ui32 >> 24)
	b[1] = byte(ui32 >> 16)
	b[2] = byte(ui32 >> 8)
	b[3] = byte(ui32)
	return b
}

func (network) Uint16ToBytes(u uint16) []byte {
	b := make([]byte, 2)
	b[0] = byte(u >> 8)
	b[1] = byte(u)
	return b

}

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
	if len(b) != 8 {
		return 0
	}
	bits := uint64(b[0])<<56 | uint64(b[1])<<48 |
		uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 |
		uint64(b[6])<<8 | uint64(b[7])
	return *(*float64)(unsafe.Pointer(&bits))
}

func (network) BytesToUint16(b []byte) uint16 {
	if len(b) != 2 {
		return 0
	}
	return uint16(b[1]) | uint16(b[0])<<8
}

func (network) BytesToUint64(b []byte) uint64 {
	if len(b) != 8 {
		return 0
	}
	bits := uint64(b[0])<<56 | uint64(b[1])<<48 |
		uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 |
		uint64(b[6])<<8 | uint64(b[7])
	return bits
}

func (network) BytesToUint32(b []byte) uint32 {
	if len(b) != 4 {
		return 0
	}
	bits := uint32(b[0])<<24 | uint32(b[1])<<16 |
		uint32(b[2])<<8 | uint32(b[3])
	return bits
}

//Nsecond2Unix is  NanoSecond To UnixTimetamp
func Nsecond2Unix(ns int64) float64 {
	return time.Duration(ns).Seconds()
}
