package utils

func SanitizeInt(val int, minVal int, maxVal int, defaultVal int) int {
    if val < minVal {
        return defaultVal
    }

    if val > maxVal {
        return defaultVal
    }

    return val;
}

func SanitizeString(str string, allowedStr []string, defaultStr string) string {
	// If the string is empty, return the default string
    if (str == "") {
        return defaultStr
    }

    // If the string is in the allowed list, return it
    for _, a := range allowedStr {
        if str == a {
            return str
        }
    }

    // Otherwise, return the default string as fallback
	return defaultStr
}
