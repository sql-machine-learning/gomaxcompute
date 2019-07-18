package gomaxcompute

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeInstanceError(t *testing.T) {
	a := assert.New(t)

	s := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
	<Code>ParseError</Code>
	<Message><![CDATA[ODPS-0130161:[1,1] Parse exception - invalid token 'SLEECT']]></Message>
	<RequestId>5583C0893F7EC2D353</RequestId>
	<HostId>odps.aliyun.com</HostId>
</Error>`

	var ie instanceError
	a.NoError(xml.Unmarshal([]byte(s), &ie))
	a.Equal(`ParseError`, ie.Code)
	a.Equal(`ODPS-0130161:[1,1] Parse exception - invalid token 'SLEECT'`, ie.Message.CDATA)
	a.Equal(`5583C0893F7EC2D353`, ie.RequestId)
	a.Equal(`odps.aliyun.com`, ie.HostId)
}

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
		a.NoError(err)
	}

	{
		_, err := decodeInstanceResult([]byte(fmt.Sprintf(s, "Transform=\"zip\"", "Format=\"csv\"", "1,2,3")))
		a.Error(err) // unsupported transform zip
	}

}
