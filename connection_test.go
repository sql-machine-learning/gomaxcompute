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

	_, err = db.Query("SELECT * FROM iris_test LIMIT 1;")
	a.NoError(err)
}
