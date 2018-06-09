package gender

import "strings"

var genders = map[string]Gender{
	"karlkyn":   Gender_Masculine,
	"kvenkyn":   Gender_Feminine,
	"hvorugkyn": Gender_Neuter,
}

// GetGender determines a Gender type based on the input string.
func GetGender(header string) Gender {
	for k, v := range genders {
		if strings.Contains(strings.ToLower(header), k) {
			return v
		}
	}

	return Gender_Unknown
}
