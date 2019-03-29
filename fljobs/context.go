package fljobs

import "github.com/nerde/fuji-lane-back/fldiagnostics"

// Context of a job execution
type Context struct {
	*fldiagnostics.Diagnostics
	Queue string
	Args  []interface{}
}

// NewContext creates a new job Context
func NewContext(queue string, args []interface{}) *Context {
	return &Context{
		Diagnostics: &fldiagnostics.Diagnostics{},
		Queue:       queue,
		Args:        args,
	}
}
