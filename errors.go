package gomaxcompute

import (
	"fmt"
	"regexp"
)

var (
	reErrorCode = regexp.MustCompile(`^.*ODPS-([0-9]+):.*$`)
)

func parseErrorCode(message string) (string, error) {
	sub := reErrorCode.FindStringSubmatch(message)
	if len(sub) != 2 {
		return "", fmt.Errorf("fail parse error: %s", message)
	}

	return sub[1], nil
}

// MaxcomputeError is an error type which represents a single Maxcompute error
// Please refer to https://www.alibabacloud.com/help/doc-detail/64654.htm
// for the list of SQL common error code
type MaxcomputeError struct {
	Code    string
	Message string
}

func (e *MaxcomputeError) Error() string {
	return e.Message
}
