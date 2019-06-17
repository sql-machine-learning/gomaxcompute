package gomaxcompute

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type instanceStatus struct {
	XMLName xml.Name `xml:"Instance"`
	Status  string   `xml:"Status"`
}

type Result struct {
	Content   string `xml:",cdata"`
	Transform string `xml:"Transform,attr"`
	Format    string `xml:"Format,attr"`
}

type instanceResult struct {
	XMLName xml.Name `xml:"Instance"`
	Result  Result   `xml:"Tasks>Task>Result"`
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
	if _, err = parseResponse(res); err != nil && err != errNilBody {
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
	body, err := parseResponse(res)
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
	body, err := parseResponse(rsp)
	if err != nil {
		return "", err
	}
	return decodeInstanceResult(body)
}

func decodeInstanceResult(result []byte) (string, error) {
	var ir instanceResult
	if err := xml.Unmarshal(result, &ir); err != nil {
		return "", err
	}

	if ir.Result.Format == "text" {
		log.Debug(ir.Result.Content)
		// ODPS errors are text begin with "ODPS-"
		if strings.HasPrefix(ir.Result.Content, "ODPS-") {
			return "", errors.WithStack(errors.New(ir.Result.Content))
		}
		// FIXME(tony): the result non-query statement usually in text format.
		// Go's database/sql API only supports lastId and affectedRows.
		return "", nil
	}

	if ir.Result.Format != "csv" {
		return "", errors.WithStack(fmt.Errorf("unsupported format %v", ir.Result.Format))
	}

	switch ir.Result.Transform {
	case "":
		return ir.Result.Content, nil
	case "Base64":
		content, err := base64.StdEncoding.DecodeString(ir.Result.Content)
		if err != nil {
			return "", err
		}
		return string(content), err
	default:
		return "", errors.WithStack(fmt.Errorf("unsupported transform %v", ir.Result.Transform))
	}
}
