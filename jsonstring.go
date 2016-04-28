package quicktemplate

// appendJSONString is a synonym to strconv.AppendQuote, but works 3x faster.
func appendJSONString(dst []byte, s string) []byte {
	j := 0
	out := ""
	dst = append(dst, `"`...)
	for i, n := 0, len(s); i < n; i++ {
		switch s[i] {
		case '"':
			out = `\"`
		case '\\':
			out = `\\`
		case '\n':
			out = `\n`
		case '\r':
			out = `\r`
		case '\t':
			out = `\t`
		case '\f':
			out = `\u000c`
		case '\b':
			out = `\u0008`
		case '<':
			out = `\u003c`
		case '\'':
			out = `\u0027`
		case 0:
			out = `\u0000`
		}
		if len(out) > 0 {
			dst = append(dst, s[j:i]...)
			dst = append(dst, out...)
			j = i + 1
			out = ""
		}
	}
	dst = append(dst, s[j:]...)
	return append(dst, `"`...)
}
