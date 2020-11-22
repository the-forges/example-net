package cmdparser

import (
	"fmt"

	"github.com/the-forges/example-net/internal"
)

type result struct {
	Input   string
	Cmd     internal.Command
	Handler internal.CommandHandler
	Err     error
}

func (r result) Error() string {
	if r.Err == nil {
		return ""
	}
	return r.Err.Error()
}

func (r result) Ok() bool {
	return r.Error() == ""
}

func (r result) Value() interface{} {
	blurb := "is valid"
	if !r.Ok() {
		blurb = fmt.Sprintf("is invalid, having the following error: %s", r.Error())
	}
	return fmt.Sprintf("%s %s", r.Input, blurb)
}

func (r result) String() string {
	return r.Value().(string)
}

func newResult(input string, cmd internal.Command, handler internal.CommandHandler, err error) result {
	return result{input, cmd, handler, err}
}
