package flexmap

import "strings"

// FlexQual is compiled path
type CompiledQual []string

func (c CompiledQual) String() string {
	escaped := []string{}
	for _, part := range c {
		escaped = append(escaped, strings.ReplaceAll(part, string(DefaultDelimiter), "\\"+string(DefaultDelimiter)))
	}
	return strings.Join(escaped, string(DefaultDelimiter))
}

func CompileQual(qual string, _delimiter ...rune) CompiledQual {
	delimiter := DefaultDelimiter
	parts := make(CompiledQual, 0)
	if len(_delimiter) > 0 {
		delimiter = _delimiter[0]
	}
	for part, rest := getNextQualPart(delimiter, qual); ; part, rest = getNextQualPart(delimiter, rest) {
		parts = append(parts, part)
		if rest == "" {
			break
		}
	}
	return parts
}

// returns first part of qual and its end position
func getNextQualPart(delimiter rune, qual string) (part string, rest string) {
	var isEscaped bool
	var qualPart string
	for i, r := range qual {
		if r == delimiter && !isEscaped {
			return qualPart, qual[i+1:]
		}
		if r == '\\' && !isEscaped {
			isEscaped = true
			continue
		} else {
			isEscaped = false
		}
		qualPart += string(r)
	}
	return qualPart, ""
}
