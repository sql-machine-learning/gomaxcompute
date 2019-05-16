package gomaxcompute

import (
	"encoding/xml"
	"errors"
	"strings"
)

type instanceStatus struct {
	XMLName xml.Name `xml:"Instance"`
	Status  string   `xml:"Status"`
}

type instanceResult struct {
	XMLName xml.Name `xml:"Instance"`
	Result  string   `xml:"Tasks>Task>Result"`
}

// instance types：SQL
func (conn *odpsConn) createInstance(job *odpsJob) (string, error) {
	if job == nil {
		return "", errors.New("nil job")
	}

	// Create
	res, err := conn.request(methodPost, conn.resource("/instances"), job.XML())
	if err != nil {
		return "", err
	}
	if _, err = parseResponseBody(res); err != nil && err != errNilBody {
		return "", err
	}

	// Parse response header "Location" to get instance ID
	ins := location2InstanceID(res.Header.Get("Location"))
	if ins == "" {
		return "", errors.New("no instance id")
	}
	return ins, nil
}

// parse instance id
func location2InstanceID(location string) string {
	pieces := strings.Split(location, "/")
	if len(pieces) < 2 {
		return ""
	}
	return pieces[len(pieces)-1]
}

// instance status：Running/Suspended/Terminated
func (conn *odpsConn) getInstanceStatus(instanceID string) (string, error) {
	res, err := conn.request(methodGet, conn.resource("/instances/"+instanceID), nil)
	if err != nil {
		return "", err
	}
	body, err := parseResponseBody(res)
	if err != nil {
		return "", err
	}

	var is instanceStatus
	err = xml.Unmarshal(body, &is)
	if err != nil {
		return "", err
	}
	return is.Status, nil
}

// getInstanceResult is valid while instance status is `Terminated`
// notice: records up to 10000 by limitation, and result type is string
func (conn *odpsConn) getInstanceResult(instanceID string) (string, error) {
	rsc := conn.resource("/instances/"+instanceID, pair{k: "result"})
	rsp, err := conn.request(methodGet, rsc, nil)
	if err != nil {
		return "", err
	}
	body, err := parseResponseBody(rsp)
	if err != nil {
		return "", err
	}
	var ir instanceResult
	if err = xml.Unmarshal(body, &ir); err != nil {
		return "", err
	}
	return ir.Result, nil
}
