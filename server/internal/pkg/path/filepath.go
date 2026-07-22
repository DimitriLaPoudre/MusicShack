package path

import "strings"

func SanitizeName(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z':
			b.WriteByte(byte(r))
		case r >= '0' && r <= '9':
			b.WriteByte(byte(r))
		case r == '.' || r == '-' || r == '_':
			b.WriteByte(byte(r))
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}
