package sanitize

import "strings"

func IsMalicious(s string) bool {
	s = strings.ToLower(s)

	// Basit XSS ve SQL Injection kontrolleri
	if strings.Contains(s, "<script") || strings.Contains(s, "</script") {
		return true
	}
	if strings.Contains(s, "drop table") || strings.Contains(s, "delete from") {
		return true
	}

	return false
}
