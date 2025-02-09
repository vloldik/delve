package delve

type stringQual struct {
	_initQual string
	qual      string
	delimiter rune
}

func (sq *stringQual) Copy() IQual {
	return &stringQual{
		_initQual: sq._initQual,
		qual:      sq.qual,
		delimiter: sq.delimiter,
	}
}

func (sq *stringQual) Next() (string, bool) {
	return sq.getNextPart(), sq.qual != ""
}

func (sq *stringQual) getDelemiterIndex() int {
	var escapeNext bool
	removedCharCount := 0
	for i, r := range sq.qual {
		if escapeNext && (r == sq.delimiter || r == '\\') {
			escapeNext = false
			continue
		}
		if r == sq.delimiter {
			return i - removedCharCount
		}
		if r == '\\' {
			sq.qual = sq.qual[0:i-removedCharCount] + sq.qual[i-removedCharCount+1:]
			removedCharCount++
			escapeNext = true
			continue
		}
		escapeNext = false
	}
	return -1
}

func (sq *stringQual) getNextPart() string {
	i := sq.getDelemiterIndex()
	if i == -1 {
		part := sq.qual
		sq.qual = ""
		return part
	}
	part := sq.qual[:i]
	if len(sq.qual) > i+1 {
		sq.qual = sq.qual[i+1:]
	} else {
		sq.qual = ""
	}
	return part
}

func (sq *stringQual) Reset() {
	sq.qual = sq._initQual
}

// Creates an uncompiled qualifier, which is cheaper to create than a compiled one, but more expensive to reuse and depends on the length of the string.
func Q(qual string, _delimiter ...rune) *stringQual {
	delimiter := DefaultDelimiter
	if len(_delimiter) > 0 {
		delimiter = _delimiter[0]
	}
	if delimiter == '\\' {
		panic(`delimiter can not be a "\"`)
	}
	return &stringQual{
		delimiter: delimiter,
		_initQual: qual,
		qual:      qual,
	}
}
