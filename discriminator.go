package jisx0208

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Is returns true if the rune r is in JIS X 0208.
func Is(r rune) bool {
	return unicode.Is(RangeTable, r)
}

// IsLevel1 returns true if the rune r is in JIS X 0208 Level1 (第一水準).
func IsLevel1(r rune) bool {
	return unicode.Is(Level1RangeTable, r)
}

// IsLevel2 returns true if the rune r is in JIS X 0208 Level2 (第二水準).
func IsLevel2(r rune) bool {
	return unicode.Is(Level2RangeTable, r)
}

// ToValid returns a copy of the string s with each run of invalid JIS X 0208
// runes replaced by the replacement string, which may be empty.
func ToValid(s, replacement string) string {
	return toValid(s, replacement, Is)
}

// Option represents an option for the discriminator.
type Option func(d *Discriminator)

// Allow is a discriminator option to set allow characters.
func Allow(r ...rune) Option {
	return func(d *Discriminator) {
		d.allow = append(d.allow, r...)
	}
}

// Disallow is a discriminator option to set disallow characters.
func Disallow(r ...rune) Option {
	return func(d *Discriminator) {
		d.disallow = append(d.disallow, r...)
	}
}

// Discriminator determines if a character is in JISX0208 or allowed/disallowed character.
type Discriminator struct {
	allow    []rune
	disallow []rune
}

// NewDiscriminator returns a character discriminator.
func NewDiscriminator(options ...Option) *Discriminator {
	var ret Discriminator
	for _, option := range options {
		option(&ret)
	}
	return &ret
}

// Is returns true if the rune r is in allowed characters, else if return false r is in disallowed characters,
// otherwise whether r is in JIS X0208 or not.
func (d *Discriminator) Is(r rune) bool {
	for _, v := range d.allow {
		if v == r {
			return true
		}
	}
	for _, v := range d.disallow {
		if v == r {
			return false
		}
	}
	return Is(r)
}

// ToValid returns a copy of the string s with each run of invalid runes
// replaced by the replacement string, which may be empty.
func (d *Discriminator) ToValid(s, replacement string) string {
	return toValid(s, replacement, d.Is)
}

func toValid(s, replacement string, is func(rune) bool) string {
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
