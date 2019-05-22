package gomaxcompute

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	a := assert.New(t)

	db, err := sql.Open("maxcompute", cfg4test.FormatDSN())
	a.NoError(err)

	const sql = `SELECT * from gomaxcompute_test;`
	rows, err := db.Query(sql)
	defer rows.Close()
	a.NoError(err)

	cols, _ := rows.Columns()
	cnt := len(cols)
	values := make([]interface{}, cnt)
	row := make([]interface{}, cnt)
	for rows.Next() {
		cts, err := rows.ColumnTypes()
		a.NoError(err)

		for i, ct := range cts {
			v, err := createByType(ct.ScanType())
			a.NoError(err)
			values[i] = v
		}
		err = rows.Scan(values...)
		a.NoError(err)

		for i, val := range values {
			v, err := parseVal(val)
			a.NoError(err)
			row[i] = v
		}
	}
}

func createByType(rt reflect.Type) (interface{}, error) {
	switch rt {
	case builtinString:
		return new(string), nil
	default:
		return nil, fmt.Errorf("unrecognized column scan type %v", rt)
	}
}

func parseVal(val interface{}) (interface{}, error) {
	switch v := val.(type) {
	case *string:
		return *v, nil
	default:
		return nil, fmt.Errorf("unrecogized type %v", v)
	}
}
