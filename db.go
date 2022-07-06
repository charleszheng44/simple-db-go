package main

import (
	"fmt"
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
	cols    []string
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
	db.Lock()
	defer db.Unlock()
	t, exist := db.tables[is.table]
	if !exist {
		return &Result{
			err: errors.Errorf("insert into non exist table %s",
				is.table),
		}
	}

	r := &Row{
		fields: make(map[string]any),
	}

	for cn := range t.schema {
		r.fields[cn] = nil
	}

	// primary key cannot be empty
	pk, exist := is.values[t.primaryKey]
	if !exist {
		return &Result{
			err: errors.Errorf("primary key is not given"),
		}
	}

	for cn, v := range is.values {
		kind, exist := t.schema[cn]
		// column is not defined in the schema
		if !exist {
			return &Result{
				err: errors.Errorf("column(%s) not exist", cn),
			}
		}
		// the given value kind is not as defined
		tk := reflect.TypeOf(v).Kind()
		if tk != kind {
			return &Result{
				err: errors.Errorf("invalid column(%s) type: "+
					"given(%s), expect(%s)", cn, tk, kind),
			}
		}

		r.fields[cn] = v
	}
	// insert the row to the table
	t.rows[pk] = r
	return &Result{
		message: "1 ROW INSERTED",
	}
}

func (db *Database) SelectFrom(ss *SelectStatement) *Result {
	db.RLock()
	defer db.RUnlock()
	table, exist := db.tables[ss.table]
	if !exist {
		return &Result{
			err: errors.Errorf("select from non-exist table %s",
				ss.table),
		}
	}

	// TODO (charleszheng44): filter by where clause
	var rs []*Row
	if ss.where == nil {
		for _, r := range table.rows {
			row := &Row{
				fields: make(map[string]any),
			}
			for _, f := range ss.fields {
				row.fields[f] = r.fields[f]
			}
			rs = append(rs, row)
		}
	}
	return &Result{
		rows: rs,
		cols: ss.fields,
	}
}

func (db *Database) DeleteFrom(ds *DeleteStatement) *Result {
	db.Lock()
	defer db.Unlock()
	table, exist := db.tables[ds.table]
	if !exist {
		return &Result{
			err: errors.Errorf("delete from non-exist table %s",
				ds.table),
		}
	}

	// check type
	given := reflect.TypeOf(ds.where.value).Kind()
	expect := table.schema[ds.where.field]
	if given != expect {
		return &Result{
			err: errors.Errorf("given value type is invalid: "+
				"expect(%s) got(%s)", expect, given),
		}
	}

	var count int
	for pk, r := range table.rows {
		// TODO(charleszheng44): support more operator type
		if reflect.DeepEqual(
			r.fields[ds.where.field],
			ds.where.value) {
			delete(table.rows, pk)
			count++
		}
	}

	return &Result{
		message: fmt.Sprintf("%d ROWS DELETED", count),
	}

}
