package delve

import (
	"strings"
)

const DefaultDelimiter = '.' // Qdelimiter is used to separate nested keys in qualified paths

type CompiledQual struct {
	parts     []string
	len       uint8
	index     uint8
	delimiter rune
}

func (c *CompiledQual) Copy() *CompiledQual {
	return &CompiledQual{
		parts:     c.parts,
		len:       c.len,
		index:     0,
		delimiter: c.delimiter,
	}
}

func (c *CompiledQual) Next() (string, bool) {
	if c.index >= c.len {
		return "", false
	}
	part := c.parts[c.index]
	hasNext := c.index < c.len-1
	c.index++
	return part, hasNext
}

func (c *CompiledQual) Reset() {
	c.index = 0
}

func (c *CompiledQual) String() string {
	if len(c.parts) == 0 {
		return ""
	}

	var builder strings.Builder
	totalLen := 0
	delimStr := string(c.delimiter)

	for i, part := range c.parts {
		if i == 0 {
			continue
		}
		totalLen += 1 + len(part) + strings.Count(part, delimStr)
	}
	builder.Grow(totalLen)

	// Build the escaped string
	for i, part := range c.parts {
		if i > 0 {
			builder.WriteByte(byte(DefaultDelimiter))
		}
		for _, r := range part {
			if r == DefaultDelimiter || r == '\\' {
				builder.WriteRune('\\')
			}
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func CQ(qual string, _delimiter ...rune) *CompiledQual {
	delimiter := DefaultDelimiter
	if len(_delimiter) > 0 {
		delimiter = _delimiter[0]
	}
	if delimiter == '\\' {
		panic("Delimiter cannot be a backslash")
	}

	expectedParts := strings.Count(qual, string(delimiter)) + 1
	parts := make([]string, 0, expectedParts)

	var currentPart strings.Builder
	currentPart.Grow(16) // Preallocate a small buffer to minimize reallocations
	var escapeNext bool

	for _, r := range qual {
		if escapeNext {
			currentPart.WriteRune(r)
			escapeNext = false
			continue
		}

		switch r {
		case '\\':
			escapeNext = true
		case delimiter:
			parts = append(parts, currentPart.String())
			currentPart.Reset()
		default:
			currentPart.WriteRune(r)
		}
	}
	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	if len(parts) > 254 {
		panic("qual len is too large!")
	}

	return &CompiledQual{
		parts:     parts,
		len:       uint8(len(parts)),
		delimiter: delimiter,
		index:     0,
	}
}
