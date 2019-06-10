package gomaxcompute

import (
	"encoding/xml"

	"github.com/twinj/uuid"
)

const anonymousSQLTask = "AnonymousSQLTask"

type odpsTask interface {
	GetName() string
	GetComment() string
	xml.Marshaler
}

type cdata struct {
	String string `xml:",cdata"`
}

type odpsSQLTask struct {
	// default: AnonymousSQLTask
	Name    string
	Query   string
	Comment string
	//	uuid settings
	Config map[string]string
}

func (s odpsSQLTask) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "SQL"}})
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "Name"}})
	e.EncodeToken(xml.CharData(anonymousSQLTask))
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Name"}})

	if s.Comment != "" {
		e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "Comment"}})
		e.EncodeToken(xml.CharData(s.Comment))
		e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Comment"}})
	}

	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "Config"}})
	for key, value := range s.Config {
		e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "Property"}})
		e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "Name"}})
		e.EncodeToken(xml.CharData(key))
		e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Name"}})
		e.EncodeElement(cdata{value}, xml.StartElement{Name: xml.Name{Local: "Value"}})
		e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Property"}})
	}
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Config"}})
	e.EncodeElement(cdata{s.Query}, xml.StartElement{Name: xml.Name{Local: "Query"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "SQL"}})
	return e.Flush()
}

func (s *odpsSQLTask) GetName() string {
	return s.Name
}

func (s *odpsSQLTask) GetComment() string {
	return s.Comment
}

func newAnonymousSQLTask(query string, config map[string]string) odpsTask {
	return newSQLTask(anonymousSQLTask, query, config)
}

func newSQLTask(name, query string, config map[string]string) odpsTask {
	if config == nil {
		config = map[string]string{
			"uuid":     uuid.NewV4().String(),
			"settings": `{"odps.sql.udf.strict.mode": "true"}`,
		}
	}
	// maxcompute sql ends with ';', different from mysql/hive/...
	return &odpsSQLTask{
		Name:   name,
		Query:  query,
		Config: config,
	}
}
