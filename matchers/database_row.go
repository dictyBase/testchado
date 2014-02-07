package matchers

import (
    "fmt"
    "github.com/dictybase/testchado"
    "github.com/onsi/gomega"
)

var dbmanager testchado.DBManager

func RegisterDBHandler(dbm testchado.DBManager) {
    dbmanager = dbm
}

//HaveRows matches the number of rows returned from arbitary SQL query in chado database
func HaveRows(expected interface{}) gomega.OmegaMatcher {
    return &HaveRowsMatcher{expected: expected}
}

type HaveRowsMatcher struct {
    expected interface{}
}

func (matcher *HaveRowsMatcher) Match(actual interface{}) (success bool, message string, err error) {
    query, ok := actual.(string)
    if !ok {
        return false, "", fmt.Errorf("HaveRows matcher expects a SQL query")
    }
    count, ok := matcher.expected.(int)
    if !ok {
        return false, "", fmt.Errorf("HaveRows matcher expects a integer value")
    }

    sqlx := dbmanager.DBHandle()
    rows, err := sqlx.Queryx(query)
    if err != nil {
        return false, "", fmt.Errorf("could not execute query: %s", err)
    }
    i := 0
    for rows.Next() {
        i += 1
    }
    if i == count {
        return true, fmt.Sprintf("Expected\n\t%#v rows: got %d from database", matcher.expected, i), nil
    }
    return false, fmt.Sprintf("Expected\n\trows %#v matches rows from database", matcher.expected), nil
}

//HaveCount is a variant of HaveRows where it matches the value of COUNT(*) or COUNT(column_name) SQL query in chado database
func HaveCount(expected interface{}) gomega.OmegaMatcher {
    return &HaveCountMatcher{expected: expected}
}

type HaveCountMatcher struct {
    expected interface{}
}

func (matcher *HaveCountMatcher) Match(actual interface{}) (success bool, message string, err error) {
    query, ok := actual.(string)
    if !ok {
        return false, "", fmt.Errorf("HaveCount matcher expects a SQL query")
    }
    count, ok := matcher.expected.(int)
    if !ok {
        return false, "", fmt.Errorf("HaveCount matcher expects a integer value")
    }

    sqlx := dbmanager.DBHandle()
    row := sqlx.QueryRowx(query)
    var dbcount int
    err = row.Scan(&dbcount)
    if err != nil {
        return false, "", fmt.Errorf("could not execute query: %s", err)
    }
    if dbcount == count {
        return true, fmt.Sprintf("Expected\n\t%#v count got %d from database", matcher.expected, dbcount), nil
    }
    return false, fmt.Sprintf("Expected\n\tcount %#v matches from database", matcher.expected), nil
}

//HaveNameCount is a variant of HaveNameCount with bind variables to run arbitary COUNT sql queries in chado database
func HaveNameCount(expected interface{}) gomega.OmegaMatcher {
    return &HaveNameCountMatcher{expected: expected}
}

type HaveNameCountMatcher struct {
    expected interface{}
}

func (matcher *HaveNameCountMatcher) Match(actual interface{}) (success bool, message string, err error) {
    query, ok := actual.(string)
    if !ok {
        return false, "", fmt.Errorf("HaveNameCount matcher expects a SQL query")
    }
    m, ok := matcher.expected.(map[string]interface{})
    if !ok {
        return false, "", fmt.Errorf("HaveNameCount matcher expects a map variable")
    }
    param, ok := m["params"]
    if !ok {
        return false, "", fmt.Errorf("The map variable does not have a param key for bind values")
    }
    args, ok := param.([]interface{})
    if !ok {
        return false, "", fmt.Errorf("Expecting slice of interface")
    }

    ct, ok := m["count"]
    if !ok {
        return false, "", fmt.Errorf("The map variable does not have a count key")
    }
    count, ok := ct.(int)
    if !ok {
        return false, "", fmt.Errorf("The count key does not have an integer count")
    }

    sqlx := dbmanager.DBHandle()
    row := sqlx.QueryRowx(query, args...)
    var dbcount int
    err = row.Scan(&dbcount)
    if err != nil {
        return false, "", fmt.Errorf("could not execute query: %s", err)
    }
    if dbcount == count {
        return true, fmt.Sprintf("Expected\n\t%d count got %d from database", count, dbcount), nil
    }
    return false, fmt.Sprintf("Expected\n\tcount %d matches from database", dbcount), nil
}
