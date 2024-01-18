package mocks

import (
	"time"
)

type Database struct {
	Data          map[string]interface{}
	OperationTime time.Duration
}

func (d *Database) GetById(id string) interface{} {
	time.Sleep(d.OperationTime)
	return d.Data[id]
}

func (d *Database) SetById(id string, data interface{}) {
	d.Data[id] = data
}
