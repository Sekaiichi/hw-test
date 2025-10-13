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
			// Можно экранировать только цифру или '\' — иначе ошибка.
			if !(unicode.IsDigit(r) || r == '\\') {
				return "", ErrInvalidString
			}
			// Новая литеральная руна (как обычный символ).
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
			// Цифра – множитель к предыдущей руне.
			if !hasLast {
				// Нет предыдущей руны -> строка некорректна (начинается с цифры или идут подряд цифры).
				return "", ErrInvalidString
			}
			repeat := int(r - '0')
			emitRepeat(lastRune, repeat)
			// После применения множителя предыдущая руна "израсходована".
			hasLast = false

		default:
			// Пришла новая обычная руна.
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
