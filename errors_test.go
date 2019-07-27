package gomaxcompute

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseErrorCode(t *testing.T) {
	a := assert.New(t)

	code, err := parseErrorCode(`ParseError: {ODPS-0130161:[1,6] Parse exception - invalid token ';'}`)
	a.Nil(err)
	a.Equal(`0130161`, code)

	code, err = parseErrorCode(`ODPS-0130131:[1,15] Table not found - table gomaxcompute_driver_w7u.i_dont_exist cannot be resolved`)
	a.Nil(err)
	a.Equal(`0130131`, code)
}
