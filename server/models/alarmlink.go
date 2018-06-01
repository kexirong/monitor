// Package models contains the types for schema 'monitor'.
package models

import (
	"errors"
)

// AlarmLink represents a row from 'monitor.alarm_link'.
type AlarmLink struct {
	AlarmName string `json:"alarm_name"` // alarm_name
	List      string `json:"list"`       // list
	Type      Type   `json:"type"`       // type
	Channel   int    `json:"channel"`    // channel

	// xo fields
	_exists, _deleted bool `json:"-"`
}

// Exists determines if the AlarmLink exists in the database.
func (al *AlarmLink) Exists() bool {
	return al._exists
}

// Deleted provides information if the AlarmLink has been deleted from the database.
func (al *AlarmLink) Deleted() bool {
	return al._deleted
}

// Insert inserts the AlarmLink to the database.
func (al *AlarmLink) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if al._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO monitor.alarm_link (` +
		`alarm_name, list, type, channel` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, al.AlarmName, al.List, al.Type, al.Channel)
	_, err = db.Exec(sqlstr, al.AlarmName, al.List, al.Type, al.Channel)
	if err != nil {
		return err
	}

	// set existence
	al._exists = true

	return nil
}

// Update updates the AlarmLink in the database.
func (al *AlarmLink) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !al._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if al._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE monitor.alarm_link SET ` +
		`list = ?, type = ?, channel = ?` +
		` WHERE alarm_name = ?`

	// run query
	XOLog(sqlstr, al.List, al.Type, al.Channel, al.AlarmName)
	_, err = db.Exec(sqlstr, al.List, al.Type, al.Channel, al.AlarmName)
	return err
}

// Save saves the AlarmLink to the database.
func (al *AlarmLink) Save(db XODB) error {
	if al.Exists() {
		return al.Update(db)
	}

	return al.Insert(db)
}

// Delete deletes the AlarmLink from the database.
func (al *AlarmLink) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !al._exists {
		return nil
	}

	// if deleted, bail
	if al._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM monitor.alarm_link WHERE alarm_name = ?`

	// run query
	XOLog(sqlstr, al.AlarmName)
	_, err = db.Exec(sqlstr, al.AlarmName)
	if err != nil {
		return err
	}

	// set deleted
	al._deleted = true

	return nil
}

// AlarmLinkByAlarmName retrieves a row from 'monitor.alarm_link' as a AlarmLink.
//
// Generated from index 'alarm_link_alarm_name_pkey'.
func AlarmLinkByAlarmName(db XODB, alarmName string) (*AlarmLink, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`alarm_name, list, type, channel ` +
		`FROM monitor.alarm_link ` +
		`WHERE alarm_name = ?`

	// run query
	XOLog(sqlstr, alarmName)
	al := AlarmLink{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, alarmName).Scan(&al.AlarmName, &al.List, &al.Type, &al.Channel)
	if err != nil {
		return nil, err
	}

	return &al, nil
}

// AlarmLinksByChannel retrieves a row from 'monitor.alarm_link' as a AlarmLink.
//
// Generated from index 'channel'.
func AlarmLinksByChannel(db XODB, channel int) ([]*AlarmLink, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`alarm_name, list, type, channel ` +
		`FROM monitor.alarm_link ` +
		`WHERE channel = ?`

	// run query
	XOLog(sqlstr, channel)
	q, err := db.Query(sqlstr, channel)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*AlarmLink{}
	for q.Next() {
		al := AlarmLink{
			_exists: true,
		}

		// scan
		err = q.Scan(&al.AlarmName, &al.List, &al.Type, &al.Channel)
		if err != nil {
			return nil, err
		}

		res = append(res, &al)
	}

	return res, nil
}
