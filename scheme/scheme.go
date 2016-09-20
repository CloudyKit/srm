package scheme

import (
	"github.com/CloudyKit/framework/validation"
	"sync/atomic"
	"reflect"
	"errors"
)

type RelKind int

const (
	HasNoRel = 0

	Belongs = 1 << iota
	HasOne
	HasMany

	BelongsManyThrough
	BelongsOneThrough
	HasManyThrough
	HasOneThrough

	BelongsEmbed
	HasOneEmbed
	HasManyEmbed
)


// Index abstract index definition
type Index struct {
	Fields []string
	Type   string
}


// Def is responsible the define the metadata for the scheme
type Def Scheme

// Scheme represents a database entity, fields and the references between different entities
type Scheme struct {
	PrimaryKey *Field
	Namespace  string

	Type       reflect.Type

	metadata   []interface{} // holds a list metadata values

	indexes    []Index
	fields     []*Field

	done       uintptr
}

type Definer interface {
	Def(*Def)
}

// Field holds metadata about the fields and the references
type Field struct {
	PrimaryKey bool                // bool if the field is the primaryKey

	Name       string              // the name of the field in the go Type
	RealName   string              // the name of the field in the database

	Type       reflect.Type        // the type of the field
	Index      []int               // the index of the field

	RelKind    RelKind             // Kind of relation, or HasNoRel for normal field
	RelDst     string              // the target fieldName in the target go Type

	RelSrc     string              // the source fieldName in the source go Type
	RelNs      string              // the name of the namespace used to store the relations through

	Testers    []validation.Tester // list of validators, to constrain the insert and update of values

	metadata   []interface{}       // metadata used by the driver, ex: declare the table
}

func (d *Def) Done() {

	// will panic if Done was called before
	d.checkDone()

	// mark the scheme as done, all methods that can modify the scheme will panic
	// after a call to .Done
	atomic.StoreUintptr(&d.done, 1)

	// if the def func didn't specify the namespace name, we generate one from the go type name
	if d.Namespace == "" {
		d.Namespace = genName(d.Type.String())
	}
}

// checkDone checks and panic if the scheme is marked as done
func (d *Def) checkDone() {
	if atomic.LoadUintptr(&d.done) > 0 {
		panic(errors.New("Trying to modify a scheme definition outside scheme call"))
	}
}

// Meta returns metadata stored during the field definition
// this metadata is intended to be used by the driver
func (scheme *Field) Meta() []interface{} {
	return scheme.metadata
}

// Meta returns metadata stored during the field definition
// this metadata is intended to be used by the driver
func (scheme *Scheme) Meta() []interface{} {
	return scheme.metadata
}

// FieldByName returns a field by name
func (scheme *Scheme) FieldByName(name string) (field *Field, found bool) {
	for _, f := range scheme.fields {
		if found = f.Name == name; found {
			field = f
			return
		}
	}
	return
}

// Indexes returns the list of indexes in the scheme
func (scheme *Scheme) Indexes() []Index {
	return scheme.indexes
}

// Fields returns the list of fields in the scheme
func (scheme *Scheme) Fields() []*Field {
	return scheme.fields
}
