// Package models contains the types for schema 'monitor'.
package models

import (
	"errors"
)

// PluginConfig represents a row from 'monitor.plugin_config'.
type PluginConfig struct {
	ID         int    `json:"id"`          // id
	HostIP     string `json:"host_ip"`     // host_ip
	HostName   string `json:"host_name"`   // host_name
	PluginName string `json:"plugin_name"` // plugin_name
	Interval   int    `json:"interval"`    // interval
	Timeout    int    `json:"timeout"`     // timeout

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the PluginConfig exists in the database.
func (pc *PluginConfig) Exists() bool {
	return pc._exists
}

// Deleted provides information if the PluginConfig has been deleted from the database.
func (pc *PluginConfig) Deleted() bool {
	return pc._deleted
}

// Insert inserts the PluginConfig to the database.
func (pc *PluginConfig) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if pc._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO monitor.plugin_config (` +
		`host_ip, host_name, plugin_name, interval, timeout` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, pc.HostIP, pc.HostName, pc.PluginName, pc.Interval, pc.Timeout)
	res, err := db.Exec(sqlstr, pc.HostIP, pc.HostName, pc.PluginName, pc.Interval, pc.Timeout)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	pc.ID = int(id)
	pc._exists = true

	return nil
}

// Update updates the PluginConfig in the database.
func (pc *PluginConfig) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !pc._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if pc._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE monitor.plugin_config SET ` +
		`host_ip = ?, host_name = ?, plugin_name = ?, interval = ?, timeout = ?` +
		` WHERE id = ?`

	// run query
	XOLog(sqlstr, pc.HostIP, pc.HostName, pc.PluginName, pc.Interval, pc.Timeout, pc.ID)
	_, err = db.Exec(sqlstr, pc.HostIP, pc.HostName, pc.PluginName, pc.Interval, pc.Timeout, pc.ID)
	return err
}

// Save saves the PluginConfig to the database.
func (pc *PluginConfig) Save(db XODB) error {
	if pc.Exists() {
		return pc.Update(db)
	}

	return pc.Insert(db)
}

// Delete deletes the PluginConfig from the database.
func (pc *PluginConfig) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !pc._exists {
		return nil
	}

	// if deleted, bail
	if pc._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM monitor.plugin_config WHERE id = ?`

	// run query
	XOLog(sqlstr, pc.ID)
	_, err = db.Exec(sqlstr, pc.ID)
	if err != nil {
		return err
	}

	// set deleted
	pc._deleted = true

	return nil
}

// Plugin returns the Plugin associated with the PluginConfig's PluginName (plugin_name).
//
// Generated from foreign key 'plugin_config_ibfk_1'.
func (pc *PluginConfig) Plugin(db XODB) (*Plugin, error) {
	return PluginByPluginName(db, pc.PluginName)
}

// PluginConfigByPluginNameHostName retrieves a row from 'monitor.plugin_config' as a PluginConfig.
//
// Generated from index 'UNIQUE_pluginConfig_pluginname_hostnamae'.
func PluginConfigByPluginNameHostName(db XODB, pluginName string, hostName string) (*PluginConfig, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, host_ip, host_name, plugin_name, interval, timeout ` +
		`FROM monitor.plugin_config ` +
		`WHERE plugin_name = ? AND host_name = ?`

	// run query
	XOLog(sqlstr, pluginName, hostName)
	pc := PluginConfig{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, pluginName, hostName).Scan(&pc.ID, &pc.HostIP, &pc.HostName, &pc.PluginName, &pc.Interval, &pc.Timeout)
	if err != nil {
		return nil, err
	}

	return &pc, nil
}

// PluginConfigByID retrieves a row from 'monitor.plugin_config' as a PluginConfig.
//
// Generated from index 'plugin_config_id_pkey'.
func PluginConfigByID(db XODB, id int) (*PluginConfig, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, host_ip, host_name, plugin_name, interval, timeout ` +
		`FROM monitor.plugin_config ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	pc := PluginConfig{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&pc.ID, &pc.HostIP, &pc.HostName, &pc.PluginName, &pc.Interval, &pc.Timeout)
	if err != nil {
		return nil, err
	}

	return &pc, nil
}
