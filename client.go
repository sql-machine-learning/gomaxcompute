package gomaxcompute

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	errNilBody    = errors.New("nil body")
	requestGMT, _ = time.LoadLocation("GMT")
)

const currentProject = "curr_project"

// http response error
type responseError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

type pair struct {
	k string
	v string
}

// optional: Body, Header
func (conn *odpsConn) request(method, resource string,
	body []byte, header ...pair) (res *http.Response, err error) {
	return conn.requestEndpoint(conn.Endpoint, method, resource, body, header...)
}

func (conn *odpsConn) requestEndpoint(endpoint, method, resource string,
	body []byte, header ...pair) (res *http.Response, err error) {
	var req *http.Request
	url := endpoint + resource
	if body != nil {
		if req, err = http.NewRequest(method, url, bytes.NewBuffer(body)); err != nil {
			return
		}
		req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	} else {
		if req, err = http.NewRequest(method, url, nil); err != nil {
			return
		}
	}

	req.Header.Set("x-odps-user-agent", "gomaxcompute/0.0.1")
	req.Header.Set("Content-Type", "application/xml")

	if dateStr := req.Header.Get("Date"); dateStr == "" {
		gmtTime := time.Now().In(requestGMT).Format(time.RFC1123)
		req.Header.Set("Date", gmtTime)
	}
	// overwrite with user-provide header
	if header != nil || len(header) == 0 {
		for _, arg := range header {
			req.Header.Set(arg.k, arg.v)
		}
	}

	// fill curr_project
	if req.URL.Query().Get(currentProject) == "" {
		req.URL.Query().Set(currentProject, conn.Project)
	}
	conn.sign(req)
	return conn.Do(req)
}

// signature
func (conn *odpsConn) sign(r *http.Request) {
	var msg, auth bytes.Buffer
	msg.WriteString(r.Method)
	msg.WriteByte('\n')
	// common header
	msg.WriteString(r.Header.Get("Content-MD5"))
	msg.WriteByte('\n')
	msg.WriteString(r.Header.Get("Content-Type"))
	msg.WriteByte('\n')
	msg.WriteString(r.Header.Get("Date"))
	msg.WriteByte('\n')
	// canonical header
	for k, v := range r.Header {
		lowerK := strings.ToLower(k)
		if strings.HasPrefix(lowerK, "x-odps-") {
			msg.WriteString(lowerK)
			msg.WriteByte(':')
			msg.WriteString(strings.Join(v, ","))
			msg.WriteByte('\n')
		}
	}

	// canonical resource
	var canonicalResource bytes.Buffer
	epURL, _ := url.Parse(conn.Endpoint)
	if strings.HasPrefix(r.URL.Path, epURL.Path) {
		canonicalResource.WriteString(r.URL.Path[len(epURL.Path):])
	} else {
		canonicalResource.WriteString(r.URL.Path)
	}
	if urlParams := r.URL.Query(); len(urlParams) > 0 {
		first := true
		for k, v := range urlParams {
			if first {
				canonicalResource.WriteByte('?')
				first = false
			} else {
				canonicalResource.WriteByte('&')
			}
			canonicalResource.WriteString(k)
			if v != nil && len(v) > 0 && v[0] != "" {
				canonicalResource.WriteByte('=')
				canonicalResource.WriteString(v[0])
			}
		}
	}
	msg.WriteString(canonicalResource.String())

	hasher := hmac.New(sha1.New, []byte(conn.AccessKey))
	hasher.Write(msg.Bytes())
	auth.WriteString("ODPS ")
	auth.WriteString(conn.AccessID)
	auth.WriteByte(':')
	auth.WriteString(base64.StdEncoding.EncodeToString(hasher.Sum(nil)))
	r.Header.Set("Authorization", auth.String())
}

func (cred *Config) resource(resource string, args ...pair) string {
	if args == nil || len(args) == 0 {
		return fmt.Sprintf("/projects/%s%s", cred.Project, resource)
	}

	ps := url.Values{}
	for _, i := range args {
		ps.Add(i.k, i.v)
	}
	return fmt.Sprintf("/projects/%s%s?%s", cred.Project, resource, ps.Encode())
}

func parseResponseBody(rsp *http.Response) ([]byte, error) {
	if rsp == nil || rsp.Body == nil {
		return nil, errNilBody
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)

	if rsp.StatusCode >= 400 {
		re := responseError{}
		if err = json.Unmarshal(body, &re); err != nil {
			return nil, fmt.Errorf("response error: %d", rsp.StatusCode)
		}
		return nil, fmt.Errorf("response error: %d, %s. %s",
			rsp.StatusCode, re.Code, re.Message)
	}
	return body, err
}
