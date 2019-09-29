package gomaxcompute

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cfg4test = &Config{
	AccessID:   os.Getenv("ODPS_ACCESS_ID"),
	AccessKey:  os.Getenv("ODPS_ACCESS_KEY"),
	Project:    os.Getenv("ODPS_PROJECT"),
	Endpoint:   os.Getenv("ODPS_ENDPOINT"),
	QueryHints: map[string]string{"odps.sql.mapper.split_size": "16"},
}

func TestSQLOpen(t *testing.T) {
	a := assert.New(t)
	db, err := sql.Open("maxcompute", cfg4test.FormatDSN())
	defer db.Close()
	a.NoError(err)
}

func TestQuerySettings(t *testing.T) {
	a := assert.New(t)
	db, err := sql.Open("maxcompute", cfg4test.FormatDSN())
	a.NoError(err)
	_, err = db.Query("SELECT * FROM gomaxcompute_test LIMIT;")
	a.NoError(err)
}
