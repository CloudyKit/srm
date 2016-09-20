package scheme

import (
	"fmt"
	"reflect"
	"regexp"
)

func (def *Def) Meta(meta ...interface{}) *Def {
	def.metadata = append(def.metadata, meta...)
	return def
}

func (def *Def) getStructField(fieldName, dbName string) (reflect.Type, []int) {

	for _, f := range def.fields {

		if f.Name == fieldName {
			panic(fmt.Errorf("Scheme definition: field %s in scheme %s was already defined", fieldName, def.Type))
		}

		if f.RealName == dbName {
			panic(fmt.Errorf("Scheme definition: field %s in scheme %s was already defined", fieldName, def.Type))
		}

		//if f.RelKind == BelongsOne {
		//	if f.RelField == fieldName {
		//		panic(fmt.Errorf("Scheme definition: field %s in scheme %s was already defined", fieldName, def.Type))
		//	}
		//}

	}

	structField, found := def.Type.FieldByName(fieldName)
	if !found {
		panic(fmt.Errorf("Scheme definition: %s is inexistent in %s", fieldName, def.Type))
	}

	return structField.Type, structField.Index
}

// Field defines a schema field with the fieldName
func (def *Def) Field(fieldName string, dbName ...string) FieldDef {
	def.checkDone()

	field := &Field{Name: fieldName}

	if len(dbName) > 0 {
		field.RealName = dbName[0]
	} else {
		field.RealName = genName(fieldName)
	}

	field.Type, field.Index = def.getStructField(fieldName, field.RealName)
	def.fields = append(def.fields, field)

	return FieldDef{field, def}
}

func (def FieldDef) Field(name string, dbName ...string) FieldDef {
	return def.d.Field(name, dbName)
}


// validateOrdering string
var validateOrdering = regexp.MustCompile("^[+-]?[a-zA-Z_][a-zA-Z_0-9]*$")

// Index define a index type _typ on fields _fields
func (def *Def) Index(_typ string, _fields ...string) *Def {
	def.checkDone()

	if len(_fields) == 0 {
		panic(fmt.Errorf("Scheme Def: defining index on scheme %s without fields", _typ))
	}

	for _, field := range _fields {
		if !validateOrdering.MatchString(field) {
			panic(fmt.Errorf("Scheme Def: defining index on scheme %s with an invalid name %s", _typ, field))
		}
	}

	def.indexes = append(def.indexes, Index{Type: _typ, Fields: _fields})
	return def
}
