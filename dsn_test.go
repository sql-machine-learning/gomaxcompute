package gomaxcompute

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_ParseDSN(t *testing.T) {
	a := assert.New(t)

	correct := "access_id:access_key@service.com/api?curr_project=test_ci&scheme=http"
	config, err := ParseDSN(correct)
	a.NoError(err)
	a.Equal("access_id", config.AccessID)
	a.Equal("access_key", config.AccessKey)
	a.Equal("test_ci", config.Project)
	a.Equal("http://service.com/api", config.Endpoint)

	badDSN := []string{
		"", // empty
		":access_key@service.com/api?curr_project=test_ci&scheme=http",              // missing access_id
		"access_idaccess_key@service.com/api?curr_project=test_ci&scheme=http",      // missing :
		"access_id:@service.com/api?curr_project=test_ci&scheme=http",               // missing @
		"access_id:access_key@?curr_project=test_ci&scheme=http",                    // missing endpoint
		"access_id:access_key@service.com/apicurr_project=test_ci&scheme=http",      // missing ?
		"access_id:access_key@service.com/api?scheme=http",                          // missing curr_project
		"access_id:access_key@service.com/api?curr_project=test_ci",                 // missing scheme
		"access_id:access_key@service.com/api?curr_project=test_ci&scheme=whatever", // invalid scheme
		"access_id:access_key@service.com/api?curr_project=test_ci&scheeeeeee=http", // invalid name
	}
	for _, dsn := range badDSN {
		_, err = ParseDSN(dsn)
		a.Error(err)
	}
}

func TestConfig_FormatDSN(t *testing.T) {
	a := assert.New(t)
	config := Config{
		AccessID: "access_id",
		AccessKey: "access_key",
		Project: "test_ci",
		Endpoint: "http://service.com/api"}
	a.Equal("access_id:access_key@service.com/api?curr_project=test_ci&scheme=http", config.FormatDSN())
}

func TestConfig_ParseAndFormatRoundTrip(t *testing.T) {
	a := assert.New(t)
	dsn := "access_id:access_key@service.com/api?curr_project=test_ci&scheme=http"

	config, err := ParseDSN(dsn)
	a.NoError(err)
	a.Equal(dsn, config.FormatDSN())
}