package matchers

import (
    "fmt"
    "github.com/dictybase/testchado"
    "github.com/onsi/gomega"
    "strings"
)

//HasCv checks for presence of a cv namespace in chado database
func HasCv(expected interface{}) gomega.OmegaMatcher {
    return &HasCvMatcher{expected: expected}
}

type HasCvMatcher struct {
    expected interface{}
}

func (matcher *HasCvMatcher) Match(actual interface{}) (success bool, message string, err error) {
    dbm, ok := actual.(testchado.DBManager)
    if !ok {
        return false, "", fmt.Errorf("HasCv matcher expects a testchado.DBManager")
    }
    cv, ok := matcher.expected.(string)
    if !ok {
        return false, "", fmt.Errorf("HasCv matcher expects a cv name")
    }

    type entries struct{ Counter int }
    e := entries{}
    sqlx := dbm.DBHandle()
    err = sqlx.Get(&e, "SELECT count(cv_id) counter FROM cv where name = $1", cv)
    if err != nil {
        return false, "", fmt.Errorf("could not execute query: %s", err)
    }
    if e.Counter == 1 {
        return true, fmt.Sprintf("Expected\n\tcv %#v does not exist in database", matcher.expected), nil
    }
    return false, fmt.Sprintf("Expected\n\tcv %#v exist in database", matcher.expected), nil
}

//HasCvterm check for presence of a cvterm in chado database.
func HasCvterm(expected interface{}) gomega.OmegaMatcher {
    return &HasCvtermMatcher{expected: expected}
}

type HasCvtermMatcher struct {
    expected interface{}
}

func (matcher *HasCvtermMatcher) Match(actual interface{}) (success bool, message string, err error) {
    dbm, ok := actual.(testchado.DBManager)
    if !ok {
        return false, "", fmt.Errorf("HasCvterm matcher expects a testchado.DBManager")
    }
    cvterm, ok := matcher.expected.(string)
    if !ok {
        return false, "", fmt.Errorf("HasCvterm matcher expects a cvterm")
    }

    type entries struct{ Counter int }
    e := entries{}
    sqlx := dbm.DBHandle()
    err = sqlx.Get(&e, "SELECT count(cvterm_id) counter FROM cvterm where name = $1", cvterm)
    if err != nil {
        return false, "", fmt.Errorf("could not execute query: %s", err)
    }
    if e.Counter > 0 {
        return true, fmt.Sprintf("Expected\n\tcvterm %#v does not exist in database", matcher.expected), nil
    }
    return false, fmt.Sprintf("Expected\n\tcvterm %#v exist in database", matcher.expected), nil
}

// HasDbxref check for presence of a xref in chado database. In case of xref in standard format(DB:Id),
// it splits and check for both id and db name.
func HasDbxref(expected interface{}) gomega.OmegaMatcher {
    return &HasDbxrefMatcher{expected: expected}
}

type HasDbxrefMatcher struct {
    expected interface{}
}

func (matcher *HasDbxrefMatcher) Match(actual interface{}) (success bool, message string, err error) {
    dbm, ok := actual.(testchado.DBManager)
    if !ok {
        return false, "", fmt.Errorf("HasDbxref matcher expects a testchado.DBManager")
    }
    dbxref, ok := matcher.expected.(string)
    if !ok {
        return false, "", fmt.Errorf("HasDbxref matcher expects a dbxref")
    }

    type entries struct{ Counter int }
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
            return false, "", fmt.Errorf("could not execute query: %s", err)
        }
    } else {
        err = sqlx.Get(&e, "SELECT count(dbxref_id) counter FROM dbxref WHERE accession = $1", dbxref)
        if err != nil {
            return false, "", fmt.Errorf("could not execute query: %s", err)
        }
    }
    if e.Counter > 0 {
        return true, fmt.Sprintf("Expected\n\tdbxref %#v does not exist in database", matcher.expected), nil
    }
    return false, fmt.Sprintf("Expected\n\tdbxref %#v exist in database", matcher.expected), nil
}
