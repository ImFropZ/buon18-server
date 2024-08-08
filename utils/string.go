package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func PascalToSnake(s string) string {
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	snake := re.ReplaceAllString(s, `${1}_${2}`)
	return strings.ToLower(snake)
}

func ContainsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func StrToInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return result
}
