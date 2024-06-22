package utils

import (
	"strings"
)

// ValidateRole is a function to validate role and return correct role format.
// It will return Capitalize role if the role is valid, otherwise it will return false
func ValidateRole(role string) (string, bool) {
	switch strings.ToLower(role) {
	case "admin":
		return "Admin", true
	case "editor":
		return "Editor", true
	case "user":
		return "User", true
	default:
		return "", false
	}
}
