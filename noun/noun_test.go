package noun

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	cases "github.com/marthjod/binquiry-experimental/case"
	"github.com/marthjod/binquiry-experimental/gender"
	"github.com/marthjod/binquiry-experimental/number"
	"github.com/marthjod/binquiry-experimental/pkg/reader"
	xmlpath "gopkg.in/xmlpath.v2"
)

var (
	noun = Noun{
		Gender: gender.Gender_Masculine,
		Cases: []*CaseForm{
			{Case: cases.Case_Nominative, Number: number.Number_Singular, Form: "penni"},
			{Case: cases.Case_Accusative, Number: number.Number_Singular, Form: "penna"},
			{Case: cases.Case_Dative, Number: number.Number_Singular, Form: "penna"},
			{Case: cases.Case_Genitive, Number: number.Number_Singular, Form: "penna"},
			{Case: cases.Case_Nominative, Number: number.Number_Plural, Form: "pennar"},
			{Case: cases.Case_Accusative, Number: number.Number_Plural, Form: "penna"},
			{Case: cases.Case_Dative, Number: number.Number_Plural, Form: "pennum"},
			{Case: cases.Case_Genitive, Number: number.Number_Plural, Form: "penna"},
		},
	}
)

func TestNoun_Json(t *testing.T) {
	expected := `{
  "gender": "Masculine",
  "cases": [
    {
      "case": "Nominative",
      "number": "Singular",
      "form": "penni"
    },
    {
      "case": "Accusative",
      "number": "Singular",
      "form": "penna"
    },
    {
      "case": "Dative",
      "number": "Singular",
      "form": "penna"
    },
    {
      "case": "Genitive",
      "number": "Singular",
      "form": "penna"
    },
    {
      "case": "Nominative",
      "number": "Plural",
      "form": "pennar"
    },
    {
      "case": "Accusative",
      "number": "Plural",
      "form": "penna"
    },
    {
      "case": "Dative",
      "number": "Plural",
      "form": "pennum"
    },
    {
      "case": "Genitive",
      "number": "Plural",
      "form": "penna"
    }
  ]
}`

	actual := noun.JSON()
	if expected != actual {
		t.Errorf("Expected: %v, actual: %v", expected, actual)
	}

}

func ExampleNoun_JSON() {

	n := Noun{
		Gender: gender.Gender_Masculine,
		Cases: []*CaseForm{
			{Case: cases.Case_Nominative, Number: number.Number_Singular, Form: "penni"},
			{Case: cases.Case_Accusative, Number: number.Number_Singular, Form: "penna"},
			{Case: cases.Case_Dative, Number: number.Number_Singular, Form: "penna"},
			{Case: cases.Case_Genitive, Number: number.Number_Singular, Form: "penna"},
			{Case: cases.Case_Nominative, Number: number.Number_Plural, Form: "pennar"},
			{Case: cases.Case_Accusative, Number: number.Number_Plural, Form: "penna"},
			{Case: cases.Case_Dative, Number: number.Number_Plural, Form: "pennum"},
			{Case: cases.Case_Genitive, Number: number.Number_Plural, Form: "penna"},
		},
	}
	fmt.Println(n.JSON())
	// Output: {
	//   "gender": "Masculine",
	//   "cases": [
	//     {
	//       "case": "Nominative",
	//       "number": "Singular",
	//       "form": "penni"
	//     },
	//     {
	//       "case": "Accusative",
	//       "number": "Singular",
	//       "form": "penna"
	//     },
	//     {
	//       "case": "Dative",
	//       "number": "Singular",
	//       "form": "penna"
	//     },
	//     {
	//       "case": "Genitive",
	//       "number": "Singular",
	//       "form": "penna"
	//     },
	//     {
	//       "case": "Nominative",
	//       "number": "Plural",
	//       "form": "pennar"
	//     },
	//     {
	//       "case": "Accusative",
	//       "number": "Plural",
	//       "form": "penna"
	//     },
	//     {
	//       "case": "Dative",
	//       "number": "Plural",
	//       "form": "pennum"
	//     },
	//     {
	//       "case": "Genitive",
	//       "number": "Plural",
	//       "form": "penna"
	//     }
	//   ]
	// }
}

func TestParseNoun(t *testing.T) {
	expected := &noun
	f, err := os.Open("testdata/penni.xml")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	header, _, root, _ := reader.Read(f)
	path := xmlpath.MustCompile("//tr/td[2]")

	actual := ParseNoun(header, path.Iter(root))
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v,\nactual: %v", expected, actual)
	}
}
