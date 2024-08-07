package utils

import "strings"

func SerializeGender(g *string) string {
	if g == nil {
		return "u"
	}
	switch strings.ToLower(*g) {
	case "m":
		return "m"
	case "f":
		return "f"
	default:
		return "u"
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
