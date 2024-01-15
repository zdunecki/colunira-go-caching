package database

import "fmt"

type DatabaseInterface interface {
	GetById(id string) (interface{}, error)
	SetById(id string, data interface{}) error
}

type Database struct {
	Data map[string]interface{}
}

func (d *Database) GetById(id string) (interface{}, error) {
	fmt.Printf("Database access for %s\n", id)
	return d.Data[id], nil
}

func (d *Database) SetById(id string, data interface{}) error {
	d.Data[id] = data
	return nil
}
