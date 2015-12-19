package models

import (
	"strings"
	"unicode/utf8"
)

var (
	reservedNames    = []string{"debug", "api", "erp", "catalog", "help", "template", "admin", "user", "create", "new", "edit", "list", "delete", "system"}
	reservedPatterns = []string{"*.git", "*.zip", "*.exe"}
)

func IsUsableName(name string) error {
	name = strings.TrimSpace(strings.ToLower(name))
	if utf8.RuneCountInString(name) == 0 {
		return ErrNameEmpty
	}

	for i := range reservedNames {
		if name == reservedNames[i] {
			return ErrNameReserved{name}
		}
	}

	for _, pat := range reservedPatterns {
		if pat[0] == '*' && strings.HasSuffix(name, pat[1:]) ||
			(pat[len(pat)-1] == '*' && strings.HasPrefix(name, pat[:len(pat)-1])) {
			return ErrNamePatternNotAllowed{pat}
		}
	}

	return nil
}
