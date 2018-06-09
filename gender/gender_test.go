package gender

import (
	"testing"
)

var expected = []struct {
	in  string
	out Gender
}{
	{
		in:  "Hvorugkynsnafnorð",
		out: Gender_Neuter,
	},
	{
		in:  " Kvenkyn ",
		out: Gender_Feminine,
	},
	{
		in:  "Eitthvað annað",
		out: Gender_Unknown,
	},
}

func TestGetGender(t *testing.T) {
	for _, exp := range expected {
		actual := GetGender(exp.in)
		if exp.out != actual {
			t.Errorf("Expected: %v, actual: %v", exp.out, actual)
		}
	}
}
