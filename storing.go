package srm

import (
	"github.com/CloudyKit/srm/change"
	"github.com/CloudyKit/srm/scheme"
	"github.com/CloudyKit/framework/validation"

	"errors"
	"reflect"
)

// storingContext a storing context holds data necessary to an store routine
type storingContext struct {
	mode       string
	validation validation.Validator
	result     validation.Result

	lastErr error

	scheme, parentScheme *scheme.Scheme
	namespace            string
	name, parentName     string
	model, parentModel   reflect.Value

	*Session
}

func (storingContext *storingContext) restoreContext(on *storingContext) {
	*on = *storingContext
}

func (storingContext *storingContext) storeFields() {

	// create a list of operations
	operations := make([]change.Set, 0, len(storingContext.scheme.Fields()))

	// walk the scheme fields
	for _, field := range storingContext.scheme.Fields() {

		switch field.RelKind {
		// if the field is a value field / has no rel,
		case scheme.HasNoRel:
			// setup validation context
			storingContext.validation.Name = field.Name
			storingContext.validation.Value = storingContext.model.FieldByIndex(field.Index)

			// run validation testers
			for i := 0; i < len(field.Testers); i++ {
				field.Testers[i](&storingContext.validation)
			}

			// check for validation errors
			if storingContext.validation.Done().Bad() {
				// stop the store routine
				panic(storingContext)
			}

			// no errors, store set operation
			operations = append(operations, change.Set{Field: field.Name, Value: storingContext.validation.Value})
		case scheme.HasOneEmbed:

			// take snapshot of the current state of the context
			snapshot := *storingContext

			// modify the context to store the child model
			storingContext.parentModel = storingContext.model
			storingContext.parentScheme = storingContext.scheme
			storingContext.parentName = storingContext.name

			storingContext.name = field.Name
			storingContext.model = storingContext.model.FieldByIndex(field.Index)
			storingContext.namespace, storingContext.scheme = storingContext.Session.namedScheme(field.Type)

			// store the child model
			storingContext.storeFields()

			// set the field value with primary key of the child model
			operations = append(operations, change.Set{
				Field: field.Name,
				Value: reflect.ValueOf(
					storingContext.parentModel.
						Interface().(IModel).
						getModel(),
				),
			})

			// restore the state of the context
			snapshot.restoreContext(storingContext)
		case scheme.HasManyEmbed:

			panic(errors.New("TODO: not implemented yet!!!!"))
		case scheme.Belongs:

			if field.Type == storingContext.parentScheme.Type && field.RelSrc == storingContext.name {
				operations = append(operations, change.Set{
					Field: field.Name,
					Value: reflect.ValueOf(
						storingContext.parentModel.
							Interface().(IModel).
							getModel(),
					),
				})
			} else {
				operations = append(operations, change.Set{
					Field: field.Name,
					Value: reflect.ValueOf(
						storingContext.model.FieldByIndex(field.Index).
							Interface().(IModel).
							getModel(),
					),
				})
			}

		}

	}

	// storingMode new will execute new in all fields
	if storingContext.mode == "creating" {

	} else if storingContext.mode == "updating" {

	} else {

	}

	storingContext.storeRelations()
	storingContext.storeThroughRelations()
}

func (storingContext *storingContext) storeThroughRelations() {

}

func (storingContext *storingContext) storeRelations() {
	// make a snapshot of the state of the storingContext
	snapshot := *storingContext
	for _, field := range storingContext.scheme.Fields() {
		switch field.RelKind {
		case scheme.HasOne:
			storingContext.parentModel = storingContext.model
			snapshot.restoreContext(storingContext)
		case scheme.HasMany:
			snapshot.restoreContext(storingContext)
		}
	}
}
