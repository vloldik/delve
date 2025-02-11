package quals

import (
	"strings"

	"github.com/vloldik/delve/v2/internal/defaultval"
	"github.com/vloldik/delve/v2/pkg/interfaces"
)

const DefaultDelimiter = '.' // Qdelimiter is used to separate nested keys in qualified paths

type compiledQual struct {
	parts     []string
	len       uint8
	index     uint8
	delimiter rune
}

func (c *compiledQual) Copy() interfaces.IQual {
	return &compiledQual{
		// No need to copy list, it's read-only
		parts:     c.parts,
		len:       c.len,
		index:     c.index,
		delimiter: c.delimiter,
	}
}

func (c *compiledQual) Next() (string, bool) {
	if c.index >= c.len {
		return "", false
	}
	part := c.parts[c.index]
	hasNext := c.index < c.len-1
	c.index++
	return part, hasNext
}

func (c *compiledQual) Reset() {
	c.index = 0
}

func (c *compiledQual) String() string {
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

// Creates a compiled qual, which is more efficient for reuse, but has a higher creation cost than string qual.
func CQ(qual string, _delimiter ...rune) *compiledQual {
	delimiter := defaultval.WithDefaultVal(DefaultDelimiter, _delimiter)
	if delimiter == '\\' {
		panic(`delimiter can not be a "\"`)
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

	return &compiledQual{
		parts:     parts,
		len:       uint8(len(parts)),
		delimiter: delimiter,
		index:     0,
	}
}
