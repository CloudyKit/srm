package srm_test

import (
	"github.com/CloudyKit/srm"
	"github.com/CloudyKit/srm/dbtest"
	"github.com/CloudyKit/srm/scheme"
	"github.com/CloudyKit/framework/validation"
	"testing"
)

type (
	Company struct {
		srm.Model
		Name      string

		Employers []*Employee
	}

	Employee struct {
		srm.Model

		Company  *Company
		Name     string

		Manager  *Employee
		Managing []*Employee
		Works    []*Work
	}

	Work struct {
		srm.Model

		Description string
		Employee    *Employee
	}

	Owner struct {
		srm.Model
		Name string

		Cats []*Cat
	}

	Cat struct {
		srm.Model

		Name  string
		Owner *Owner
	}
)

var (
	_ = srm.Scheme(&Company{}, func(d *scheme.Def) {
		d.Field("Name").Validation(
			validation.MinLength("Min Length", 5),
		)
		d.Field("Employers").HasMany("Company")
	})

	_ = srm.Scheme(
		&Employee{},
		func(d *scheme.Def) {

			d.Field("Name").Validation(
				validation.MinLength("Min Length", 5),
			)

			d.Field("Company").Belongs("Id")
			d.Field("Managing").HasOne("Manager")
			d.Field("Manager").Belongs("Managing")
			d.Field("Works").HasMany("Employee")

		})

	_ = srm.Scheme(
		&Work{},
		func(d *scheme.Def) {
			d.Field("Description").Validation(
				validation.MinLength("Min length", 5),
				validation.MaxLength("Max length", 500),
			)
			d.Field("Employee").Belongs()
		})

	_ = srm.Scheme(
		&Cat{},
		func(d *scheme.Def) {
			d.Field("Name")
			d.Field("Owner").
				Belongs()

		})

	_ = srm.Scheme(
		&Owner{},
		func(d *scheme.Def) {
			d.Field("Name")
			d.Field("Cats").
				HasMany("Owner")
		})
)

func TestDB_SaveHaveOne(t *testing.T) {
	fakeDB := dbtest.NewFakeDB()
	cat := &Cat{
		Name:  "Whiskers",
		Owner: &Owner{Name: "David"},
	}
	fakeDB.Store(cat)

	w, g := fakeDB.OPLogExpect("INSERT: table(owner) key(1) set(Name)=\"David\"\nINSERT: table(cat) key(2) set(Name)=\"Whiskers\" set(owner_id)=\"1\"")
	if w != g {
		t.Errorf("OPLog mismatch: want:\n%s\n got:\n%s\n", w, g)
	}

	fakeDB.ResetOPLog()

	fakeDB.Store(&Cat{
		Name: "Whiskers brother",
	})

	w, g = fakeDB.OPLogExpect("INSERT: table(cat) key(3) set(Name)=\"Whiskers brother\" set(owner_id)=\"\"")
	if w != g {
		t.Errorf("OPLog mismatch: want:\n%s\n got:\n%s\n", w, g)
	}
	fakeDB.ResetOPLog()

	cat.Owner = nil
	fakeDB.Store(cat)

	w, g = fakeDB.OPLogExpect("UPDATE: table(cat) key(2) set(Name)=\"Whiskers\" set(owner_id)=\"\"")
	if w != g {
		t.Errorf("OPLog mismatch: want:\n%s\n got:\n%s\n", w, g)
	}
}

func TestDB_SaveHaveMany(t *testing.T) {

	company := &Company{
		Name: "Cloudy Kit",
		Employers: []*Employee{
			{Name: "Henrique"},
			{Name: "Henrique"},
			{Name: "Henrique", Works: []*Work{{Description: "Clean the shop"}}},
			{Name: "Henrique"},
		},
	}

	company.Employers[0].Managing = company.Employers[1:]

	fakeDB := dbtest.NewFakeDB()

	result, err := fakeDB.Store(company)
	if err != nil {
		t.Fatalf("Something bad happend can't insert the records to the database, err: %s", err)
	}

	if result.Bad() {
		t.Fatalf("Something bad happend can't insert the records to the database, validation: %+v", result)
	}

	want, got := fakeDB.OPLogExpect("\tINSERT: table(company) key(1) set(Name)=\"Cloudy Kit\"\nINSERT: table(employee) key(2) set(Name)=\"Henrique\" set(inchargeid)=\"\" set(companyid)=\"1\"\nINSERT: table(employee) key(3) set(Name)=\"Henrique\" set(inchargeid)=\"2\" set(companyid)=\"\"\nINSERT: table(employee) key(4) set(Name)=\"Henrique\" set(inchargeid)=\"2\" set(companyid)=\"\"\nINSERT: table(work) key(5) set(Description)=\"Clean the shop\" set(employeeid)=\"4\"\nINSERT: table(employee) key(6) set(Name)=\"Henrique\" set(inchargeid)=\"2\" set(companyid)=\"\"\nUPDATE: table(employee) key(3) set(Name)=\"Henrique\" set(inchargeid)=\"2\" set(companyid)=\"1\"\nUPDATE: table(employee) key(4) set(Name)=\"Henrique\" set(inchargeid)=\"2\" set(companyid)=\"1\"\nUPDATE: table(work) key(5) set(Description)=\"Clean the shop\" set(employeeid)=\"4\"\nUPDATE: table(employee) key(6) set(Name)=\"Henrique\" set(inchargeid)=\"2\" set(companyid)=\"1\"")
	if got != want {
		t.Fatalf("OPLog mismatch\nWant:\n%s\nGot:\n%s", want, got)
	}

	fakeDB.FakeDriver.PanicUpdate = false
	fakeDB.ResetOPLog()
	for _, e := range company.Employers {
		if e.Company != company {
			t.Error("Employee Company should be pointing to the parent Company")
		}
		e.Company = nil
		fakeDB.Store(e)
	}

	want, got = fakeDB.OPLogExpect("UPDATE: table(EmployeeScheme) key(2) set(Name)=\"Henrique\" set(CompanyID)=\"\" set(InCharge)=\"\"\nUPDATE: table(EmployeeScheme) key(3) set(Name)=\"Henrique\" set(CompanyID)=\"\" set(InCharge)=\"\"\nUPDATE: table(EmployeeScheme) key(4) set(Name)=\"Henrique\" set(CompanyID)=\"\" set(InCharge)=\"\"\nUPDATE: table(WorkScheme) key(5) set(Description)=\"Clean the shop\" set(EmployeeID)=\"4\"\nUPDATE: table(EmployeeScheme) key(6) set(Name)=\"Henrique\" set(CompanyID)=\"\" set(InCharge)=\"\"")
	if want != got {
		t.Fatalf("OPLog mismatch want:\n%s\ngot:\n%s", want, got)
	}

}
