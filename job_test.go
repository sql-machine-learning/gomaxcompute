package gomaxcompute

import (
	"encoding/xml"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJob_GenCreateInstanceXml(t *testing.T) {
	a := assert.New(t)
	task := newAnonymousSQLTask("SELECT 1;", map[string]string{
		"uuid":     uuid.NewString(),
		"settings": `{"odps.sql.udf.strict.mode": "true"}`,
	})
	job := newJob(task)
	if job == nil {
		t.Error("fail to create new job")
	}

	_, err := xml.MarshalIndent(job, "", "    ")
	a.NoError(err)
}
