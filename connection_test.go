package gomaxcompute

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	a := assert.New(t)
	db, err := sql.Open("maxcompute", cfg4test.FormatDSN())
	a.NoError(err)

	const stmt = `SELECT * FROM gomaxcompute_test LIMIT 1;`
	_, err = db.Query(stmt)
	a.NoError(err)
}

func TestQueryBase64(t *testing.T) {
	a := assert.New(t)
	db, err := sql.Open("maxcompute", cfg4test.FormatDSN())
	a.NoError(err)

	row, err := db.Query(`SELECT CAST("\001" AS string) AS a;`)
	a.NoError(err)
	for row.Next() {
		var s string
		row.Scan(&s)
		a.Equal("\001", s)
	}
}
