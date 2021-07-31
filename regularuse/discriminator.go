package regularuse

import (
	"unicode"
)

// Is returns true if the rune r is in the regular use list (常用漢字表).
func Is(r rune) bool {
	return unicode.Is(RangeTable, r)
}

// IsNotRegularIdeographic returns true if the rune r is in the ideographic and is not in the regular use list (常用漢字表).
func IsNotRegularIdeographic(r rune) bool {
	return unicode.Is(unicode.Ideographic, r) && !unicode.Is(RangeTable, r)
}
