// Package slug turns idea names into kebab-case path segments.
// PHASE-01 fold rules: DEC-014 (not Unicode NFKD).
package slug

import (
	"errors"
	"strings"
	"unicode"
)

const maxLen = 80

// ErrEmpty is returned when slugify yields nothing usable.
var ErrEmpty = errors.New("slug: empty")

// ErrReserved is returned for "." and "..".
var ErrReserved = errors.New("slug: reserved")

// latinFold is the fixed PHASE-01 compatibility map (DEC-014).
// Do not grow this map this phase; unlisted runes are dropped.
var latinFold = map[rune]string{
	'À': "a", 'Á': "a", 'Â': "a", 'Ã': "a", 'Ä': "a", 'Å': "a",
	'à': "a", 'á': "a", 'â': "a", 'ã': "a", 'ä': "a", 'å': "a",
	'È': "e", 'É': "e", 'Ê': "e", 'Ë': "e",
	'è': "e", 'é': "e", 'ê': "e", 'ë': "e",
	'Ì': "i", 'Í': "i", 'Î': "i", 'Ï': "i",
	'ì': "i", 'í': "i", 'î': "i", 'ï': "i",
	'Ò': "o", 'Ó': "o", 'Ô': "o", 'Õ': "o", 'Ö': "o",
	'ò': "o", 'ó': "o", 'ô': "o", 'õ': "o", 'ö': "o",
	'Ù': "u", 'Ú': "u", 'Û': "u", 'Ü': "u",
	'ù': "u", 'ú': "u", 'û': "u", 'ü': "u",
	'Ý': "y", 'ý': "y", 'ÿ': "y",
	'Ñ': "n", 'ñ': "n",
	'Ç': "c", 'ç': "c",
	'Š': "s", 'š': "s",
	'Ž': "z", 'ž': "z",
	'Æ': "ae", 'æ': "ae",
	'Œ': "oe", 'œ': "oe",
}

// Slugify converts name to a kebab-case slug per DEC-014.
// latinFold + ASCII [a-zA-Z0-9], map space/_ to -, drop other runes
// (and combining marks), collapse --, trim -, lower, max 80.
// Empty / "." / ".." refuse. Unlisted letters are dropped, not folded.
func Slugify(name string) (string, error) {
	var b strings.Builder
	b.Grow(len(name))
	prevDash := false
	for _, r := range name {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		if repl, ok := latinFold[r]; ok {
			for _, c := range repl {
				writeAlnum(&b, &prevDash, c)
			}
			continue
		}
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			writeAlnum(&b, &prevDash, unicode.ToLower(r))
		case r == ' ' || r == '_':
			if b.Len() > 0 && !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}

	out := strings.Trim(b.String(), "-")
	if len(out) > maxLen {
		out = strings.Trim(out[:maxLen], "-")
	}
	if out == "" {
		return "", ErrEmpty
	}
	if out == "." || out == ".." {
		return "", ErrReserved
	}
	return out, nil
}

func writeAlnum(b *strings.Builder, prevDash *bool, r rune) {
	b.WriteRune(r)
	*prevDash = false
}
