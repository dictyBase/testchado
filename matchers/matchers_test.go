package matchers

import (
    "github.com/dictybase/testchado"
    . "github.com/onsi/gomega"
    "testing"
)

func TestCommonMatcher(t *testing.T) {
    RegisterTestingT(t)
    chado := testchado.NewDBManager()
    chado.DeploySchema()
    chado.LoadDefaultFixture()
    defer chado.DropSchema()

    Expect(chado).Should(HaveCv("sequence"))
    Expect(chado).ShouldNot(HaveCv("gene_ontology"))

    for _, name := range []string{"gene", "match_part", "has_agent"} {
        Expect(chado).Should(HaveCvterm(name))
    }
    for _, name := range []string{"perl", "golang", "python"} {
        Expect(chado).ShouldNot(HaveCvterm(name))
    }

    for _, dbxref := range []string{"sequence:variant_of", "relationship:OBO_REL:transformed_into", "member_of"} {
        Expect(chado).Should(HaveDbxref(dbxref))
    }

    for _, organism := range []string{"dicty", "frog", "rice", "mouse"} {
        Expect(chado).Should(HaveOrganism(organism))
    }
}

func TestDatabaseMatchers(t *testing.T) {
    RegisterTestingT(t)
    chado := testchado.NewDBManager()
    RegisterDBHandler(chado)
    chado.DeploySchema()
    chado.LoadDefaultFixture()
    defer chado.DropSchema()

    q := "SELECT count(*) FROM organism"
    Expect(q).Should(HaveCount(12))

    q = "SELECT * FROM dbxref JOIN db ON dbxref.db_id = db.db_id where db.name = 'PMID'"
    Expect(q).Should(HaveRows(6))

    query := `
     SELECT count(cvterm.cvterm_id) counter from CVTERM join CV on CV.CV_ID=CVTERM.CV_ID
     WHERE CV.NAME = $1
    `

    m := make(map[string]interface{})
    m["params"] = append(make([]interface{}, 0), "sequence")
    m["count"] = 286
    Expect(query).Should(HaveNameCount(m))
}
