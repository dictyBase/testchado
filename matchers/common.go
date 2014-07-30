package matchers

import (
	"fmt"
	"strings"

	"github.com/dictybase/testchado"
	"github.com/onsi/gomega"
)

type entries struct{ Counter int }

//HaveCv matches a cv namespace in chado database
//	chado := testchado.NewDBManager()
//  Expect(chado).Should(HaveCv("sequence"))
//  Expect(chado).ShouldNot(HaveCv("gene_ontology"))
func HaveCv(expected interface{}) gomega.OmegaMatcher {
	return &HaveCvMatcher{expected: expected}
}

type HaveCvMatcher struct {
	expected interface{}
}

func (matcher *HaveCvMatcher) Match(actual interface{}) (success bool, err error) {
	dbm, ok := actual.(testchado.DBManager)
	if !ok {
		return false, fmt.Errorf("HaveCv matcher expects a testchado.DBManager")
	}
	cv, ok := matcher.expected.(string)
	if !ok {
		return false, fmt.Errorf("HaveCv matcher expects a cv name")
	}

	//type entries struct{ Counter int }
	e := entries{}
	sqlx := dbm.DBHandle()
	err = sqlx.Get(&e, "SELECT count(cv_id) counter FROM cv where name = $1", cv)
	if err != nil {
		return false, fmt.Errorf("could not execute query: %s", err)
	}
	if e.Counter == 1 {
		return true, nil
	}
	return false, nil
}

func (matcher *HaveCvMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%cv #v to exist in database", matcher.expected)
}

func (matcher *HaveCvMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%cv #v not to exist in database", matcher.expected)
}

//HaveCvterm matches cvterm in chado database.
//	chado := testchado.NewDBManager()
//	Expect(chado).Should(HaveCvterm("match_part"))
func HaveCvterm(expected interface{}) gomega.OmegaMatcher {
	return &HaveCvtermMatcher{expected: expected}
}

type HaveCvtermMatcher struct {
	expected interface{}
}

func (matcher *HaveCvtermMatcher) Match(actual interface{}) (success bool, err error) {
	dbm, ok := actual.(testchado.DBManager)
	if !ok {
		return false, fmt.Errorf("HaveCvterm matcher expects a testchado.DBManager")
	}
	cvterm, ok := matcher.expected.(string)
	if !ok {
		return false, fmt.Errorf("HaveCvterm matcher expects a cvterm")
	}

	//type entries struct{ Counter int }
	e := entries{}
	sqlx := dbm.DBHandle()
	err = sqlx.Get(&e, "SELECT count(cvterm_id) counter FROM cvterm where name = $1", cvterm)
	if err != nil {
		return false, fmt.Errorf("could not execute query: %s", err)
	}
	if e.Counter > 0 {
		return true, nil
	}
	return false, nil
}

func (matcher *HaveCvtermMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%cvterm #v to exist in database", matcher.expected)
}

func (matcher *HaveCvtermMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%cvterm #v not to exist in database", matcher.expected)
}

// HaveDbxref matches xref in chado database. In case of xref in standard format(DB:Id),
// it splits and check for both id and db name.
//	chado := testchado.NewDBManager()
//	Expect(chado).Should(HaveDbxref("sequence:variant_of"))
//	Expect(chado).Should(HaveDbxref("member_of"))
func HaveDbxref(expected interface{}) gomega.OmegaMatcher {
	return &HaveDbxrefMatcher{expected: expected}
}

type HaveDbxrefMatcher struct {
	expected interface{}
}

func (matcher *HaveDbxrefMatcher) Match(actual interface{}) (success bool, err error) {
	dbm, ok := actual.(testchado.DBManager)
	if !ok {
		return false, fmt.Errorf("HaveDbxref matcher expects a testchado.DBManager")
	}
	dbxref, ok := matcher.expected.(string)
	if !ok {
		return false, fmt.Errorf("HaveDbxref matcher expects a dbxref")
	}

	//type entries struct{ Counter int }
	e := entries{}
	sqlx := dbm.DBHandle()
	if strings.Contains(dbxref, ":") {
		d := strings.SplitN(dbxref, ":", 2)
		q := `
        SELECT count(dbxref_id) counter FROM dbxref JOIN db
        ON dbxref.db_id = db.db_id
        WHERE dbxref.accession = $1
        AND db.name = $2
        `
		err = sqlx.Get(&e, q, d[1], d[0])
		if err != nil {
			return false, fmt.Errorf("could not execute query: %s", err)
		}
	} else {
		err = sqlx.Get(&e, "SELECT count(dbxref_id) counter FROM dbxref WHERE accession = $1", dbxref)
		if err != nil {
			return false, fmt.Errorf("could not execute query: %s", err)
		}
	}
	if e.Counter > 0 {
		return true, nil
	}
	return false, nil
}

func (matcher *HaveDbxrefMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%dbxref #v to exist in database", matcher.expected)
}

func (matcher *HaveDbxrefMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%dbxref #v not to exist in database", matcher.expected)
}

// HaveFeature matches uniquename of a feature in chado database.
//	chado := testchado.NewDBManager()
//	Expect(chado).Should(HaveFeature("sadA"))
func HaveFeature(expected interface{}) gomega.OmegaMatcher {
	return &HaveFeatureMatcher{expected: expected}
}

type HaveFeatureMatcher struct {
	expected interface{}
}

func (matcher *HaveFeatureMatcher) Match(actual interface{}) (success bool, err error) {
	dbm, ok := actual.(testchado.DBManager)
	if !ok {
		return false, fmt.Errorf("HaveFeature matcher expects a testchado.DBManager")
	}
	feature, ok := matcher.expected.(string)
	if !ok {
		return false, fmt.Errorf("HaveFeature matcher expects a feature")
	}

	//type entries struct{ Counter int }
	e := entries{}
	sqlx := dbm.DBHandle()
	err = sqlx.Get(&e, "SELECT count(feature_id) counter FROM feature where uniquename = $1", feature)
	if err != nil {
		return false, fmt.Errorf("could not execute query: %s", err)
	}
	if e.Counter > 1 {
		return true, nil
	}
	return false, nil
}

func (matcher *HaveFeatureMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%feature #v to exist in database", matcher.expected)
}

func (matcher *HaveFeatureMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%feature #v not to exist in database", matcher.expected)
}

// HaveOrganism matches common name(common_name column) of an organism in chado database.
//	chado := testchado.NewDBManager()
//	Expect(chado).Should(HaveOrganism("mouse"))
func HaveOrganism(expected interface{}) gomega.OmegaMatcher {
	return &HaveOrganismMatcher{expected: expected}
}

type HaveOrganismMatcher struct {
	expected interface{}
}

func (matcher *HaveOrganismMatcher) Match(actual interface{}) (success bool, err error) {
	dbm, ok := actual.(testchado.DBManager)
	if !ok {
		return false, fmt.Errorf("HaveOrganism matcher expects a testchado.DBManager")
	}
	organism, ok := matcher.expected.(string)
	if !ok {
		return false, fmt.Errorf("HaveOrganism matcher expects a organism")
	}

	//type entries struct{ Counter int }
	e := entries{}
	sqlx := dbm.DBHandle()
	err = sqlx.Get(&e, "SELECT count(organism_id) counter FROM organism where common_name = $1", organism)
	if err != nil {
		return false, fmt.Errorf("could not execute query: %s", err)
	}
	if e.Counter > 0 {
		return true, nil
	}
	return false, nil
}
func (matcher *HaveOrganismMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%organism #v to exist in database", matcher.expected)
}

func (matcher *HaveOrganismMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%organism #v not to exist in database", matcher.expected)
}
