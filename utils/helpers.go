package utils

import "strings"

func TrimAndCapitalize(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func GenerateBankAbbreviation(name string) string {
	name = strings.TrimSpace(name)
	if len(name) < 4 {
		return strings.ToLower(name)
	}
	return strings.ToLower(name[:2] + name[len(name)-2:])
}
