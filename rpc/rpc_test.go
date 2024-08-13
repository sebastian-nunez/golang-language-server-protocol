package rpc

import (
	"testing"
)

func TestEncodeMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncodeMessage(tc.msg)
			if got != tc.want {
				t.Errorf("EncodeMessage got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDecodeMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		msg               []byte
		wantContent       []byte
		wantMethod        string
		wantContentLength int
	}{
		{
			name:              "simple message",
			msg:               []byte("Content-Length: 17\r\n\r\n{\"Method\":\"post\"}"),
			wantContent:       []byte("{\"Method\":\"post\"}"),
			wantMethod:        "post",
			wantContentLength: 17,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotMethod, gotContent, err := DecodeMessage(tc.msg)
			if err != nil {
				t.Errorf("DecodeMessage got unexpected error: %v", err)
			}

			gotContentLength := len(gotContent)
			if gotContentLength != tc.wantContentLength {
				t.Errorf("DecodeMessage content length got %v, want %v", gotContentLength, tc.wantContentLength)
			}

			if gotMethod != tc.wantMethod {
				t.Errorf("DecodeMessage method got %v, want %v", gotMethod, tc.wantMethod)
			}

			if string(gotContent) != string(tc.wantContent) {
				t.Errorf("DecodeMessage content got %v, want %v", string(gotContent), string(tc.wantContent))
			}
		})
	}
}
