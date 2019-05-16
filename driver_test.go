package gomaxcompute

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cfg4test = &Config{
	AccessID:  os.Getenv("ODPS_ACCESS_ID"),
	AccessKey: os.Getenv("ODPS_ACCESS_KEY"),
	Project:   os.Getenv("ODPS_PROJECT"),
	Endpoint:  os.Getenv("ODPS_ENDPOINT"),
}

func TestSQLOpen(t *testing.T) {
	a := assert.New(t)
	db, err := sql.Open("maxcompute", cfg4test.FormatDSN())
	defer db.Close()
	a.NoError(err)
}
