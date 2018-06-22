package wordtype

import (
	"testing"
)

func TestGetWordType(t *testing.T) {
	var expected = []struct {
		in  string
		out WordType
	}{
		{
			in:  "nafnorð",
			out: WordType_Noun,
		},
		{
			in:  "lýsingarorð",
			out: WordType_Adjective,
		},
		{
			in:  "sagnorð",
			out: WordType_Verb,
		},
		{
			in:  "Eitthvað annað",
			out: WordType_Unknown,
		},
	}

	for _, exp := range expected {
		actual := GetWordType(exp.in)
		if exp.out != actual {
			t.Errorf("Expected: %v, actual: %v", exp.out, actual)
		}
	}
}
