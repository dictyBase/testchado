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
//  Expect("SELECT * FROM feature").Should(HaveRows(20))
func HaveRows(expected interface{}) gomega.OmegaMatcher {
	return &HaveRowsMatcher{expected: expected}
}

type HaveRowsMatcher struct {
	expected interface{}
}

func (matcher *HaveRowsMatcher) Match(actual interface{}) (success bool, err error) {
	query, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("HaveRows matcher expects a SQL query")
	}
	count, ok := matcher.expected.(int)
	if !ok {
		return false, fmt.Errorf("HaveRows matcher expects a integer value")
	}

	sqlx := dbmanager.DBHandle()
	rows, err := sqlx.Queryx(query)
	if err != nil {
		return false, fmt.Errorf("could not execute query: %s", err)
	}
	i := 0
	for rows.Next() {
		i += 1
	}
	if i == count {
		return true, nil
	}
	return false, nil
}

func (matcher *HaveRowsMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v to match \n\t%#vrows from database", actual, matcher.expected)
}

func (matcher *HaveRowsMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v not to match \n\t%#vrows from database", actual, matcher.expected)
}

//HaveCount is a variant of HaveRows where it matches the value of COUNT(*) or COUNT(column_name) SQL query in chado database
//      Expect("SELECT count(*) FROM dbxref").Should(HaveCount(23))
//      Expect("SELECT count(pub_id) FROM pub").Should(HaveCount(34))
func HaveCount(expected interface{}) gomega.OmegaMatcher {
	return &HaveCountMatcher{expected: expected}
}

type HaveCountMatcher struct {
	expected interface{}
}

func (matcher *HaveCountMatcher) Match(actual interface{}) (success bool, err error) {
	query, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("HaveCount matcher expects a SQL query")
	}
	count, ok := matcher.expected.(int)
	if !ok {
		return false, fmt.Errorf("HaveCount matcher expects a integer value")
	}

	sqlx := dbmanager.DBHandle()
	row := sqlx.QueryRowx(query)
	var dbcount int
	err = row.Scan(&dbcount)
	if err != nil {
		return false, fmt.Errorf("could not execute query: %s", err)
	}
	if dbcount == count {
		return true, nil
	}
	return false, nil
}

func (matcher *HaveCountMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v to match \n\t%#vcount from database", actual, matcher.expected)
}

func (matcher *HaveCountMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v not to match \n\t%#count from database", actual, matcher.expected)
}

//HaveNameCount is a variant of HaveNameCount with bind variables to run arbitary COUNT sql queries in chado database
/*The bind variables and expected count are passed through a map structure
  params : An interface slice containing bind values in order
  count  : The expected count from SQL statement

      query := `
          SELECT count(*) FROM feature JOIN organism ON
              feature.organism_id = organism.organism_id
              WHERE feature.is_obsolete = $1
              AND
              organism.genus = $2
              AND
              organism.species = $3
          `
      m := make(map[string]interface{})
      m["params"] = append(make([]interface, 0), 1, "Homo", "sapiens")
      m["count"] = 50
      Expect(query).Should(HaveNameCount(m))
*/
func HaveNameCount(expected interface{}) gomega.OmegaMatcher {
	return &HaveNameCountMatcher{expected: expected}
}

type HaveNameCountMatcher struct {
	expected interface{}
}

func (matcher *HaveNameCountMatcher) Match(actual interface{}) (success bool, err error) {
	query, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("HaveNameCount matcher expects a SQL query")
	}
	m, ok := matcher.expected.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("HaveNameCount matcher expects a map variable")
	}
	param, ok := m["params"]
	if !ok {
		return false, fmt.Errorf("The map variable does not have a param key for bind values")
	}
	args, ok := param.([]interface{})
	if !ok {
		return false, fmt.Errorf("Expecting slice of interface")
	}

	ct, ok := m["count"]
	if !ok {
		return false, fmt.Errorf("The map variable does not have a count key")
	}
	count, ok := ct.(int)
	if !ok {
		return false, fmt.Errorf("The count key does not have an integer count")
	}

	sqlx := dbmanager.DBHandle()
	row := sqlx.QueryRowx(query, args...)
	var dbcount int
	err = row.Scan(&dbcount)
	if err != nil {
		return false, fmt.Errorf("could not execute query: %s", err)
	}
	if dbcount == count {
		return true, nil
	}
	return false, nil
}

func (matcher *HaveNameCountMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v to match \n\t%#vcount from database", actual, matcher.expected)
}

func (matcher *HaveNameCountMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n\t%#v not to match \n\t%#count from database", actual, matcher.expected)
}
