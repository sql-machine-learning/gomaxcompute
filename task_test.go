package gomaxcompute

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLTask_MarshalXML(t *testing.T) {
	a := assert.New(t)
	task := newAnonymousSQLTask("SELECT 1;", nil)
	_, err := xml.Marshal(task)
	a.NoError(err)
}
