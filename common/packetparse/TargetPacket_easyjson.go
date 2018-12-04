// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package packetparse

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson1bfd9a9bDecodeGithubComKexirongMonitorCommonPacketparse(in *jlexer.Lexer, out *TargetPackets) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(TargetPackets, 0, 8)
			} else {
				*out = TargetPackets{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 *TargetPacket
			if in.IsNull() {
				in.Skip()
				v1 = nil
			} else {
				if v1 == nil {
					v1 = new(TargetPacket)
				}
				(*v1).UnmarshalEasyJSON(in)
			}
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1bfd9a9bEncodeGithubComKexirongMonitorCommonPacketparse(out *jwriter.Writer, in TargetPackets) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			if v3 == nil {
				out.RawString("null")
			} else {
				(*v3).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v TargetPackets) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1bfd9a9bEncodeGithubComKexirongMonitorCommonPacketparse(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TargetPackets) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1bfd9a9bEncodeGithubComKexirongMonitorCommonPacketparse(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TargetPackets) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1bfd9a9bDecodeGithubComKexirongMonitorCommonPacketparse(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TargetPackets) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1bfd9a9bDecodeGithubComKexirongMonitorCommonPacketparse(l, v)
}
func easyjson1bfd9a9bDecodeGithubComKexirongMonitorCommonPacketparse1(in *jlexer.Lexer, out *TargetPacket) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "hostname":
			out.HostName = string(in.String())
		case "timestamp":
			out.TimeStamp = float64(in.Float64())
		case "plugin":
			out.Plugin = string(in.String())
		case "instance":
			out.Instance = string(in.String())
		case "type":
			out.Type = string(in.String())
		case "value":
			if in.IsNull() {
				in.Skip()
				out.Value = nil
			} else {
				in.Delim('[')
				if out.Value == nil {
					if !in.IsDelim(']') {
						out.Value = make([]float64, 0, 8)
					} else {
						out.Value = []float64{}
					}
				} else {
					out.Value = (out.Value)[:0]
				}
				for !in.IsDelim(']') {
					var v4 float64
					v4 = float64(in.Float64())
					out.Value = append(out.Value, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "vltags":
			out.VlTags = string(in.String())
		case "message":
			out.Message = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1bfd9a9bEncodeGithubComKexirongMonitorCommonPacketparse1(out *jwriter.Writer, in TargetPacket) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"hostname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.HostName))
	}
	{
		const prefix string = ",\"timestamp\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float64(float64(in.TimeStamp))
	}
	{
		const prefix string = ",\"plugin\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Plugin))
	}
	{
		const prefix string = ",\"instance\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Instance))
	}
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"value\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Value == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Value {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.Float64(float64(v6))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"vltags\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.VlTags))
	}
	{
		const prefix string = ",\"message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v TargetPacket) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1bfd9a9bEncodeGithubComKexirongMonitorCommonPacketparse1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TargetPacket) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1bfd9a9bEncodeGithubComKexirongMonitorCommonPacketparse1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TargetPacket) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1bfd9a9bDecodeGithubComKexirongMonitorCommonPacketparse1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TargetPacket) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1bfd9a9bDecodeGithubComKexirongMonitorCommonPacketparse1(l, v)
}
