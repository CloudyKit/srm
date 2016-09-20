package scheme

import (
	"github.com/CloudyKit/framework/validation"
	"fmt"
)

// FieldDef holds methods
type FieldDef struct {
	f *Field
	d *Def
}

func (f FieldDef) PrimaryKey() FieldDef {
	f.d.checkDone()
	if f.d.PrimaryKey != nil {
		panic(fmt.Errorf("Scheme(%s): a primary key is already defined for this scheme", f.d.Type.Name()))
	}
	f.f.PrimaryKey = true
	f.d.PrimaryKey = f.f
	return f
}

// Validation adds validators to the field
func (f FieldDef) Validation(testers ...validation.Tester) FieldDef {
	f.d.checkDone()
	f.f.Testers = append(f.f.Testers, testers...)
	return f
}

// Meta adds special metadata to the field
func (f FieldDef) Meta(m ...interface{}) FieldDef {
	f.d.checkDone()
	f.f.metadata = append(f.f.metadata, m...)
	return f
}
