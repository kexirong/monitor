// Package models contains the types for schema 'monitor'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
	"time"
)

// Plugin represents a row from 'monitor.plugin'.
type Plugin struct {
	PluginName string    `json:"plugin_name"` // plugin_name
	PluginType string    `json:"plugin_type"` // plugin_type
	FileName   string    `json:"file_name"`   // file_name
	Comment    string    `json:"comment"`     // comment
	CreatedAt  time.Time `json:"created_at"`  // created_at

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Plugin exists in the database.
func (p *Plugin) Exists() bool {
	return p._exists
}

// Deleted provides information if the Plugin has been deleted from the database.
func (p *Plugin) Deleted() bool {
	return p._deleted
}

// Insert inserts the Plugin to the database.
func (p *Plugin) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if p._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO monitor.plugin (` +
		`plugin_name, plugin_type, file_name, comment, created_at` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, p.PluginName, p.PluginType, p.FileName, p.Comment, p.CreatedAt)
	_, err = db.Exec(sqlstr, p.PluginName, p.PluginType, p.FileName, p.Comment, p.CreatedAt)
	if err != nil {
		return err
	}

	// set existence
	p._exists = true

	return nil
}

// Update updates the Plugin in the database.
func (p *Plugin) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !p._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if p._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE monitor.plugin SET ` +
		`plugin_type = ?, file_name = ?, comment = ?, created_at = ?` +
		` WHERE plugin_name = ?`

	// run query
	XOLog(sqlstr, p.PluginType, p.FileName, p.Comment, p.CreatedAt, p.PluginName)
	_, err = db.Exec(sqlstr, p.PluginType, p.FileName, p.Comment, p.CreatedAt, p.PluginName)
	return err
}

// Save saves the Plugin to the database.
func (p *Plugin) Save(db XODB) error {
	if p.Exists() {
		return p.Update(db)
	}

	return p.Insert(db)
}

// Delete deletes the Plugin from the database.
func (p *Plugin) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !p._exists {
		return nil
	}

	// if deleted, bail
	if p._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM monitor.plugin WHERE plugin_name = ?`

	// run query
	XOLog(sqlstr, p.PluginName)
	_, err = db.Exec(sqlstr, p.PluginName)
	if err != nil {
		return err
	}

	// set deleted
	p._deleted = true

	return nil
}

func PluginsAll(db XODB) ([]*Plugin, error) {
	const sqlstr = `SELECT ` +
		`plugin_name, plugin_type, file_name, comment, created_at ` +
		`FROM monitor.plugin `
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Plugin{}
	for q.Next() {
		p := Plugin{
			_exists: true,
		}

		// scan
		err = q.Scan(&p.PluginName, &p.PluginType, &p.FileName, &p.Comment, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, &p)
	}

	return res, nil
} // PluginByPluginName retrieves a row from 'monitor.plugin' as a Plugin.
//
// Generated from index 'plugin_plugin_name_pkey'.
func PluginByPluginName(db XODB, pluginName string) (*Plugin, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`plugin_name, plugin_type, file_name, comment, created_at ` +
		`FROM monitor.plugin ` +
		`WHERE plugin_name = ?`

	// run query
	XOLog(sqlstr, pluginName)
	p := Plugin{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, pluginName).Scan(&p.PluginName, &p.PluginType, &p.FileName, &p.Comment, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
