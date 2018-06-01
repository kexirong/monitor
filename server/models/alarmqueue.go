// Package models contains the types for schema 'monitor'.
package models

import (
	"errors"
	"fmt"
	"time"
)

// AlarmQueue represents a row from 'monitor.alarm_queue'.
type AlarmQueue struct {
	ID        int64     `json:"id"`         // id
	HostName  string    `json:"host_name"`  // host_name
	AlarmName string    `json:"alarm_name"` // alarm_name
	Alarmele  string    `json:"alarmele"`   // alarmele
	Value     float64   `json:"value"`      // value
	Message   string    `json:"message"`    // message
	HandleMan string    `json:"handle_man"` // handle_man
	Stat      int       `json:"stat"`       // stat
	Level     Level     `json:"level"`      // level
	CreatedAt time.Time `json:"created_at"` // created_at
	UpdatedAt time.Time `json:"updated_at"` // updated_at

	// xo fields
	_exists, _deleted bool
}

func (aq *AlarmQueue) String() string {
	return fmt.Sprintf("[%s]seq: %d, Time: %s,HostName: %s,  Plugin: %s, Instance: %s, Value: %g, Message: %s",
		aq.Level, aq.ID, aq.CreatedAt, aq.HostName, aq.AlarmName, aq.Alarmele, aq.Value, aq.Message)
}

// Exists determines if the AlarmQueue exists in the database.
func (aq *AlarmQueue) Exists() bool {
	return aq._exists
}

// Deleted provides information if the AlarmQueue has been deleted from the database.
func (aq *AlarmQueue) Deleted() bool {
	return aq._deleted
}

// Insert inserts the AlarmQueue to the database.
func (aq *AlarmQueue) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if aq._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO monitor.alarm_queue (` +
		`host_name, alarm_name, alarmele, value, message, handle_man, stat, level, created_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?,  ?` +
		`)`

	// run query
	XOLog(sqlstr, aq.HostName, aq.AlarmName, aq.Alarmele, aq.Value, aq.Message, aq.HandleMan, aq.Stat, aq.Level, aq.CreatedAt)
	res, err := db.Exec(sqlstr, aq.HostName, aq.AlarmName, aq.Alarmele, aq.Value, aq.Message, aq.HandleMan, aq.Stat, aq.Level, aq.CreatedAt)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	aq.ID = int64(id)
	aq._exists = true

	return nil
}

// Update updates the AlarmQueue in the database.
func (aq *AlarmQueue) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !aq._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if aq._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE monitor.alarm_queue SET ` +
		`host_name = ?, alarm_name = ?, alarmele = ?, value = ?, message = ?, handle_man = ?, stat = ?, level = ? ` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, aq.HostName, aq.AlarmName, aq.Alarmele, aq.Value, aq.Message, aq.HandleMan, aq.Stat, aq.Level, aq.ID)
	_, err = db.Exec(sqlstr, aq.HostName, aq.AlarmName, aq.Alarmele, aq.Value, aq.Message, aq.HandleMan, aq.Stat, aq.Level, aq.ID)
	return err
}

// Save saves the AlarmQueue to the database.
func (aq *AlarmQueue) Save(db XODB) error {
	if aq.Exists() {
		return aq.Update(db)
	}

	return aq.Insert(db)
}

// Delete deletes the AlarmQueue from the database.
func (aq *AlarmQueue) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !aq._exists {
		return nil
	}

	// if deleted, bail
	if aq._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM monitor.alarm_queue WHERE id = ?`

	// run query
	XOLog(sqlstr, aq.ID)
	_, err = db.Exec(sqlstr, aq.ID)
	if err != nil {
		return err
	}

	// set deleted
	aq._deleted = true

	return nil
}

// AlarmQueueByID retrieves a row from 'monitor.alarm_queue' as a AlarmQueue.
//
// Generated from index 'alarm_queue_id_pkey'.
func AlarmQueueByID(db XODB, id int64) (*AlarmQueue, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, host_name, alarm_name, alarmele, value, message, handle_man, stat, level, created_at, updated_at ` +
		`FROM monitor.alarm_queue ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	aq := AlarmQueue{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&aq.ID, &aq.HostName, &aq.AlarmName, &aq.Alarmele, &aq.Value, &aq.Message, &aq.HandleMan, &aq.Stat, &aq.Level, &aq.CreatedAt, &aq.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &aq, nil
}

func AlarmQueueByStat(db XODB, stat int) ([]*AlarmQueue, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, host_name, alarm_name, alarmele, value, message, handle_man, stat, level, created_at, updated_at ` +
		`FROM monitor.alarm_queue where stat = ? `

	// run query
	XOLog(sqlstr, stat)
	q, err := db.Query(sqlstr, stat)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*AlarmQueue{}
	for q.Next() {
		aq := AlarmQueue{
			_exists: true,
		}

		// scan
		err = q.Scan(&aq.ID, &aq.HostName, &aq.AlarmName, &aq.Alarmele, &aq.Value, &aq.Message, &aq.HandleMan, &aq.Stat, &aq.Level, &aq.CreatedAt, &aq.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &aq)
	}

	return res, nil
}
