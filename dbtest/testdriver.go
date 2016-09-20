package dbtest

import (
	"github.com/CloudyKit/srm/change"
	"github.com/CloudyKit/srm/driver"
	"github.com/CloudyKit/srm/query"
	"github.com/CloudyKit/srm/scheme"

	"bytes"
	"fmt"
	"reflect"
)

var _ = driver.Driver(&FakeDriver{})

type FakeRecord map[string]reflect.Value
type fakeTable map[string]FakeRecord
type fakeDB map[string]fakeTable

func NewFakeDriver() *FakeDriver {
	return &FakeDriver{}
}

type FakeDriver struct {
	driver.NoTransactions
	ids         int
	db          fakeDB
	oplog       bytes.Buffer

	PanicNew    bool
	PanicUpdate bool
	PanicDelete bool
	PanicModify bool
	PanicRemove bool
}

func (fk *FakeDriver) ResetOPLog() {
	fk.oplog.Reset()
}

func (fk *FakeDriver) OPLog() *bytes.Buffer {
	return &fk.oplog
}

func (d *FakeDriver) UseScheme(name string, s *scheme.Scheme) error {
	return nil
}

func (d *FakeDriver) Search(name string, s *scheme.Scheme, q *query.Query) driver.Result {
	return nil
}

func (d *FakeDriver) Retrieve(name string, s *scheme.Scheme, key string) driver.Result {
	return nil
}

func (d *FakeDriver) getTable(tableName string) fakeTable {

	if d.db == nil {
		d.db = make(fakeDB)
	}

	table, ok := d.db[tableName]
	if !ok {
		table = make(fakeTable)
		d.db[tableName] = table
	}

	return table
}

func (d *FakeDriver) printf(format string, v ...interface{}) {
	fmt.Fprintf(&d.oplog, format, v...)
}

func (d *FakeDriver) Create(name string, s *scheme.Scheme, operations ...change.Set) (key string, err error) {
	if d.PanicNew {
		panic(fmt.Errorf("Panic on New enabled"))
	}

	d.ids++
	key = fmt.Sprint(d.ids)

	d.printf("INSERT: table(%s) key(%s)", name, key)
	table := d.getTable(name)
	record := make(FakeRecord)

	for _, set := range operations {
		record[set.Field] = set.Value

		if set.Value.IsValid() {
			d.printf(" set(%s)=%q", set.Field, set.Value.Interface())
		} else {
			d.printf(" set(%s)=%q", set.Field, "")
		}
	}

	//record[keyField] = reflect.ValueOf(key)
	table[key] = record

	d.oplog.WriteString("\n")
	return
}

func (d *FakeDriver) getRecord(table, key string) (record FakeRecord, found bool) {
	record, found = d.getTable(table)[key]
	return
}

func (d *FakeDriver) Modify(name string, s *scheme.Scheme, key string, operations ...change.Set) (numofmodified int, err error) {
	if d.PanicUpdate {
		panic(fmt.Errorf("Panic on Update enabled"))
	}
	d.printf("UPDATE: table(%s) key(%s)", name, key)

	record, found := d.getRecord(name, key)
	if found {
		numofmodified++
		for _, set := range operations {
			record[set.Field] = set.Value
			if set.Value.IsValid() {
				d.printf(" set(%s)=%q", set.Field, set.Value.Interface())
			} else {
				d.printf(" set(%s)=%q", set.Field, "")
			}
		}
		//record[keyField] = reflect.ValueOf(key)
	} else {
		d.printf(" NOT FOUND")
	}
	d.oplog.WriteString("\n")
	return
}

func (d *FakeDriver) SearchAndModify(name string, s *scheme.Scheme, q *query.Query, operations ...change.Operation) (numofmodified int, err error) {
	if d.PanicModify {
		panic(fmt.Errorf("Panic on Modify enabled"))
	}

	//record, found := d.getRecord(s.Entity(), primaryKey)
	//
	//if found {
	//	numofmodified++
	//	for _, op := range operations {
	//		switch op := op.(type) {
	//		case change.Set:
	//			record[op.Field] = op.Value
	//		default:
	//			err = errors.New("operation is not supported")
	//		}
	//	}
	//}

	return
}

func (d *FakeDriver) Remove(name string, s *scheme.Scheme, key string) (numofmodified int, err error) {
	if d.PanicDelete {
		panic(fmt.Errorf("Panic on Delete enabled"))
	}

	_, found := d.getRecord(name, key)
	if found {
		numofmodified++
		delete(d.getTable(name), key)
	}
	return
}

func (d *FakeDriver) SearchAndRemove(name string, s *scheme.Scheme, q *query.Query) (numofmodified int, err error) {
	if d.PanicRemove {
		panic(fmt.Errorf("Panic on Remove enabled"))
	}

	return
}
