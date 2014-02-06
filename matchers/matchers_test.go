package matchers

import (
    "github.com/dictybase/testchado"
    . "github.com/onsi/gomega"
    "testing"
)

func TestCommonMatcher(t *testing.T) {
    RegisterTestingT(t)
    chado := testchado.NewChadoSchema()
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
}
