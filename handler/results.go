package handler

import "fmt"

type result struct {
	Val interface{}
	Err error
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
	return r.Val
}

func (r result) String() string {
	return fmt.Sprintf("%v", r.Value())
}

func newResult(val interface{}, err error) result {
	return result{val, err}
}
