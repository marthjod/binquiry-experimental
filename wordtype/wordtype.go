package wordtype

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Word represents features every word type must exhibit.
type Word interface {
	JSON() string
	CanonicalForm() string
}

// Words is a list of Word types.
type Words []Word

var wordTypes = map[string]WordType{
	"nafnorð":     WordType_Noun,
	"lýsingarorð": WordType_Adjective,
	"sagnorð":     WordType_Verb,
}

// GetWordType determines a word type based on the input string.
func GetWordType(header string) WordType {
	for k, v := range wordTypes {
		if strings.Contains(strings.ToLower(header), k) {
			return v
		}
	}

	return WordType_Unknown
}

// JSON representation of Words.
func (w *Words) JSON() string {
	j, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(j)
}
