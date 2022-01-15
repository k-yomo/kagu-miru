package xspanner

import "reflect"

// getColumnNames extract columns from struct
// dto must be struct like Item{} that has `spanner:"column"` tag
func getColumnNames(dto interface{}) []string {
	t := reflect.TypeOf(dto)

	var columns []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columns = append(columns, field.Tag.Get("spanner"))
	}
	return columns

}
