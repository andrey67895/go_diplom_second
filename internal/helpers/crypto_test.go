package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeHashSha256(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test #1",
			args: args{value: "TEST_ENCODE"},
			want: "6d222efbcfc554d3acf5446947765438be89d5d395bf3b0c03f926ddb5d3184a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, EncodeHashSha256(tt.args.value), "EncodeHashSha256(%v)", tt.args.value)
		})
	}
}

func TestEncodeHashSha512(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test #1",
			args: args{value: "TEST_ENCODE"},
			want: "efdb8a54032f88502ae0c01356a68620e43293ddc7f6cf8dadbaf9ae5a77bcdd722b4a90c042daddbc18b7de24af91b4d64434650e9dd5c92943db5986d9f588",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, EncodeHashSha512(tt.args.value), "EncodeHashSha512(%v)", tt.args.value)
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	type args struct {
		input string
		key   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test #1",
			args: args{input: "TEST_ENCODE", key: "KEY"},
			want: "\x1f\x00\n\x1f\x1a\x1c\x05\x06\x16\x0f\x00",
		},
		{
			name: "positive test #2",
			args: args{input: "\u001F\u0000\n\u001F\u001A\u001C\u0005\u0006\u0016\u000F\u0000", key: "KEY"},
			want: "TEST_ENCODE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, EncryptDecrypt(tt.args.input, tt.args.key), "EncryptDecrypt(%v, %v)", tt.args.input, tt.args.key)
		})
	}
}
