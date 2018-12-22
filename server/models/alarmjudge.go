// Package models contains the types for schema 'monitor'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
)

// AlarmJudge represents a row from 'monitor.alarm_judge'.
type AlarmJudge struct {
	ID          int64  `json:"id"`           // id
	AnchorPoint string `json:"anchor_point"` // anchor_point
	Express     string `json:"express"`      // express
	Level       Level  `json:"level"`        // level

	// xo fields
	_exists, _deleted bool
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

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO monitor.alarm_judge (` +
		`anchor_point, express, level` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, aj.AnchorPoint, aj.Express, aj.Level)
	res, err := db.Exec(sqlstr, aj.AnchorPoint, aj.Express, aj.Level)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	aj.ID = int64(id)
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

	// sql query
	const sqlstr = `UPDATE monitor.alarm_judge SET ` +
		`anchor_point = ?, express = ?, level = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, aj.AnchorPoint, aj.Express, aj.Level, aj.ID)
	_, err = db.Exec(sqlstr, aj.AnchorPoint, aj.Express, aj.Level, aj.ID)
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

	// sql query
	const sqlstr = `DELETE FROM monitor.alarm_judge WHERE id = ?`

	// run query
	XOLog(sqlstr, aj.ID)
	_, err = db.Exec(sqlstr, aj.ID)
	if err != nil {
		return err
	}

	// set deleted
	aj._deleted = true

	return nil
}

func AlarmJudgesAll(db XODB) ([]*AlarmJudge, error) {
	const sqlstr = `SELECT ` +
		`id, anchor_point, express, level ` +
		`FROM monitor.alarm_judge `
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*AlarmJudge{}
	for q.Next() {
		aj := AlarmJudge{
			_exists: true,
		}

		// scan
		err = q.Scan(&aj.ID, &aj.AnchorPoint, &aj.Express, &aj.Level)
		if err != nil {
			return nil, err
		}

		res = append(res, &aj)
	}

	return res, nil
} // AlarmJudgeByAnchorPointExpress retrieves a row from 'monitor.alarm_judge' as a AlarmJudge.
//
// Generated from index 'UNI_AlarmJudge_AnchorPoint_Express'.
func AlarmJudgeByAnchorPointExpress(db XODB, anchorPoint string, express string) (*AlarmJudge, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, anchor_point, express, level ` +
		`FROM monitor.alarm_judge ` +
		`WHERE anchor_point = ? AND express = ?`

	// run query
	XOLog(sqlstr, anchorPoint, express)
	aj := AlarmJudge{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, anchorPoint, express).Scan(&aj.ID, &aj.AnchorPoint, &aj.Express, &aj.Level)
	if err != nil {
		return nil, err
	}

	return &aj, nil
}

// AlarmJudgeByID retrieves a row from 'monitor.alarm_judge' as a AlarmJudge.
//
// Generated from index 'alarm_judge_id_pkey'.
func AlarmJudgeByID(db XODB, id int64) (*AlarmJudge, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, anchor_point, express, level ` +
		`FROM monitor.alarm_judge ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	aj := AlarmJudge{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&aj.ID, &aj.AnchorPoint, &aj.Express, &aj.Level)
	if err != nil {
		return nil, err
	}

	return &aj, nil
}
