package noun

import (
	"encoding/json"
	"fmt"

	cases "github.com/marthjod/binquiry-experimental/case"
	gender "github.com/marthjod/binquiry-experimental/gender"
	number "github.com/marthjod/binquiry-experimental/number"
	xmlpath "gopkg.in/xmlpath.v2"
)

// ParseNoun parses XML input into a Noun struct.
func ParseNoun(header string, iter *xmlpath.Iter) *Noun {
	n := Noun{
		Gender: gender.GetGender(header),
	}
	count := 1
	for iter.Next() {
		node := iter.Node()
		switch count {
		case 1:
			n.Cases = append(n.Cases, &CaseForm{
				Case:   cases.Case_Nominative,
				Number: number.Number_Singular,
				Form:   node.String(),
			})
		case 2:
			n.Cases = append(n.Cases, &CaseForm{
				Case:   cases.Case_Accusative,
				Number: number.Number_Singular,
				Form:   node.String(),
			})
		case 3:
			n.Cases = append(n.Cases, &CaseForm{
				Case:   cases.Case_Dative,
				Number: number.Number_Singular,
				Form:   node.String(),
			})
		case 4:
			n.Cases = append(n.Cases, &CaseForm{
				Case:   cases.Case_Genitive,
				Number: number.Number_Singular,
				Form:   node.String(),
			})
		case 5:
			n.Cases = append(n.Cases, &CaseForm{
				Case:   cases.Case_Nominative,
				Number: number.Number_Plural,
				Form:   node.String(),
			})
		case 6:
			n.Cases = append(n.Cases, &CaseForm{
				Case:   cases.Case_Accusative,
				Number: number.Number_Plural,
				Form:   node.String(),
			})
		case 7:
			n.Cases = append(n.Cases, &CaseForm{
				Case:   cases.Case_Dative,
				Number: number.Number_Plural,
				Form:   node.String(),
			})
		case 8:
			n.Cases = append(n.Cases, &CaseForm{
				Case:   cases.Case_Genitive,
				Number: number.Number_Plural,
				Form:   node.String(),
			})
		}
		count++
	}

	return &n
}

// JSON representation of a Noun.
func (n Noun) JSON() string {
	j, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(j)
}

// CanonicalForm returns a noun's canonical form, ie. nominative singular.
func (n Noun) CanonicalForm() string {
	for _, f := range n.Cases {
		if f.Case == cases.Case_Nominative && f.Number == number.Number_Singular {
			return f.Form
		}
	}
	// TODO this should not be reached
	return ""
}
