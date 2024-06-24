package utils

import "strings"

func SerializeGender(g string) string {
	switch strings.ToLower(g) {
	case "m":
		return "M"
	case "f":
		return "F"
	default:
		return "U"
	}
}

func DeserializeGender(g string) string {
	switch strings.ToUpper(g) {
	case "M":
		return "male"
	case "F":
		return "female"
	default:
		return "unknown"
	}
}
