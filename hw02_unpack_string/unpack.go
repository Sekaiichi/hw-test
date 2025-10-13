package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var (
		b          strings.Builder
		lastRune   rune
		hasLast    bool
		escapeNext bool
	)

	emitRepeat := func(r rune, n int) {
		if n <= 0 {
			return
		}
		b.WriteString(strings.Repeat(string(r), n))
	}

	for _, r := range s {
		if escapeNext {
			if !(unicode.IsDigit(r) || r == '\\') {
				return "", ErrInvalidString
			}
			if hasLast {
				emitRepeat(lastRune, 1)
			}
			lastRune = r
			hasLast = true
			escapeNext = false
			continue
		}

		switch {
		case r == '\\':
			escapeNext = true

		case unicode.IsDigit(r):
			if !hasLast {
				return "", ErrInvalidString
			}
			repeat := int(r - '0')
			emitRepeat(lastRune, repeat)
			hasLast = false

		default:
			if hasLast {
				emitRepeat(lastRune, 1)
			}
			lastRune = r
			hasLast = true
		}
	}

	// Конец строки
	if escapeNext {
		// Висячий слэш
		return "", ErrInvalidString
	}
	if hasLast {
		emitRepeat(lastRune, 1)
	}

	return b.String(), nil
}
