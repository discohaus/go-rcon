package rcon

import (
	"errors"
	"unicode"
)

type CharSet int

const (
	CharSetASCII CharSet = iota + 1
	CharSetLatin_1
	CharSetUTF8
)

var (
	ErrNonASCII       = errors.New("payload contains non-ASCII characters")
	ErrNonLatin       = errors.New("payload contains non-ISO-8859-1 characters")
	ErrNonUTF8        = errors.New("payload contains non-UTF-8 characters")
	ErrInvalidCharSet = errors.New("invalid character set")
)

func (cs CharSet) ValidateByte(b byte) error {
	switch cs {
	case CharSetASCII:
		if b > unicode.MaxASCII {
			return ErrNonASCII
		}
	case CharSetLatin_1:
		if b > unicode.MaxLatin1 {
			return ErrNonLatin
		}
	case CharSetUTF8:
		return nil
	default:
		return ErrInvalidCharSet
	}
	return nil
}
