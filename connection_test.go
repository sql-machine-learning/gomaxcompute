package gomaxcompute

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	a := assert.New(t)
	db, err := sql.Open("maxcompute", cfg4test.FormatDSN())
	a.NoError(err)

	for i := 0; i < 100; i++ {
		const stmt = `SELECT * FROM iris_test LIMIT 1;`
		_, err = db.Query(stmt)
		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(-1)
		}
	}
}
