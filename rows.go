package gomaxcompute

import (
	"bytes"
	"database/sql/driver"
	"encoding/csv"
	"reflect"
)

var builtinString = reflect.TypeOf(string(""))

type odpsRows struct {
	meta      *resultMeta
	reader    *csv.Reader
	header    []string
	headerLen int
}

func (rs *odpsRows) Close() error {
	return nil
}

func (rs *odpsRows) Columns() []string {
	return rs.header
}

// Notice: odps using `\N` to denote nil, even for not-string field
func (rs *odpsRows) Next(dst []driver.Value) error {
	records, err := rs.reader.Read()
	if err != nil {
		return err
	}
	for i := 0; i < rs.headerLen; i++ {
		dst[i] = records[i]
	}
	return nil
}

func (rs *odpsRows) ColumnTypeScanType(i int) reflect.Type {
	return builtinString
}

func (rs *odpsRows) ColumnTypeDatabaseTypeName(i int) string {
	return rs.meta.Schema.Columns[i].Type
}

func newRows(m *resultMeta, res string) (*odpsRows, error) {
	rd := csv.NewReader(bytes.NewBufferString(res))
	// hr equas to [m.Schema.Columns.Name]
	hr, err := rd.Read()
	if err != nil {
		return nil, err
	}
	return &odpsRows{meta: m, reader: rd, header: hr, headerLen: len(hr)}, nil
}
