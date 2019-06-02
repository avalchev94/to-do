package database

import (
	"net/url"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/clock"
)

type Parameter struct {
	Table    string
	Name     string
	Operator string
	Value    interface{}
}

var allowedParameters = map[string]Parameter{
	"date":            Parameter{"task_schedule", "date", "=", date.Date{}},
	"dateGreater":     Parameter{"task_schedule", "date", ">", date.Date{}},
	"time":            Parameter{"task_schedule", "time", "=", clock.Clock(0)},
	"timeGreater":     Parameter{"task_schedule", "time", ">", clock.Clock(0)},
	"finished":        Parameter{"task_schedule", "finished", "=", &time.Time{}},
	"finishedGreater": Parameter{"task_schedule", "finished", ">", &time.Time{}},
}

func (p Parameter) encodeSQL() string {
	var value string
	switch p.Value.(type) {
	case date.Date:
		value = "'" + p.Value.(date.Date).String() + "'"
	case clock.Clock:
		value = "'" + p.Value.(clock.Clock).HhMmSs() + "'"
	}
	return p.Table + "." + p.Name + p.Operator + value
}

type Parameters []Parameter

func ParseParameters(url *url.URL) Parameters {
	params := Parameters{}

	for k, v := range url.Query() {
		if p, ok := allowedParameters[k]; ok {
			switch p.Value.(type) {
			case date.Date:
				date, err := date.AutoParse(v[0])
				if err != nil {
					continue
				}
				p.Value = date
			case clock.Clock:
				clock, err := clock.Parse(v[0])
				if err != nil {
					continue
				}
				p.Value = clock

			}
			params = append(params, p)
		}
	}
	return params
}

func (p Parameters) encodeSQL(mappings map[string]string) string {
	where := ""
	for _, param := range p {
		if len(where) > 0 {
			where += " AND "
		}

		if value, ok := mappings[param.Table]; ok {
			param.Table = value
			where += param.encodeSQL()
		}
	}
	return where
}
