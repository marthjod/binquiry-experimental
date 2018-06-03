package noun

import (
	"encoding/json"
	"fmt"

	"github.com/marthjod/bingo/model/case"
	"github.com/marthjod/bingo/model/gender"
	"github.com/marthjod/bingo/model/number"
	"gopkg.in/xmlpath.v2"
)

// CaseForm represents a single case form, i.e. case name, number, and actual form.
type CaseForm struct {
	Case   cases.Case    `json:"case"`
	Number number.Number `json:"number"`
	Form   string        `json:"form"`
}

// Noun is defined as a combination of a gender and a list of case forms.
type Noun struct {
	Gender    gender.Gender `json:"gender"`
	CaseForms []CaseForm    `json:"cases"`
}

func FromJSON(b []byte) (Noun, error) {
	var n Noun
	err := json.Unmarshal(b, &n)
	if err != nil {
		return Noun{}, fmt.Errorf("%s %s", b, err)
	}
	return n, nil
}

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
			n.CaseForms = append(n.CaseForms, CaseForm{
				Case:   cases.Nominative,
				Number: number.Singular,
				Form:   node.String(),
			})
		case 2:
			n.CaseForms = append(n.CaseForms, CaseForm{
				Case:   cases.Accusative,
				Number: number.Singular,
				Form:   node.String(),
			})
		case 3:
			n.CaseForms = append(n.CaseForms, CaseForm{
				Case:   cases.Dative,
				Number: number.Singular,
				Form:   node.String(),
			})
		case 4:
			n.CaseForms = append(n.CaseForms, CaseForm{
				Case:   cases.Genitive,
				Number: number.Singular,
				Form:   node.String(),
			})
		case 5:
			n.CaseForms = append(n.CaseForms, CaseForm{
				Case:   cases.Nominative,
				Number: number.Plural,
				Form:   node.String(),
			})
		case 6:
			n.CaseForms = append(n.CaseForms, CaseForm{
				Case:   cases.Accusative,
				Number: number.Plural,
				Form:   node.String(),
			})
		case 7:
			n.CaseForms = append(n.CaseForms, CaseForm{
				Case:   cases.Dative,
				Number: number.Plural,
				Form:   node.String(),
			})
		case 8:
			n.CaseForms = append(n.CaseForms, CaseForm{
				Case:   cases.Genitive,
				Number: number.Plural,
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
	for _, f := range n.CaseForms {
		if f.Case == cases.Nominative && f.Number == number.Singular {
			return f.Form
		}
	}
	// TODO this should not be reached
	return ""
}

func (n Noun) String() string {
	return n.CanonicalForm()
}
