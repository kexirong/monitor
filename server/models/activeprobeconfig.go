// Package models contains the types for schema 'monitor'.
package models

import (
	"errors"
	"time"
)

// ActiveProbeConfig represents a row from 'monitor.active_probe_config'.
type ActiveProbeConfig struct {
	ID            int       `json:"id"`              // id
	ActiveProbeID int64     `json:"active_probe_id"` // active_probe_id
	Target        string    `json:"target"`          // target
	Arg1          string    `json:"arg1"`            // arg1
	Arg2          string    `json:"arg2"`            // arg2
	UpdatedAt     time.Time `json:"updated_at"`      // updated_at

	// xo fields
	_exists, _deleted bool `json:"-"`
}

// Exists determines if the ActiveProbeConfig exists in the database.
func (apc *ActiveProbeConfig) Exists() bool {
	return apc._exists
}

// Deleted provides information if the ActiveProbeConfig has been deleted from the database.
func (apc *ActiveProbeConfig) Deleted() bool {
	return apc._deleted
}

// Insert inserts the ActiveProbeConfig to the database.
func (apc *ActiveProbeConfig) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if apc._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO monitor.active_probe_config (` +
		`active_probe_id, target, arg1, arg2` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, apc.ActiveProbeID, apc.Target, apc.Arg1, apc.Arg2)
	res, err := db.Exec(sqlstr, apc.ActiveProbeID, apc.Target, apc.Arg1, apc.Arg2)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	apc.ID = int(id)
	apc._exists = true

	return nil
}

// Update updates the ActiveProbeConfig in the database.
func (apc *ActiveProbeConfig) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !apc._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if apc._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE monitor.active_probe_config SET ` +
		`active_probe_id = ?, target = ?, arg1 = ?, arg2 = ?  ` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, apc.ActiveProbeID, apc.Target, apc.Arg1, apc.Arg2, apc.UpdatedAt, apc.ID)
	_, err = db.Exec(sqlstr, apc.ActiveProbeID, apc.Target, apc.Arg1, apc.Arg2, apc.ID)
	return err
}

// Save saves the ActiveProbeConfig to the database.
func (apc *ActiveProbeConfig) Save(db XODB) error {
	if apc.Exists() {
		return apc.Update(db)
	}

	return apc.Insert(db)
}

// Delete deletes the ActiveProbeConfig from the database.
func (apc *ActiveProbeConfig) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !apc._exists {
		return nil
	}

	// if deleted, bail
	if apc._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM monitor.active_probe_config WHERE id = ?`

	// run query
	XOLog(sqlstr, apc.ID)
	_, err = db.Exec(sqlstr, apc.ID)
	if err != nil {
		return err
	}

	// set deleted
	apc._deleted = true

	return nil
}

// ActiveProbe returns the ActiveProbe associated with the ActiveProbeConfig's ActiveProbeID (active_probe_id).
//
// Generated from foreign key 'active_probe_config_id'.
func (apc *ActiveProbeConfig) ActiveProbe(db XODB) (*ActiveProbe, error) {
	return ActiveProbeByID(db, apc.ActiveProbeID)
}

// ActiveProbeConfigsByActiveProbeID retrieves a row from 'monitor.active_probe_config' as a ActiveProbeConfig.
//
// Generated from index 'IDX_ActiveProbeConfig_id'.
func ActiveProbeConfigsByActiveProbeID(db XODB, activeProbeID int64) ([]*ActiveProbeConfig, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, active_probe_id, target, arg1, arg2, updated_at ` +
		`FROM monitor.active_probe_config ` +
		`WHERE active_probe_id = ?`

	// run query
	XOLog(sqlstr, activeProbeID)
	q, err := db.Query(sqlstr, activeProbeID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*ActiveProbeConfig{}
	for q.Next() {
		apc := ActiveProbeConfig{
			_exists: true,
		}

		// scan
		err = q.Scan(&apc.ID, &apc.ActiveProbeID, &apc.Target, &apc.Arg1, &apc.Arg2, &apc.UpdatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &apc)
	}

	return res, nil
}

// ActiveProbeConfigByActiveProbeIDTarget retrieves a row from 'monitor.active_probe_config' as a ActiveProbeConfig.
//
// Generated from index 'UNIQUE_ActiveProbeConfig_id_target'.
func ActiveProbeConfigByActiveProbeIDTarget(db XODB, activeProbeID int64, target string) (*ActiveProbeConfig, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, active_probe_id, target, arg1, arg2, updated_at ` +
		`FROM monitor.active_probe_config ` +
		`WHERE active_probe_id = ? AND target = ? `

	// run query
	XOLog(sqlstr, activeProbeID, target)
	apc := ActiveProbeConfig{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, activeProbeID, target).Scan(&apc.ID, &apc.ActiveProbeID, &apc.Target, &apc.Arg1, &apc.Arg2, &apc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &apc, nil
}

// ActiveProbeConfigByID retrieves a row from 'monitor.active_probe_config' as a ActiveProbeConfig.
//
// Generated from index 'active_probe_config_id_pkey'.
func ActiveProbeConfigByID(db XODB, id int) (*ActiveProbeConfig, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, active_probe_id, target, arg1, arg2, updated_at ` +
		`FROM monitor.active_probe_config ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	apc := ActiveProbeConfig{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&apc.ID, &apc.ActiveProbeID, &apc.Target, &apc.Arg1, &apc.Arg2, &apc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &apc, nil
}
