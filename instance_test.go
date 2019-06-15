package gomaxcompute

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeInstanceResult(t *testing.T) {
	a := assert.New(t)

	s := `<?xml version="1.0"?>
<Instance>
 <Tasks>
   <Task Type="SQL">
     <Name>AnonymousSQLTask</Name>
     <Result %s %s><![CDATA[%s]]></Result>
   </Task>
 </Tasks>
</Instance>`

	{
		data := "abc123!?$*&()'-=@~"
		sEnc := base64.StdEncoding.EncodeToString([]byte(data))
		content, err := decodeInstanceResult([]byte(fmt.Sprintf(s, "Transform=\"Base64\"", "Format=\"csv\"", sEnc)))
		a.NoError(err)
		a.Equal(data, content)
	}

	{
		data := "1,2,3"
		content, err := decodeInstanceResult([]byte(fmt.Sprintf(s, "", "Format=\"csv\"", data)))
		a.NoError(err)
		a.Equal(data, content)
	}

	{
		_, err := decodeInstanceResult([]byte(fmt.Sprintf(s, "", "Format=\"text\"", "1,2,3")))
		a.Error(err) // unsupported format text
	}

	{
		_, err := decodeInstanceResult([]byte(fmt.Sprintf(s, "Transform=\"zip\"", "Format=\"csv\"", "1,2,3")))
		a.Error(err) // unsupported transform zip
	}

}
