package models

import (
	"database/sql"
	"strconv"
)

//NullInt64 Int is an nullable int64.
type NullInt64 struct {
	sql.NullInt64
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *NullInt64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	i.Valid = err == nil
	return err
}

// MarshalText implements encoding.TextMarshaler.
func (i NullInt64) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}
