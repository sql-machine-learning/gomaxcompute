package gomaxcompute

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
)

// register driver
func init() {
	sql.Register("maxcompute", &MaxcomputeDriver{})
}

// MaxcomputeDriver impls database/sql/driver.Driver
type MaxcomputeDriver struct {
	conn odpsConn
}

func (o MaxcomputeDriver) Open(dsn string) (driver.Conn, error) {
	cfg, err := ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	o.conn = odpsConn{&http.Client{}, cfg, nil}
	return &o.conn, nil
}

// SetQuerySettings sets the global query settings.
// TODO(Yancey1989): add the one-off query settings interface.
func (o MaxcomputeDriver) SetQuerySettings(hints map[string]string) error {
	o.conn.hints = hints
	return nil
}
