package wordtype

import (
	"testing"

	cases "github.com/marthjod/binquiry-experimental/case"
	"github.com/marthjod/binquiry-experimental/gender"
	"github.com/marthjod/binquiry-experimental/noun"
	"github.com/marthjod/binquiry-experimental/number"
)

var words = &Words{
	&noun.Noun{
		Gender: gender.Gender_Masculine,
		Cases: []*noun.CaseForm{
			{Case: cases.Case_Nominative, Number: number.Number_Singular, Form: "penni"},
			{Case: cases.Case_Accusative, Number: number.Number_Singular, Form: "penna"},
			{Case: cases.Case_Dative, Number: number.Number_Singular, Form: "penna"},
			{Case: cases.Case_Genitive, Number: number.Number_Singular, Form: "penna"},
			{Case: cases.Case_Nominative, Number: number.Number_Plural, Form: "pennar"},
			{Case: cases.Case_Accusative, Number: number.Number_Plural, Form: "penna"},
			{Case: cases.Case_Dative, Number: number.Number_Plural, Form: "pennum"},
			{Case: cases.Case_Genitive, Number: number.Number_Plural, Form: "penna"},
		},
	},
	&noun.Noun{
		Gender: gender.Gender_Feminine,
		Cases: []*noun.CaseForm{
			{Case: cases.Case_Nominative, Number: number.Number_Singular, Form: "kona"},
			{Case: cases.Case_Accusative, Number: number.Number_Singular, Form: "konu"},
			{Case: cases.Case_Dative, Number: number.Number_Singular, Form: "konu"},
			{Case: cases.Case_Genitive, Number: number.Number_Singular, Form: "konu"},
			{Case: cases.Case_Nominative, Number: number.Number_Plural, Form: "konur"},
			{Case: cases.Case_Accusative, Number: number.Number_Plural, Form: "konur"},
			{Case: cases.Case_Dative, Number: number.Number_Plural, Form: "konum"},
			{Case: cases.Case_Genitive, Number: number.Number_Plural, Form: "kvenna"},
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
