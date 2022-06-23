package main

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

func stringToKind(kindStr string) (reflect.Kind, error) {
	kindStr = strings.ToLower(kindStr)
	switch kindStr {
	case "integer":
		return reflect.Int, nil
	case "float":
		return reflect.Float64, nil
	case "string":
		return reflect.String, nil
	case "boolean":
		return reflect.Bool, nil
	default:
		return reflect.Kind(^uint(0)),
			errors.Errorf("unsupported kind %s", kindStr)
	}
}

type Table struct {
	primaryKey string
	schema     map[string]reflect.Kind
	rows       map[any]*Row
}

type Row struct {
	fields map[string]any
}

func NewTable(pk string, schema map[string]reflect.Kind) *Table {
	return &Table{
		primaryKey: pk,
		schema:     schema,
		rows:       make(map[any]*Row),
	}
}
