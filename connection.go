package gomaxcompute

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	waitInteveralMs    = 1000
	tunnelHTTPProtocal = "http"
	terminated         = "Terminated"
	methodGet          = "GET"
	methodPost         = "POST"
)

type odpsConn struct {
	*http.Client
	*Config
}

// ODPS does not support transaction
func (*odpsConn) Begin() (driver.Tx, error) {
	return nil, nil
}

// ODPS does not support Prepare
func (*odpsConn) Prepare(query string) (driver.Stmt, error) {
	panic("Not implemented")
}

// Goodps accesses server by restful, so Close() do nth.
func (*odpsConn) Close() error {
	return nil
}

// Implements database/sql/driver.Execer. Notice result is nil
func (conn *odpsConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	ins, err := conn.wait(query, args)
	if err != nil {
		return nil, err
	}
	_, err = conn.getInstanceResult(ins)
	if err != nil {
		return nil, err
	}
	// FIXME(weiguo): precise result
	return &odpsResult{-1, -1}, nil
}

// Implements database/sql/driver.Queryer
func (conn *odpsConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	ins, err := conn.wait(query, args)
	if err != nil {
		return nil, err
	}

	// get tunnel server
	tunnelServer, err := conn.getTunnelServer()
	if err != nil {
		return nil, err
	}
	log.Infof("++Query:[%s] tunnel:[%s]", query, tunnelServer)
	// get meta by tunnel
	meta, err := conn.getResultMeta(ins, tunnelServer)
	if err != nil {
		log.Error("------------above error------------------")
		return nil, err
	}

	res, err := conn.getInstanceResult(ins)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(res, "ODPS-") {
		return nil, errors.New(res)
	}
	return newRows(meta, res)
}

func (conn *odpsConn) getResultMeta(instance, tunnelServer string) (*resultMeta, error) {
	endpoint := fmt.Sprintf("%s://%s", tunnelHTTPProtocal, tunnelServer)
	rsc := fmt.Sprintf("/projects/%s/instances/%s", conn.Project, instance)
	params := url.Values{}
	params.Add(currentProject, conn.Project)
	params.Add("downloads", "")
	url := rsc + "?" + params.Encode()

	rsp, err := conn.requestEndpoint(endpoint, methodPost, url, nil)
	if err != nil {
		return nil, err
	}
	body, err := parseResponseBody(rsp)
	if err != nil {
		return nil, err
	}

	meta := resultMeta{}
	err = json.Unmarshal(body, &meta)
	return &meta, err
}

func (conn *odpsConn) getTunnelServer() (string, error) {
	rsp, err := conn.request(methodGet, conn.resource("/tunnel"), nil)
	if err != nil {
		return "", err
	}

	url, err := parseResponseBody(rsp)
	if err != nil {
		return "", err
	}
	return string(url), nil
}

func (conn *odpsConn) wait(query string, args []driver.Value) (string, error) {
	if len(args) > 0 {
		query = fmt.Sprintf(query, args)
	}

	ins, err := conn.createInstance(newSQLJob(query))
	if err != nil {
		return "", err
	}
	if err := conn.poll(ins, waitInteveralMs); err != nil {
		return "", err
	}
	return ins, nil
}

func (conn *odpsConn) poll(instanceID string, interval int) error {
	du := time.Duration(interval) * time.Millisecond
	for {
		status, err := conn.getInstanceStatus(instanceID)
		if err != nil {
			return err
		}
		if status == terminated {
			break
		}
		time.Sleep(du)
	}
	return nil
}
