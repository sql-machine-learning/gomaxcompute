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

	const sql = `select * from yiyang_test_table1;`
	_, err = db.Query(sql)
	a.NoError(err)
}
