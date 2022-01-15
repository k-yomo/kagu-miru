package xspanner

import (
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/google/go-cmp/cmp"
)

func Test_getColumnNames(t *testing.T) {
	t.Parallel()

	type test struct {
		Column1 string           `spanner:"column1"`
		Column2 int              `spanner:"column2"`
		Column3 spanner.NullTime `spanner:"column3"`
	}

	tests := []struct {
		name string
		dto  interface{}
		want []string
	}{
		{
			name: "returns all columns",
			dto:  test{},
			want: []string{
				"column1",
				"column2",
				"column3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getColumnNames(tt.dto)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ItemAllColumns(), (-want +got): %s", diff)
			}
		})
	}
}
