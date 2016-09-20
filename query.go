package srm

import (
	"github.com/CloudyKit/srm/driver"
	"github.com/CloudyKit/srm/query"
	"github.com/CloudyKit/srm/scheme"
)

type Result struct {
	d *Session
	q *query.Query
	s *scheme.Scheme
	r driver.Result
}

func (q *Result) Select(fields ...string) *Result {
	return nil
}

func (q *Result) For(i interface{}) error {
	return nil
}

func (r *Result) NumOfRecords() int {
	return r.r.NumOfRecords()
}

func (r *Result) Fetch(target interface{}) error {
	return nil
}

func (r *Result) FetchNext(target interface{}) bool {
	return true
}

func (r *Result) FetchAll(target interface{}) error {
	return nil
}
