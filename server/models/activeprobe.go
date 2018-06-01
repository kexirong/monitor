// Package models contains the types for schema 'monitor'.
package models

import (
	"errors"
	"time"
)

// ActiveProbe represents a row from 'monitor.active_probe'.
type ActiveProbe struct {
	ID         int64     `json:"id"`          // id
	PluginName string    `json:"plugin_name"` // plugin_name
	HostName   string    `json:"host_name"`   // host_name
	IP         string    `json:"ip"`          // ip
	Interval   int       `json:"interval"`    // interval
	UpdatedAt  time.Time `json:"updated_at"`  // updated_at

	// xo fields
	_exists, _deleted bool `json:"-"`
}

// Exists determines if the ActiveProbe exists in the database.
func (ap *ActiveProbe) Exists() bool {
	return ap._exists
}

// Deleted provides information if the ActiveProbe has been deleted from the database.
func (ap *ActiveProbe) Deleted() bool {
	return ap._deleted
}

// Insert inserts the ActiveProbe to the database.
func (ap *ActiveProbe) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if ap._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO monitor.active_probe (` +
		`plugin_name, host_name, ip, 'interval'` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, ap.PluginName, ap.HostName, ap.IP, ap.Interval)
	res, err := db.Exec(sqlstr, ap.PluginName, ap.HostName, ap.IP, ap.Interval)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	ap.ID = int64(id)
	ap._exists = true

	return nil
}

// Update updates the ActiveProbe in the database.
func (ap *ActiveProbe) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !ap._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if ap._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE monitor.active_probe SET ` +
		`plugin_name = ?, host_name = ?, ip = ?, 'interval' = ? ` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, ap.PluginName, ap.HostName, ap.IP, ap.Interval, ap.ID)
	_, err = db.Exec(sqlstr, ap.PluginName, ap.HostName, ap.IP, ap.Interval, ap.ID)
	return err
}

// Save saves the ActiveProbe to the database.
func (ap *ActiveProbe) Save(db XODB) error {
	if ap.Exists() {
		return ap.Update(db)
	}

	return ap.Insert(db)
}

// Delete deletes the ActiveProbe from the database.
func (ap *ActiveProbe) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !ap._exists {
		return nil
	}

	// if deleted, bail
	if ap._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM monitor.active_probe WHERE id = ?`

	// run query
	XOLog(sqlstr, ap.ID)
	_, err = db.Exec(sqlstr, ap.ID)
	if err != nil {
		return err
	}

	// set deleted
	ap._deleted = true

	return nil
}

// ActiveProbeByPluginNameHostName retrieves a row from 'monitor.active_probe' as a ActiveProbe.
//
// Generated from index 'UNIQUE_ActiveProbeConfig_name_hostname'.
func ActiveProbeByPluginNameHostName(db XODB, pluginName string, hostName string) (*ActiveProbe, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, plugin_name, host_name, ip, 'interval', updated_at ` +
		`FROM monitor.active_probe ` +
		`WHERE plugin_name = ? AND host_name = ?`

	// run query
	XOLog(sqlstr, pluginName, hostName)
	ap := ActiveProbe{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, pluginName, hostName).Scan(&ap.ID, &ap.PluginName, &ap.HostName, &ap.IP, &ap.Interval, &ap.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &ap, nil
}

// ActiveProbeByID retrieves a row from 'monitor.active_probe' as a ActiveProbe.
//
// Generated from index 'active_probe_id_pkey'.
func ActiveProbeByID(db XODB, id int64) (*ActiveProbe, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, plugin_name, host_name, ip, 'interval', updated_at ` +
		`FROM monitor.active_probe ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	ap := ActiveProbe{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&ap.ID, &ap.PluginName, &ap.HostName, &ap.IP, &ap.Interval, &ap.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &ap, nil
}

func ActiveProbeAll(db XODB) ([]*ActiveProbe, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		"id, plugin_name, host_name, ip, `interval`, updated_at " +
		`FROM monitor.active_probe `

	// run query
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*ActiveProbe{}
	for q.Next() {
		ap := ActiveProbe{
			_exists: true,
		}

		// scan
		err = q.Scan(&ap.ID, &ap.PluginName, &ap.HostName, &ap.IP, &ap.Interval, &ap.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &ap)
	}

	return res, nil
}
