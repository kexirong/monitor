// Package models contains the types for schema 'monitor'.
package models

import (
	"database/sql/driver"
	"errors"
)

// Type is the 'type' enum type from schema 'monitor'.
type Type uint16

const (
	// TypeTeam is the 'team' Type.
	TypeTeam = Type(1)

	// TypeStaff is the 'staff' Type.
	TypeStaff = Type(2)

	// TypeOther is the 'other' Type.
	TypeOther = Type(3)
)

// String returns the string value of the Type.
func (t Type) String() string {
	var enumVal string

	switch t {
	case TypeTeam:
		enumVal = "team"

	case TypeStaff:
		enumVal = "staff"

	case TypeOther:
		enumVal = "other"
	}

	return enumVal
}

// MarshalText marshals Type into text.
func (t Type) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText unmarshals Type from text.
func (t *Type) UnmarshalText(text []byte) error {
	switch string(text) {
	case "team":
		*t = TypeTeam

	case "staff":
		*t = TypeStaff

	case "other":
		*t = TypeOther

	default:
		return errors.New("invalid Type")
	}

	return nil
}

// Value satisfies the sql/driver.Valuer interface for Type.
func (t Type) Value() (driver.Value, error) {
	return t.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for Type.
func (t *Type) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid Type")
	}

	return t.UnmarshalText(buf)
}
