// Package models contains the types for schema 'monitor'.
package models

import (
	"database/sql"
	"errors"
)

// AlarmJudge represents a row from 'monitor.alarm_judge'.
type AlarmJudge struct {
	AlarmName string        `json:"alarm_name"` // alarm_name
	Alarmele  string        `json:"alarmele"`   // alarmele
	Ajtype    Ajtype        `json:"ajtype"`     // ajtype
	Level1    sql.NullInt64 `json:"level1"`     // level1
	Level2    sql.NullInt64 `json:"level2"`     // level2
	Level3    sql.NullInt64 `json:"level3"`     // level3

	// xo fields
	_exists, _deleted bool `json:"-"`
}

// Exists determines if the AlarmJudge exists in the database.
func (aj *AlarmJudge) Exists() bool {
	return aj._exists
}

// Deleted provides information if the AlarmJudge has been deleted from the database.
func (aj *AlarmJudge) Deleted() bool {
	return aj._deleted
}

// Insert inserts the AlarmJudge to the database.
func (aj *AlarmJudge) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if aj._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO monitor.alarm_judge (` +
		`alarm_name, alarmele, ajtype, level1, level2, level3` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, aj.AlarmName, aj.Alarmele, aj.Ajtype, aj.Level1, aj.Level2, aj.Level3)
	_, err = db.Exec(sqlstr, aj.AlarmName, aj.Alarmele, aj.Ajtype, aj.Level1, aj.Level2, aj.Level3)
	if err != nil {
		return err
	}

	// set existence
	aj._exists = true

	return nil
}

// Update updates the AlarmJudge in the database.
func (aj *AlarmJudge) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !aj._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if aj._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query with composite primary key
	const sqlstr = `UPDATE monitor.alarm_judge SET ` +
		`ajtype = ?, level1 = ?, level2 = ?, level3 = ?` +
		` WHERE alarm_name = ? AND alarmele = ?`

	// run query
	XOLog(sqlstr, aj.Ajtype, aj.Level1, aj.Level2, aj.Level3, aj.AlarmName, aj.Alarmele)
	_, err = db.Exec(sqlstr, aj.Ajtype, aj.Level1, aj.Level2, aj.Level3, aj.AlarmName, aj.Alarmele)
	return err
}

// Save saves the AlarmJudge to the database.
func (aj *AlarmJudge) Save(db XODB) error {
	if aj.Exists() {
		return aj.Update(db)
	}

	return aj.Insert(db)
}

// Delete deletes the AlarmJudge from the database.
func (aj *AlarmJudge) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !aj._exists {
		return nil
	}

	// if deleted, bail
	if aj._deleted {
		return nil
	}

	// sql query with composite primary key
	const sqlstr = `DELETE FROM monitor.alarm_judge WHERE alarm_name = ? AND alarmele = ?`

	// run query
	XOLog(sqlstr, aj.AlarmName, aj.Alarmele)
	_, err = db.Exec(sqlstr, aj.AlarmName, aj.Alarmele)
	if err != nil {
		return err
	}

	// set deleted
	aj._deleted = true

	return nil
}

// AlarmJudgeByAlarmele retrieves a row from 'monitor.alarm_judge' as a AlarmJudge.
//
// Generated from index 'alarm_judge_alarmele_pkey'.
func AlarmJudgeByAlarmele(db XODB, alarmele string) (*AlarmJudge, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`alarm_name, alarmele, ajtype, level1, level2, level3 ` +
		`FROM monitor.alarm_judge ` +
		`WHERE alarmele = ?`

	// run query
	XOLog(sqlstr, alarmele)
	aj := AlarmJudge{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, alarmele).Scan(&aj.AlarmName, &aj.Alarmele, &aj.Ajtype, &aj.Level1, &aj.Level2, &aj.Level3)
	if err != nil {
		return nil, err
	}

	return &aj, nil
}
