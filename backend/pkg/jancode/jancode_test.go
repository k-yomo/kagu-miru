package jancode

import "testing"

func TestExtractJANCode(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{s: "prefix4547366419429suffix"},
			want: "4547366419429",
		},
		{
			name: "no match",
			args: args{s: "abc454736641942"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractJANCode(tt.args.s); got != tt.want {
				t.Errorf("ExtractJANCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
