package gomaxcompute

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
)

// register driver
func init() {
	sql.Register("maxcompute", &ODPSDriver{})
}

// ODPSDriver impls database/sql/driver.Driver
type ODPSDriver struct {
	conn odpsConn
}

func (o ODPSDriver) Open(dsn string) (driver.Conn, error) {
	cfg, err := ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	o.conn = odpsConn{&http.Client{}, cfg, nil}
	return &o.conn, nil
}

// SetQuerySettings sets the global query settings.
// TODO(Yancey1989): add the one-off query settings interface.
func (o ODPSDriver) SetQuerySettings(hints map[string]string) error {
	o.conn.hints = hints
	return nil
}
