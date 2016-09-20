// SRM struct relation manager
package srm

import (
	"github.com/CloudyKit/srm/change"
	"github.com/CloudyKit/srm/driver"
	"github.com/CloudyKit/srm/query"
	"github.com/CloudyKit/srm/scheme"
	"github.com/CloudyKit/framework/validation"

	"bytes"
	"reflect"
	"unicode"
)

// Model holds database data, like relations, primary key
type Model struct {
	PrimaryKey string
}

// IModel represents a struct containing a Model data
type IModel interface {
	getModel() *Model
}

// New returns a new database context
func New(driver driver.Driver) *Session {
	db := new(Session)
	db.driver = driver
	return db
}

type namedScheme struct {
	Scheme *scheme.Scheme
	Name   string
}

// AbstractDB holds a database context, an AbstractDB is responsible to operate on a database
type Session struct {
	driver  driver.Driver
	schemes map[reflect.Type]namedScheme
}

// Use maps an entity name to a model
func (db *Session) Use(entityName string, model IModel) *scheme.Scheme {
	t := reflect.TypeOf(model)
	_scheme := getScheme(t)
	db.schemes[t] = namedScheme{Name: entityName, Scheme: _scheme}
	return _scheme
}

func genName(text string) string {
	buf := bytes.NewBuffer(make([]byte, 0, len(text)))
	for k, r := range text {
		if unicode.IsUpper(r) {
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
			if unicode.IsLower(r) {
				if k+1 < len(text) && unicode.IsUpper(text[k+1]) {
					buf.WriteByte('_')
				}
			}
		}
	}
	return buf.String()
}

// namedScheme used by the database context to get a scheme with it's entity name
func (db *Session) namedScheme(t reflect.Type) (string, *scheme.Scheme) {
	_namedScheme, ok := db.schemes[t]
	if ok == false {
		return genName(t.String()), getScheme(t)
	}
	return _namedScheme.Name, _namedScheme.Scheme
}

// Begin begins a transaction
func (db *Session) Begin() error {
	return db.driver.Begin()
}

// Commit all operations
func (db *Session) Commit() error {
	return db.driver.Commit()
}

// RowBack cancel all performed operations
func (db *Session) RowBack() error {
	return db.driver.RowBack()
}

// Store stores a model into the database
func (db *Session) Store(m IModel) (validation.Result, error) {

	defer func() {
		if err := recover(); err != nil {
			// recover from a panic
			if _, ok := err.(*storingContext); !ok {
				panic(err)
			}
		}
	}()

	// holds the state of the store routine
	context := storingContext{
		scheme: getScheme(reflect.TypeOf(m)),
		model:  reflect.ValueOf(m),
	}

	context.storeFields()

	return context.validation.Done(), context.lastErr
}

// Retrieve  retrieves a stored model or any decedent field specified by the fields ...key
func (db *Session) Retrieve(m IModel, fields ...string) error {
	return nil
}

// Remove remove a model, will remove the model from de database
func (db *Session) Remove(p IModel) error {
	return nil
}

// Search search's for the records matching the specified query and returns the matching records
// in the result set
func (db *Session) Search(m IModel, q *query.Query) Result {
	return Result{q: q, d: db}
}

// SearchAndModify search's and modify the records matching the passed query expression with a list of changes
// passed as variadic list of operations
func (db *Session) SearchAndModify(s IModel, q *query.Query, changes ...change.Operation) (int, error) {
	return 0, nil
}

// SearchAndRemove search's and remove the records matching the passed query expression
func (db *Session) SearchAndRemove(s IModel, q *query.Query) (int, error) {
	return 0, nil
}
