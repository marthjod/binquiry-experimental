package wordtype

import (
	"github.com/marthjod/bingo/model/case"
	"github.com/marthjod/bingo/model/gender"
	"github.com/marthjod/bingo/model/noun"
	"github.com/marthjod/bingo/model/number"
	"reflect"
	"testing"
)

var words = &Words{
	&noun.Noun{
		Gender: gender.Masculine,
		CaseForms: []noun.CaseForm{
			{Case: cases.Nominative, Number: number.Singular, Form: "penni"},
			{Case: cases.Accusative, Number: number.Singular, Form: "penna"},
			{Case: cases.Dative, Number: number.Singular, Form: "penna"},
			{Case: cases.Genitive, Number: number.Singular, Form: "penna"},
			{Case: cases.Nominative, Number: number.Plural, Form: "pennar"},
			{Case: cases.Accusative, Number: number.Plural, Form: "penna"},
			{Case: cases.Dative, Number: number.Plural, Form: "pennum"},
			{Case: cases.Genitive, Number: number.Plural, Form: "penna"},
		},
	},
	&noun.Noun{
		Gender: gender.Feminine,
		CaseForms: []noun.CaseForm{
			{Case: cases.Nominative, Number: number.Singular, Form: "kona"},
			{Case: cases.Accusative, Number: number.Singular, Form: "konu"},
			{Case: cases.Dative, Number: number.Singular, Form: "konu"},
			{Case: cases.Genitive, Number: number.Singular, Form: "konu"},
			{Case: cases.Nominative, Number: number.Plural, Form: "konur"},
			{Case: cases.Accusative, Number: number.Plural, Form: "konur"},
			{Case: cases.Dative, Number: number.Plural, Form: "konum"},
			{Case: cases.Genitive, Number: number.Plural, Form: "kvenna"},
		},
	},
}

func TestGetWordType(t *testing.T) {
	var expected = []struct {
		in  string
		out WordType
	}{
		{
			in:  "nafnorð",
			out: Noun,
		},
		{
			in:  "lýsingarorð",
			out: Adjective,
		},
		{
			in:  "sagnorð",
			out: Verb,
		},
		{
			in:  "Eitthvað annað",
			out: Unknown,
		},
	}

	for _, exp := range expected {
		actual := GetWordType(exp.in)
		if exp.out != actual {
			t.Errorf("Expected: %v, actual: %v", exp.out, actual)
		}
	}
}

func TestWords_ParseTemplate(t *testing.T) {
	expected := []byte(`penni (Masculine), 8 forms
kona (Feminine), 8 forms`)

	actual, err := words.ParseTemplate(`{{ range $i, $e := . }}{{ if $i }}{{ println }}{{ end }}{{ (index .CaseForms 0).Form }} ({{ .Gender }}), {{ .CaseForms | len }} forms{{ end }}`)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: --%v--,\nactual: --%v--", string(expected), string(actual))
	}
}
