package driver

import (
	"errors"
	"github.com/CloudyKit/srm/change"
	"github.com/CloudyKit/srm/query"
	"github.com/CloudyKit/srm/scheme"
)

type NoTransactions struct{}

func (NoTransactions) Begin() error {
	return errors.New("Transactions are not supported by this driver")
}
func (NoTransactions) Commit() error {
	return errors.New("Transactions are not supported by this driver")
}
func (NoTransactions) RowBack() error {
	return errors.New("Transactions are not supported by this driver")
}

type Result interface {
	NumOfRecords() int
	ScanRow(v ...interface{}) error
}

type Driver interface {
	UseScheme(name string, s *scheme.Scheme) error

	Begin() error
	Commit() error
	RowBack() error

	Retrieve(name string, s *scheme.Scheme, key string) Result

	Create(name string, s *scheme.Scheme, operations ...change.Set) (string, error)
	Modify(name string, s *scheme.Scheme, key string, operations ...change.Set) (int, error)
	Remove(name string, s *scheme.Scheme, key string) (int, error)

	Search(name string, s *scheme.Scheme, q *query.Query) Result
	SearchAndModify(name string, s *scheme.Scheme, q *query.Query, operations ...change.Operation) (numofmodified int, err error)
	SearchAndRemove(name string, s *scheme.Scheme, q *query.Query) (numofmodified int, err error)
}
