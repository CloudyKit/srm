package srm

import (
	"errors"
	"github.com/CloudyKit/srm/scheme"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	scmMap = map[reflect.Type]*scheme.Scheme{}
	scmRwMutex = sync.RWMutex{}
)

func getScheme(t reflect.Type) *scheme.Scheme {

	scmRwMutex.RLock()
	_scheme, found := scmMap[t]
	scmRwMutex.RUnlock()

	if !found {
		panic(errors.New("Trying to use an undefined Scheme! call database.Scheme(Type,func("))
	}
	return _scheme
}

func Scheme(typ IModel, defFn func(*scheme.Def)) *scheme.Def {
	scmRwMutex.Lock()
	defer scmRwMutex.Unlock()

	t := reflect.TypeOf(typ)

	scm, found := scmMap[t]
	if found {
		panic(errors.New("Scheme redefinition is not permited!"))
	}

	scm = &scheme.Scheme{Type: t}

	// interprets the new allocated scheme as type *Def
	scmDef := (*scheme.Def)(scm)

	defFn(scmDef)
	scmDef.Done()
	scmMap[t] = scm

	return scm
}
