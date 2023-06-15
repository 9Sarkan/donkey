package donkey

import "strings"

type Reader interface {
	SetPrefix(string)
	SetReplacer(replacer *strings.Replacer)
	Read() error
	Validate() error
}
