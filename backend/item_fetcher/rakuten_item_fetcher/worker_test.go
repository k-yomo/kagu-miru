package main

import (
	"reflect"
	"testing"
)

func Test_getDimensionRangeByTagID(t *testing.T) {
	type args struct {
		tagID     int64
		fromTagID int64
	}
	tests := []struct {
		name string
		args args
		want *intRange
	}{
		{
			name: "~ 19cm tag id found",
			args: args{
				tagID:     1,
				fromTagID: 1,
			},
			want: &intRange{
				Gte: 0,
				Lte: func() *int { i := 19; return &i }(),
			},
		},
		{
			name: "200cm ~ tag id found",
			args: args{
				tagID:     20,
				fromTagID: 1,
			},
			want: &intRange{
				Gte: 200,
				Lte: nil,
			},
		},
		{
			name: "tag id not found",
			args: args{
				tagID:     21,
				fromTagID: 1,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDimensionRangeByTagID(tt.args.tagID, tt.args.fromTagID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDimensionRangeByTagID() = %v, want %v", got, tt.want)
			}
		})
	}
}
