package storage

import (
	"time"
)

type StringObject struct {
	last  time.Time
	value string
}

func NewStringObject(value string) *StringObject {
	return &StringObject{
		last: time.Now(),
		value: value,
	}
}


func (c *StringObject) LastAccess() time.Time {
	return c.last
}

func (c *StringObject) Touch() {
	c.last = time.Now()
}

func (c *StringObject) Type() string {
	return "string"
}

func (c *StringObject) String() string {
	return c.value
}

func (c *StringObject) Len() int {
	return len(c.value)
}


