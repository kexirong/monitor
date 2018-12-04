// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

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

func easyjsonAecc8163DecodeGithubComKexirongMonitorServerModels(in *jlexer.Lexer, out *PluginConfig) {
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
		case "id":
			out.ID = int64(in.Int64())
		case "host_ip":
			out.HostIP = string(in.String())
		case "host_name":
			out.HostName = string(in.String())
		case "plugin_name":
			out.PluginName = string(in.String())
		case "interval":
			out.Interval = int(in.Int())
		case "timeout":
			out.Timeout = int(in.Int())
		case "updated_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.UpdatedAt).UnmarshalJSON(data))
			}
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
func easyjsonAecc8163EncodeGithubComKexirongMonitorServerModels(out *jwriter.Writer, in PluginConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"host_ip\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.HostIP))
	}
	{
		const prefix string = ",\"host_name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.HostName))
	}
	{
		const prefix string = ",\"plugin_name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.PluginName))
	}
	{
		const prefix string = ",\"interval\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Interval))
	}
	{
		const prefix string = ",\"timeout\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Timeout))
	}
	{
		const prefix string = ",\"updated_at\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.UpdatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PluginConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonAecc8163EncodeGithubComKexirongMonitorServerModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PluginConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonAecc8163EncodeGithubComKexirongMonitorServerModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PluginConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonAecc8163DecodeGithubComKexirongMonitorServerModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PluginConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonAecc8163DecodeGithubComKexirongMonitorServerModels(l, v)
}
