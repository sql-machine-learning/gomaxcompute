package gomaxcompute

import (
	"encoding/xml"
	"strconv"
)

// musts: Priority, Tasks
type odpsJob struct {
	xml.Marshaler
	Name     string              `xml:"Name"`
	Comment  string              `xml:"Comment"`
	Priority int                 `xml:"Priority"`
	Tasks    map[string]odpsTask `xml:"Tasks"`
}

func newJob(tasks ...odpsTask) *odpsJob {
	taskMap := make(map[string]odpsTask, len(tasks))
	for _, t := range tasks {
		taskMap[t.GetName()] = t
	}
	return &odpsJob{
		Priority: 1,
		Tasks:    taskMap,
	}
}

func newSQLJob(sql string) *odpsJob {
	return newJob(newAnonymousSQLTask(sql, nil))
}

func (j *odpsJob) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "Instance"}})
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "Job"}})

	// optional job name
	if j.Name != "" {
		e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "Comment"}})
		e.EncodeToken(xml.CharData(j.Name))
		e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Comment"}})
	}

	// optional job comment
	if j.Comment != "" {
		e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "Comment"}})
		e.EncodeToken(xml.CharData(j.Comment))
		e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Comment"}})
	}

	// Job.Priority
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "Priority"}})
	e.EncodeToken(xml.CharData(strconv.Itoa(j.Priority)))
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Priority"}})

	// Job.Tasks
	e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "Tasks"}})
	for _, task := range j.Tasks {
		if err = e.EncodeElement(task, xml.StartElement{Name: xml.Name{Local: "Task"}}); err != nil {
			return
		}
	}
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Tasks"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Job"}})
	e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "Instance"}})
	return e.Flush()
}

func (j *odpsJob) SetName(name string) {
	j.Name = name
}

func (j *odpsJob) SetComment(comment string) {
	j.Comment = comment
}

func (j *odpsJob) SetPriority(priority int) {
	j.Priority = priority
}

func (j *odpsJob) SetTasks(ts ...odpsTask) {
	taskMap := make(map[string]odpsTask, len(ts))
	for _, t := range ts {
		taskMap[t.GetName()] = t
	}
	j.Tasks = taskMap
}

func (j *odpsJob) AddTask(t odpsTask) {
	j.Tasks[t.GetName()] = t
}

func (j *odpsJob) XML() []byte {
	body, _ := xml.MarshalIndent(j, "    ", "    ")
	return body
}
