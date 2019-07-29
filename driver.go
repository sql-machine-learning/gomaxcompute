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
type Driver struct {
	cfg *Config
}

func (d Driver) Open(dsn string) (driver.Conn, error) {
	var err error
	if d.cfg, err = ParseDSN(dsn); err != nil {
		return nil, err
	}
	return &odpsConn{&http.Client{}, d.cfg}, nil
}
