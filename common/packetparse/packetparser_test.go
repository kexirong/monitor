package packetparse

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func Benchmark_MsgTargetPacket(b *testing.B) {
	v := TargetPacket{
		HostName:  "kk-debian",
		TimeStamp: 1524205995.484389,
		Plugin:    "cpus1",
		Type:      "percent",
		Instance:  "0",
		Value:     []float64{12.1, 0.0, 12.6, 73.7},
		VlTags:    "user|nice|system|idle",
	}

	//b.SetBytes(int64(len(bts)))
	//b.ReportAllocs()
	bts, err := v.MarshalMsg(nil)
	_ = bts
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = v.UnmarshalMsg(bts)
		if err != nil {
			b.Fatal(err)
		}

	}
}

func Benchmark_gencode(b *testing.B) {

	var p = TargetPacket{
		HostName:  "kk-debian",
		TimeStamp: 1524205995.484389,
		Plugin:    "cpus1",
		Type:      "percent",
		Instance:  "0",
		Value:     []float64{12.1, 0.0, 12.6, 73.7},
		VlTags:    "user|nice|system|idle",
	}
	bb, err := p.Marshal(nil)
	_ = bb
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		ubb, err := p.Unmarshal(bb)
		_ = ubb
		if err != nil {
			b.Error(err)
		}

	}

}

func Benchmark_easyjson(b *testing.B) {
	t := TargetPacket{
		HostName:  "kk-debian",
		TimeStamp: 1524205995.484389,
		Plugin:    "cpus1",
		Type:      "percent",
		Instance:  "0",
		Value:     []float64{12.1, 0.0, 12.6, 73.7},
		VlTags:    "user|nice|system|idle",
	}
	bs, err := t.MarshalJSON()
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		err = t.UnmarshalJSON(bs)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_jsoniter(b *testing.B) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonstr := `{"plugin": "cpus1", "timestamp": 1524205995.484389, "hostname": "kk-debian", "value": [12.1, 0.0, 12.6, 73.7], "instance": "0", "vltags": "user|nice|system|idle", "type": "percent"}`
	var p TargetPacket
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal([]byte(jsonstr), &p)
		if err != nil {
			b.Error(err)
		}

	}

}

func Benchmark_stdjson(b *testing.B) {

	jsonstr := `{"plugin": "cpus1", "timestamp": 1524205995.484389, "hostname": "kk-debian", "value": [12.1, 0.0, 12.6, 73.7], "instance": "0", "vltags": "user|nice|system|idle", "type": "percent"}`
	var p TargetPacket
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal([]byte(jsonstr), &p)
		if err != nil {
			b.Error(err)
		}

	}

}

func Test_jsoniter(t *testing.T) {
	var json = jsoniter.ConfigFastest
	jsonstr := `[{"plugin": "cpus1", "timestamp": 1524205995.484389, "hostname": "kk-debian", "value": [12.1, 0.0, 12.6, 73.7], "instance": "0", "vltags": "user|nice|system|idle", "type": "percent"}, 
	{"plugin": "cpu1", "timestamp": 1524205995.484393, "hostname": "kk-debian", "value": [11.5, 0.0, 12.5, 74.5], "instance": "1", "vltags": "user|nice|system|idle", "type": "percent"},
	{"plugin": "cpu1", "timestamp": 1524205995.484393, "hostname": "kk-debian", "value": [11.5, 0.0, 12.5, 74.5], "instance": "1", "vltags": "user|nice|system|idle", "type": "percent"},
	{"plugin": "cpu1", "timestamp": 1524205995.484393, "hostname": "kk-debian", "value": [11.5, 0.0, 12.5, 74.5], "instance": "1", "vltags": "user|nice|system|idle", "type": "percent"}]`
	var pp []*TargetPacket

	err := json.Unmarshal([]byte(jsonstr), &pp)
	if err != nil {
		t.Error(err)
	}
	_, err = json.Marshal(pp)
	if err != nil {
		t.Error(err)
	}

	for _, p := range pp {
		t.Logf("%#v", p)
	}

}

func Test_gencode(t *testing.T) {

	var tps = TargetPackets{
		&TargetPacket{
			HostName:  "kk-debian",
			TimeStamp: 1524205995.484389,
			Plugin:    "cpus",
			Type:      "percent",
			Instance:  "0",
			Value:     []float64{12.1, 0.0, 12.6, 73.7},
			VlTags:    "user|nice|system|idle",
		},
		&TargetPacket{
			HostName:  "kk-debian",
			TimeStamp: 1524205996.484389,
			Plugin:    "cpus1",
			Type:      "percent",
			Instance:  "1",
			Value:     []float64{12.1, 0.0, 12.6, 73.7},
			VlTags:    "user|nice|system|idle",
		},
		&TargetPacket{
			HostName:  "kk-debian",
			TimeStamp: 1524205997.484389,
			Plugin:    "cpus2",
			Type:      "percent",
			Instance:  "2",
			Value:     []float64{12.1, 0.0, 12.6, 73.7},
			VlTags:    "user|nice|system|idle",
		},
	}
	bb, err := tps.Marshal()
	_ = bb
	if err != nil {
		t.Error(err)
	}
	tps = TargetPackets{}
	l, err := tps.Unmarshal(bb)
	if err != nil || len(bb) != int(l) {
		t.Error(err)
	}
	t.Log(tps)
}
