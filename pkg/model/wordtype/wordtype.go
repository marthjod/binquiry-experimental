package wordtype

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"
)

// Word represents features every word type must exhibit.
type Word interface {
	JSON() string
	List() []string
	ParseTemplate(tpl string) ([]byte, error)
}

// Words is a list of Word types.
type Words []Word

// WordType is an enum representing word types.
type WordType int

//go:generate jsonenums -type=WordType
//go:generate stringer -type=WordType
const (
	Noun      WordType = iota
	Adjective WordType = iota
	Verb      WordType = iota
	Unknown   WordType = iota
)

var wordTypes = map[string]WordType{
	"nafnorð":     Noun,
	"lýsingarorð": Adjective,
	"sagnorð":     Verb,
}

// GetWordType determines a word type based on the input string.
func GetWordType(header string) WordType {
	for k, v := range wordTypes {
		if strings.Contains(strings.ToLower(header), k) {
			return v
		}
	}

	return Unknown
}

func Determine(r io.Reader) WordType {
	s := bufio.NewScanner(r)
	for s.Scan() {
		t := GetWordType(s.Text())
		if t != Unknown {
			return t
		}
	}
	return Unknown

}

// JSON representation of Words.
func (w *Words) JSON() string {
	j, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error())
	}
	return string(j)
}

// List representation of Words.
func (w *Words) List() string {
	var buffer bytes.Buffer

	for _, word := range *w {
		buffer.WriteString(fmt.Sprintf("%s\n", word.List()))
	}

	return buffer.String()
}

// ParseTemplate returns a parsed template based on Words' fields.
func (w *Words) ParseTemplate(tpl string) ([]byte, error) {
	var buf bytes.Buffer

	tmpl, err := template.New("").Parse(tpl)
	if err != nil {
		return []byte{}, err
	}

	err = tmpl.Execute(&buf, w)
	return buf.Bytes(), err
}
