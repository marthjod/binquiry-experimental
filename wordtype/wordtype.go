package wordtype

import (
	"strings"
)

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
