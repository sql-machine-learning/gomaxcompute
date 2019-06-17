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

	queries := []string{
		`SELECT * FROM gomaxcompute_test LIMIT 1`,
		`SELECT * FROM gomaxcompute_test LIMIT 1;`,
	}
	for _, query := range queries {
		_, err = db.Query(query)
		a.NoError(err)
	}
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
