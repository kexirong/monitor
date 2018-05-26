// Package models contains the types for schema 'monitor'.
package models

import (
	"database/sql/driver"
	"errors"
)

// Ajtype is the 'ajtype' enum type from schema 'monitor'.
type Ajtype uint16

const (
	// AjtypeLe is the 'le' Ajtype.
	AjtypeLe = Ajtype(1)

	// AjtypeNe is the 'ne' Ajtype.
	AjtypeNe = Ajtype(2)

	// AjtypeGe is the 'ge' Ajtype.
	AjtypeGe = Ajtype(3)
)

// String returns the string value of the Ajtype.
func (a Ajtype) String() string {
	var enumVal string

	switch a {
	case AjtypeLe:
		enumVal = "le"

	case AjtypeNe:
		enumVal = "ne"

	case AjtypeGe:
		enumVal = "ge"
	}

	return enumVal
}

// MarshalText marshals Ajtype into text.
func (a Ajtype) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

// UnmarshalText unmarshals Ajtype from text.
func (a *Ajtype) UnmarshalText(text []byte) error {
	switch string(text) {
	case "le":
		*a = AjtypeLe

	case "ne":
		*a = AjtypeNe

	case "ge":
		*a = AjtypeGe

	default:
		return errors.New("invalid Ajtype")
	}

	return nil
}

// Value satisfies the sql/driver.Valuer interface for Ajtype.
func (a Ajtype) Value() (driver.Value, error) {
	return a.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for Ajtype.
func (a *Ajtype) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid Ajtype")
	}

	return a.UnmarshalText(buf)
}
