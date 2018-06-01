package noun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

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

// List of a Noun's forms.
func (n Noun) List() []string {
	l := []string{}
	for _, c := range n.CaseForms {
		l = append(l, c.Form)
	}

	return l
}

// ParseTemplate returns a parsed template based on Noun fields.
func (n Noun) ParseTemplate(tpl string) ([]byte, error) {
	var buf bytes.Buffer

	tmpl, err := template.New("").Parse(tpl)
	if err != nil {
		return []byte{}, err
	}

	err = tmpl.Execute(&buf, n)
	return buf.Bytes(), err
}
