package gomaxcompute

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
)

// register driver
func init() {
	sql.Register("maxcompute", &Driver{})
}

// impls database/sql/driver.Driver
type Driver struct{}

func (d Driver) Open(dsn string) (driver.Conn, error) {
	cfg, err := ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	return &odpsConn{&http.Client{}, cfg}, nil
}
