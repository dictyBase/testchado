/*
Package testchado is a golang library to write unit tests for chado
database based packages and applications.
Chado is an open source modular database schema for storing biological
data. Package testchado provides a resuable API to write test cases for
softwares that uses chado database schema for storage. It supports two
RDBMS backends, postgresql and sqlite and uses custom chado matchers
based on gomega (http://onsi.github.io/gomega/#adding_your_own_matchers)
package.

Quick Start

    package chado_quickstart

    import (
        "github.com/dictybase/testchado"
        "testing"
      . "github.com/onsi/gomega"
    )

    func TestQuickStart (t *testing.T) {
        //gomega setup
        RegisterTestingT(t)

        //setup
        chado := NewDBManager()
        chado.LoadDefaultFixtures()
        chado.DeploySchema()
        //teardown
        defer chado.DropSchema()

        //matchers for testing
        Expect(chado).Should(HaveCv("sequence"))
        for _, name := range []string{"gene", "match_part", "has_agent"} {
            Expect(chado).Should(HaveCvterm(name))
        }
        Expect(chado).Should(HaveDbxref("SO:0000704"))
    }

Run against a sqlite backend.

    go test

To run against an postgresql backend set the TC_DSOURCE variable.

    TC_DSOURCE="dbname=chado user=chado password=chado host=localhost sslmode=disable"
                                \ go test



Testing Arbitary SQL

Though the custom chado matchers are easy to get started, however they could
not complement running arbitary SQL. The following matchers allows to test
number of rows from adhoc SQL statements with bind parameters.


    package chado_adhoc

    import (
        "github.com/dictybase/testchado"
        "testing"
      . "github.com/onsi/gomega"
    )

    func TestQuickStart (t *testing.T) {
        //gomega setup
        RegisterTestingT(t)

        //setup
        chado := NewDBManager()
        RegisterDBHandler(chado)
        chado.DeploySchema()
        chado.LoadDefaultFixtures()
        //teardown
        defer chado.DropSchema()

        //matchers for testing
        q := "SELECT count(*) FROM cvterm"
        Expect(q).Should(HaveCount(13))

        q = "SELECT * FROM organism"
        Expect(q).Should(HaveRows(12))


        //named parameters
        query := `
            SELECT count(*) counter from CVTERM join CV on CV.CV_ID=CVTERM.CV_ID
            WHERE CV.NAME = $1 AND CVTERM.IS_OBSOLETE = $2
        `
        m := make(map[string]interface{})
        m["params"] = append(make([]interface{},0), "sequence", 8)
        m["count"] = 8
        Expect(query).Should(HaveNameCount(m))
    }


Loading Fixtures

Package testchado provides supports loading database fixtures. Currently, it
expects a flat file containing multiple INSERT statements for loading data.
However, it could be also any custom SQL statements that populates data, for
example, COPY statements for postgresql backend. There are three methods for
handling fixtures ..

        LoadDefaultFixture() // basic and general fixture, should work for
        most cases
        LoadPresetFixture("cvprop") // Either of cvprop or eco
        LoadCustomFixture("path") // A file containing SQL statements


*/

package testchado
