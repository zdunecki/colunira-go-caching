package integration_test

import (
	"time"
)

type MockDatabase struct {
	Data map[string]interface{}
	OperationTime time.Duration
}

func (d *MockDatabase) GetById(id string) interface{} {
	time.Sleep(d.OperationTime)
	return d.Data[id]
}

func (d *MockDatabase) SetById(id string, data interface{}) {
	d.Data[id] = data
}
