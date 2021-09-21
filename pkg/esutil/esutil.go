package esutil

import "fmt"

func BoostFieldForMultiMatch(field string, boost int) string {
	return fmt.Sprintf("%s^%d", field, boost)
}
