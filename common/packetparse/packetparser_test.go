package packetparse

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func Benchmark_packetparse(b *testing.B) {

	var pp = []TargetPacket{TargetPacket{
		HostName:  "kk-debian",
		TimeStamp: 1524205995.484389,
		Plugin:    "cpus1",
		Type:      "percent",
		Instance:  "0",
		Value:     []float64{12.1, 0.0, 12.6, 73.7},
		VlTags:    "user|nice|system|idle",
	},
		TargetPacket{
			HostName:  "kk-debian",
			TimeStamp: 1524205995.484393,
			Plugin:    "cpus1",
			Type:      "percent",
			Instance:  "1",
			Value:     []float64{11.5, 0.0, 12.5, 74.5},
			VlTags:    "user|nice|system|idle",
		},
		TargetPacket{
			HostName:  "kk-debian",
			TimeStamp: 1524205995.484393,
			Plugin:    "cpus1",
			Type:      "percent",
			Instance:  "1",
			Value:     []float64{11.5, 0.0, 12.5, 74.5},
			VlTags:    "user|nice|system|idle",
		},
		TargetPacket{
			HostName:  "kk-debian",
			TimeStamp: 1524205995.484393,
			Plugin:    "cpus1",
			Type:      "percent",
			Instance:  "1",
			Value:     []float64{11.5, 0.0, 12.5, 74.5},
			VlTags:    "user|nice|system|idle",
		},
	}

	for i := 0; i < b.N; i++ {
		for _, p := range pp {
			bb, err := TargetPackage(p)

			if err != nil {
				b.Error(err)
			}

			_, err = TargetParse(bb)
			if err != nil {
				b.Error(err)
			}
			//b.Log(st.String())
		}
	}

}

func Benchmark_json(b *testing.B) {
	jsonstr := `[{"plugin": "cpus1", "timestamp": 1524205995.484389, "hostname": "kk-debian", "value": [12.1, 0.0, 12.6, 73.7], "instance": "0", "vltags": "user|nice|system|idle", "type": "percent"}, 
	{"plugin": "cpu1", "timestamp": 1524205995.484393, "hostname": "kk-debian", "value": [11.5, 0.0, 12.5, 74.5], "instance": "1", "vltags": "user|nice|system|idle", "type": "percent"}]`
	var pp []TargetPacket
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal([]byte(jsonstr), &pp)
		if err != nil {
			b.Error(err)
		}
		_, err = json.Marshal(pp)
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_jsoniter(b *testing.B) {
	var json = jsoniter.ConfigFastest
	jsonstr := `[{"plugin": "cpus1", "timestamp": 1524205995.484389, "hostname": "kk-debian", "value": [12.1, 0.0, 12.6, 73.7], "instance": "0", "vltags": "user|nice|system|idle", "type": "percent"}, 
	{"plugin": "cpu1", "timestamp": 1524205995.484393, "hostname": "kk-debian", "value": [11.5, 0.0, 12.5, 74.5], "instance": "1", "vltags": "user|nice|system|idle", "type": "percent"},
	{"plugin": "cpu1", "timestamp": 1524205995.484393, "hostname": "kk-debian", "value": [11.5, 0.0, 12.5, 74.5], "instance": "1", "vltags": "user|nice|system|idle", "type": "percent"},
	{"plugin": "cpu1", "timestamp": 1524205995.484393, "hostname": "kk-debian", "value": [11.5, 0.0, 12.5, 74.5], "instance": "1", "vltags": "user|nice|system|idle", "type": "percent"}]`
	var pp []TargetPacket
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal([]byte(jsonstr), &pp)
		if err != nil {
			b.Error(err)
		}
		_, err = json.Marshal(pp)
		if err != nil {
			b.Error(err)
		}
	}
}
