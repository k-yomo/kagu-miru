package interfaceconv

func StringArrayToInterfaceArray(strArray []string) []interface{} {
	arr := make([]interface{}, 0, len(strArray))
	for _, str := range strArray {
		arr = append(arr, str)
	}
	return arr
}
