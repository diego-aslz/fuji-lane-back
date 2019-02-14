package fujilane

// IsBlankStr returns true if the given string pointer is nil or empty
func IsBlankStr(str *string) bool {
	return str == nil || *str == ""
}

// IsBlankUint returns true if the given uint pointer is nil or 0
func IsBlankUint(i *uint) bool {
	return i == nil || *i == 0
}
