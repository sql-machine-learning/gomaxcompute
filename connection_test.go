package gomaxcompute

import (
	"database/sql"
	"fmt"
	"math/rand"
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

func TestExec(t *testing.T) {
	a := assert.New(t)
	db, err := sql.Open("maxcompute", cfg4test.FormatDSN())
	a.NoError(err)

	tn := fmt.Sprintf("unitest%d", rand.Int())
	_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(shop_name STRING);", tn))
	a.NoError(err)
	_, err = db.Exec(fmt.Sprintf("DROP TABLE %s;", tn))
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
