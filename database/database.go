package database

import (
	"time"
)

type DatabaseInterface interface {
	GetById(id string) interface{}
	SetById(id string, data interface{})
}

type Database struct {
	Data map[string]interface{}
}

func (d *Database) GetById(id string) interface{} {	
	time.Sleep(time.Second)
	return d.Data[id]
}

func (d *Database) SetById(id string, data interface{}) {
	d.Data[id] = data
}
