package main

import (
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

type Database struct {
	sync.RWMutex
	tables map[string]*Table
}

func NewDatabase() *Database {
	return &Database{
		tables: make(map[string]*Table),
	}
}

type Result struct {
	err     error
	rows    []*Row
	message string
}

func (db *Database) Interpret(sts any) *Result {
	switch s := sts.(type) {
	case *CreateStatement:
		return db.CreateTable(s)
	case *SelectStatement:
		return db.SelectFrom(s)
	case *InsertStatement:
		return db.InsertInto(s)
	case *DeleteStatement:
		return db.DeleteFrom(s)
	default:
		return &Result{
			err: errors.Errorf("unsupported statement %v",
				reflect.TypeOf(sts).Kind()),
		}
	}
}

func (db *Database) CreateTable(cs *CreateStatement) *Result {
	if _, exist := cs.schema[cs.primaryKey]; !exist {
		return &Result{
			err: errors.New("primary key not defined"),
		}
	}
	db.Lock()
	defer db.Unlock()
	if _, exist := db.tables[cs.table]; exist {
		return &Result{
			err: errors.Errorf("table %s already exists"),
		}
	}
	db.tables[cs.table] = NewTable(cs.primaryKey, cs.schema)
	return &Result{
		message: "TABLE CREATED",
	}
}

func (db *Database) DeleteTable(ds *DeleteStatement) *Result {
	db.Lock()
	defer db.Unlock()
	if _, exist := db.tables[ds.table]; !exist {
		return &Result{
			err: errors.Errorf("delete non exist table %s", ds.table),
		}
	}
	delete(db.tables, ds.table)
	return nil
}

func (db *Database) InsertInto(is *InsertStatement) *Result {
	panic("NOT IMPLEMENT YET")
}

func (db *Database) SelectFrom(ss *SelectStatement) *Result {
	/*
		db.RLock()
		defer db.RUnlock()
		table := db.tables[ss.table]
		if table == nil {
			return &Result{
				err: errors.Errorf("select from non-exist table %s", ss.table),
			}
		}

		var ret []*Row

		if ss.where == nil {
			for _, r := range table.rows {
				row := &Row{
					fields: make(map[string]any),
				}
				for _, f := range ss.fields {
					r.fields[f]
				}
			}
		}
	*/
	panic("NOT IMPLEMENT YET")

}

func (db *Database) DeleteFrom(ds *DeleteStatement) *Result {
	panic("NOT IMPLEMENT YET")
}
