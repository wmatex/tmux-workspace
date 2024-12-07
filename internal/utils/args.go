package utils

import (
	"strings"
)

func SplitArgs(cmd string) []string {
	var args []string
	var sb strings.Builder
	var quote rune

	for _, c := range cmd {
		if c == quote {
			quote = 0
		} else if quote == 0 && (c == '\'' || c == '"') {
			quote = c
		} else if quote == 0 && c == ' ' {
			if sb.Len() > 0 {
				args = append(args, sb.String())
				sb.Reset()
			}
		} else {
			sb.WriteRune(c)
		}
	}

	if sb.Len() > 0 {
		args = append(args, sb.String())
	}

	return args
}
