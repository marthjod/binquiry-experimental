package word

import (
	"encoding/json"
	"fmt"

	"github.com/marthjod/binquiry-experimental/wordtype"
)

// Words is a list of Word types.
type Words []*Word

// CanonicalForm returns a word's canonical form.
func (w Word) CanonicalForm() string {
	switch w.Type {
	case wordtype.WordType_Noun:
		return w.Noun.CanonicalForm()
	}
	return ""
}

// JSON representation of Words.
func (w Words) JSON() string {
	j, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(j)
}
