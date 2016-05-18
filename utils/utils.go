package utils

import (
	"regexp"
	"unicode"
)

func Replace(str string, tag string, replacer string) string {
	r := regexp.MustCompile(tag)
	return r.ReplaceAllString(str, replacer)
}

func UpcaseInitial(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i + 1:]
	}
	return ""
}

func Int64InSlice(a int64, list []int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}