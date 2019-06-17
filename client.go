package gomaxcompute

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
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
func (conn *odpsConn) request(method, resource string, body []byte, header ...pair) (res *http.Response, err error) {
	return conn.requestEndpoint(conn.Endpoint, method, resource, body, header...)
}

func (conn *odpsConn) requestEndpoint(endpoint, method, resource string, body []byte, header ...pair) (res *http.Response, err error) {
	var req *http.Request
	url := endpoint + resource
	if body != nil {
		if req, err = http.NewRequest(method, url, bytes.NewBuffer(body)); err != nil {
			return nil, errors.WithStack(err)
		}
		req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	} else {
		if req, err = http.NewRequest(method, url, nil); err != nil {
			return nil, errors.WithStack(err)
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
	log.Debug("--------------------------------")
	log.Debugf("request.url: %v", req.URL.String())
	log.Debugf("request.header: %v", req.Header)
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
		// query parameters need to be hashed in alphabet order
		keys := make([]string, len(urlParams))
		i := 0
		for k := range urlParams {
			keys[i] = k
			i++
		}
		sort.Strings(keys)

		first := true
		for _, k := range keys {
			if first {
				canonicalResource.WriteByte('?')
				first = false
			} else {
				canonicalResource.WriteByte('&')
			}
			canonicalResource.WriteString(k)
			v := urlParams[k]
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

func parseResponse(rsp *http.Response) ([]byte, error) {
	if rsp == nil || rsp.Body == nil {
		return nil, errNilBody
	}
	log.Debugf("response code: %v", rsp.StatusCode)

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode >= 400 {
		return nil, parseResponseError(rsp.StatusCode, body)
	}
	return body, err
}

func parseResponseError(statusCode int, body []byte) error {
	re := responseError{}
	if err := json.Unmarshal(body, &re); err != nil {
		ie := instanceError{}
		if err := xml.Unmarshal([]byte(body), &ie); err != nil {
			return errors.WithStack(fmt.Errorf("response error: %d, %s", statusCode, string(body)))
		}
		return errors.WithStack(fmt.Errorf("%s: %s", ie.Code, ie.Message))
	}
	return errors.WithStack(fmt.Errorf("response error: %d, %s. %s",
		statusCode, re.Code, re.Message))
}
