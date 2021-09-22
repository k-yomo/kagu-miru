package jancode

import "regexp"

// ExtractJANCode 13 digits JAN code
func ExtractJANCode(s string) string {
	janCode := regexp.MustCompile("[0-9]{13}").FindString(s)
	return janCode
}
