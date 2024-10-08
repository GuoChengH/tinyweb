package tinyweb

import (
	"strings"
	"unicode"
)

// SubStringLast 将字符串str中最后一次出现substr的后面内容返回
// 例如:SubStringLast("/api/user/1","/api") -> "/user/1"
func SubStringLast(str string, substr string) string {
	index := strings.Index(str, substr)
	if index < 0 {
		return ""
	}
	return str[index+len(substr):]
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
