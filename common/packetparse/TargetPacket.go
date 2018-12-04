package packetparse

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
	"unsafe"
)

var (
	_ = unsafe.Sizeof(0)
	_ = io.ReadFull
	_ = time.Now()
)

//go:generate msgp -io=false
//easyjson:json
type TargetPackets []*TargetPacket

//easyjson:json
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

func (tp *TargetPacket) CheckRecord() error {

	if tp.Type == "" {
		return errors.New("TargetPacket.Type is none")
	}

	if tp.HostName == "" {
		return errors.New("TargetPacket.HostName is none")
	}

	if tp.TimeStamp <= 0 {
		return errors.New("TargetPacket.TimeStamp le 0")
	}

	if tp.Plugin == "" {
		return errors.New("TargetPacket.Plugin is none")
	}

	if len(tp.Value) == 0 {
		return errors.New("TargetPacket.Value is none")
	}
	if len(strings.Split(tp.VlTags, "|")) < len(tp.Value) {
		return fmt.Errorf(" vltags:%s is not equals ", tp.VlTags)
	}

	return nil

}

func (d *TargetPacket) Size() (s uint64) {

	{
		l := uint64(len(d.HostName))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Plugin))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Instance))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Type))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Value))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}

		s += 8 * l

	}
	{
		l := uint64(len(d.VlTags))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Message))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	s += 8
	return
}
func (d *TargetPacket) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		l := uint64(len(d.HostName))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		copy(buf[i+0:], d.HostName)
		i += l
	}
	{

		v := *(*uint64)(unsafe.Pointer(&(d.TimeStamp)))

		buf[i+0+0] = byte(v >> 0)

		buf[i+1+0] = byte(v >> 8)

		buf[i+2+0] = byte(v >> 16)

		buf[i+3+0] = byte(v >> 24)

		buf[i+4+0] = byte(v >> 32)

		buf[i+5+0] = byte(v >> 40)

		buf[i+6+0] = byte(v >> 48)

		buf[i+7+0] = byte(v >> 56)

	}
	{
		l := uint64(len(d.Plugin))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.Plugin)
		i += l
	}
	{
		l := uint64(len(d.Instance))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.Instance)
		i += l
	}
	{
		l := uint64(len(d.Type))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.Type)
		i += l
	}
	{
		l := uint64(len(d.Value))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		for k0 := range d.Value {

			{

				v := *(*uint64)(unsafe.Pointer(&(d.Value[k0])))

				buf[i+0+8] = byte(v >> 0)

				buf[i+1+8] = byte(v >> 8)

				buf[i+2+8] = byte(v >> 16)

				buf[i+3+8] = byte(v >> 24)

				buf[i+4+8] = byte(v >> 32)

				buf[i+5+8] = byte(v >> 40)

				buf[i+6+8] = byte(v >> 48)

				buf[i+7+8] = byte(v >> 56)

			}

			i += 8

		}
	}
	{
		l := uint64(len(d.VlTags))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.VlTags)
		i += l
	}
	{
		l := uint64(len(d.Message))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.Message)
		i += l
	}
	return buf[:i+8], nil
}

func (d *TargetPacket) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		l := uint64(0)
		{
			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++
			l = t
		}

		d.HostName = string(buf[i+0 : i+0+l])
		i += l
	}
	{

		v := 0 | (uint64(buf[i+0+0]) << 0) | (uint64(buf[i+1+0]) << 8) | (uint64(buf[i+2+0]) << 16) | (uint64(buf[i+3+0]) << 24) | (uint64(buf[i+4+0]) << 32) | (uint64(buf[i+5+0]) << 40) | (uint64(buf[i+6+0]) << 48) | (uint64(buf[i+7+0]) << 56)
		d.TimeStamp = *(*float64)(unsafe.Pointer(&v))

	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Plugin = string(buf[i+8 : i+8+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Instance = string(buf[i+8 : i+8+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Type = string(buf[i+8 : i+8+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Value)) >= l {
			d.Value = d.Value[:l]
		} else {
			d.Value = make([]float64, l)
		}
		for k0 := range d.Value {

			{

				v := 0 | (uint64(buf[i+0+8]) << 0) | (uint64(buf[i+1+8]) << 8) | (uint64(buf[i+2+8]) << 16) | (uint64(buf[i+3+8]) << 24) | (uint64(buf[i+4+8]) << 32) | (uint64(buf[i+5+8]) << 40) | (uint64(buf[i+6+8]) << 48) | (uint64(buf[i+7+8]) << 56)
				d.Value[k0] = *(*float64)(unsafe.Pointer(&v))

			}

			i += 8

		}
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.VlTags = string(buf[i+8 : i+8+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Message = string(buf[i+8 : i+8+l])
		i += l
	}
	return i + 8, nil
}

func (tps TargetPackets) Size() (s uint64) {

	l := uint64(len(tps))
	{
		t := l
		for t >= 0x80 {
			t >>= 7
			s++
		}
		s++
	}

	for k0 := range tps {
		s += tps[k0].Size()
	}

	return
}

func (tps TargetPackets) Marshal() ([]byte, error) {
	size := tps.Size()

	buf := make([]byte, size)

	i := uint64(0)

	{
		l := uint64(len(tps))
		t := uint64(l)

		for t >= 0x80 {
			buf[i+0] = byte(t) | 0x80
			t >>= 7
			i++
		}

		buf[i+0] = byte(t)
		i++

		for k0 := range tps {
			nbuf, err := tps[k0].Marshal(buf[i+0:])
			if err != nil {
				return nil, err
			}
			i += uint64(len(nbuf))
		}
	}
	return buf[:i+0], nil
}

func (tps *TargetPackets) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	l := uint64(0)
	{
		bs := uint8(7)
		t := uint64(buf[i+0] & 0x7F)
		for buf[i+0]&0x80 == 0x80 {
			i++
			t |= uint64(buf[i+0]&0x7F) << bs
			bs += 7
		}
		i++

		l = t

	}

	*tps = make(TargetPackets, l)

	for k0 := range *tps {
		(*tps)[k0] = &TargetPacket{}
		ni, err := (*tps)[k0].Unmarshal(buf[i+0:])
		if err != nil {
			return 0, err
		}
		i += ni
	}

	return i, nil
}
