package gomaxcompute

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twinj/uuid"
)

func TestJob_GenCreateInstanceXml(t *testing.T) {
	a := assert.New(t)
	task := newAnonymousSQLTask("SELECT 1;", map[string]string{
		"uuid":     uuid.NewV4().String(),
		"settings": `{"odps.sql.udf.strict.mode": "true"}`,
	})
	job := newJob(task)
	if job == nil {
		t.Error("fail to create new job")
	}

	_, err := xml.MarshalIndent(job, "", "    ")
	a.NoError(err)
}

func TestNewSQLJob(t *testing.T) {
	job := newSQLJob(`SELECT 1;`)
	fmt.Println(string(string(job.XML())))
}
