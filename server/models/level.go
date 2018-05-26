// Package models contains the types for schema 'monitor'.
package models

import (
	"database/sql/driver"
	"errors"
)

// Level is the 'level' enum type from schema 'monitor'.
type Level uint16

const (
	// LevelLevel1 is the 'level1' Level.
	LevelLevel1 = Level(1)

	// LevelLevel2 is the 'level2' Level.
	LevelLevel2 = Level(2)

	// LevelLevel3 is the 'level3' Level.
	LevelLevel3 = Level(3)
)

// String returns the string value of the Level.
func (l Level) String() string {
	var enumVal string

	switch l {
	case LevelLevel1:
		enumVal = "level1"

	case LevelLevel2:
		enumVal = "level2"

	case LevelLevel3:
		enumVal = "level3"
	}

	return enumVal
}

// MarshalText marshals Level into text.
func (l Level) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// UnmarshalText unmarshals Level from text.
func (l *Level) UnmarshalText(text []byte) error {
	switch string(text) {
	case "level1":
		*l = LevelLevel1

	case "level2":
		*l = LevelLevel2

	case "level3":
		*l = LevelLevel3

	default:
		return errors.New("invalid Level")
	}

	return nil
}

// Value satisfies the sql/driver.Valuer interface for Level.
func (l Level) Value() (driver.Value, error) {
	return l.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for Level.
func (l *Level) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid Level")
	}

	return l.UnmarshalText(buf)
}
