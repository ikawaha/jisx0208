package kanji

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// IsRegularUse returns true if the rune r is in the regular-use list (常用漢字表).
func IsRegularUse(r rune) bool {
	return unicode.Is(RegularUseRangeTable, r)
}

// IsNotRegularUseHan returns true if the rune r is in the Han and is not in the regular use list (常用漢字表).
func IsNotRegularUseHan(r rune) bool {
	return unicode.Is(unicode.Han, r) && !unicode.Is(RegularUseRangeTable, r)
}

// Option represents an option for the discriminator.
type Option func(d *RegularUseHanDiscriminator)

// Allow is a discriminator option to set allow characters.
func Allow(r ...rune) Option {
	return func(d *RegularUseHanDiscriminator) {
		d.allow = append(d.allow, r...)
	}
}

// Disallow is a discriminator option to set disallow characters.
func Disallow(r ...rune) Option {
	return func(d *RegularUseHanDiscriminator) {
		d.disallow = append(d.disallow, r...)
	}
}

// RegularUseHanDiscriminator determines if a character is in regular-use kanji or allowed/disallowed character.
type RegularUseHanDiscriminator struct {
	allow    []rune
	disallow []rune
}

// NewRegularHanDiscriminator returns a regular Han character discriminator.
func NewRegularHanDiscriminator(options ...Option) *RegularUseHanDiscriminator {
	var ret RegularUseHanDiscriminator
	for _, option := range options {
		option(&ret)
	}
	return &ret
}

// IsNotRegularUseHan returns true if the rune r is in allowed characters, else if return false r is in disallowed characters,
// otherwise whether r is in JIS X0208 or not.
func (d *RegularUseHanDiscriminator) IsNotRegularUseHan(r rune) bool {
	for _, v := range d.allow {
		if v == r {
			return false
		}
	}
	for _, v := range d.disallow {
		if v == r {
			return true
		}
	}
	return IsNotRegularUseHan(r)
}

// ReplaceNotRegularUseHanAll returns a copy of the string s with each run of not in regular-use Hans
// replaced by the replacement string, which may be empty.
func (d *RegularUseHanDiscriminator) ReplaceNotRegularUseHanAll(s, replacement string) string {
	return replace(s, replacement, func(r rune) bool { return !d.IsNotRegularUseHan(r) })
}

func replace(s, replacement string, is func(rune) bool) string {
	var b strings.Builder

	for i, c := range s {
		if c != utf8.RuneError && is(c) {
			continue
		}
		if !is(c) {
			b.Grow(len(s) + len(replacement))
			b.WriteString(s[:i])
			s = s[i:]
			break
		}
		_, wid := utf8.DecodeRuneInString(s[i:])
		if wid == 1 {
			b.Grow(len(s) + len(replacement))
			b.WriteString(s[:i])
			s = s[i:]
			break
		}
	}

	// Fast path for unchanged input
	if b.Cap() == 0 { // didn't call b.Grow above
		return s
	}

	invalid := false // previous byte was from an invalid UTF-8 sequence
	for i := 0; i < len(s); {
		c := s[i]
		if c < utf8.RuneSelf {
			i++
			invalid = false
			b.WriteByte(c)
			continue
		}
		r, wid := utf8.DecodeRuneInString(s[i:])
		if wid == 1 {
			i++
			if !invalid {
				invalid = true
				b.WriteString(replacement)
			}
			continue
		}
		invalid = false
		if is(r) {
			b.WriteString(s[i : i+wid])
		} else {
			b.WriteString(replacement)
		}
		i += wid
	}
	return b.String()
}
