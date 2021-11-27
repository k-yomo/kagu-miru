package pointerconv

func StringToPointer(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}
