package xesquery

import "fmt"

func Boost(field string, boost float64) string {
	return fmt.Sprintf("%s^%f", field, boost)
}
