package kanji

import (
	"unicode"
)

// IsForPersonalNames returns true if the rune r is kanji designated for personal names.
func IsForPersonalNames(r rune) bool {
	return IsRegularUse(r) || unicode.Is(DesignatedForPersonalNamesRangeTable, r)
}
