package chirouter

import "strings"

type valueSanitizer struct {
	replacer *strings.Replacer
}

func newSanitizer() *valueSanitizer {
	forbiddenChars := "#;%$\"`'"
	replacePairs := []string{}
	for _, char := range forbiddenChars {
		replacePairs = append(replacePairs, string(char), "")
	}

	r := strings.NewReplacer(replacePairs...)
	return &valueSanitizer{
		replacer: r,
	}
}

func (s *valueSanitizer) Sanitize(value string) string {
	return s.replacer.Replace(value)
}
