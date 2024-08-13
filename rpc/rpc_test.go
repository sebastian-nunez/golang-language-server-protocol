package rpc

import (
	"testing"
)

func TestEncodeMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		msg  any
		want string
	}{
		{
			name: "simple string",
			msg:  "hello",
			want: "Content-Length: 7\r\n\r\n\"hello\"",
		},
		{
			name: "simple map",
			msg:  map[string]string{"key": "value"},
			want: "Content-Length: 15\r\n\r\n{\"key\":\"value\"}",
		},
		{
			name: "nil input",
			msg:  nil,
			want: "Content-Length: 4\r\n\r\nnull",
		},
		{
			name: "custom struct",
			msg: struct {
				Name string
				Age  int
			}{"Alice", 30},
			want: "Content-Length: 25\r\n\r\n{\"Name\":\"Alice\",\"Age\":30}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeMessage(tt.msg)
			if got != tt.want {
				t.Errorf("EncodeMessage got %v, want %v", got, tt.want)
			}
		})
	}
}
