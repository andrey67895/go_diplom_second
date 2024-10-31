package helpers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeJWTError(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeJWT(tt.args.tokenString)
			assert.Equal(t, "", got)
			assert.True(t, strings.Contains(err.Error(), "ошибка разбора: signature is invalid"))
		})
	}
}

func TestGenerateJWTAndCheck(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive test #1",
			args: args{username: "TEST_USER"},
			want: "TEST_USER",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateJWTAndCheck(tt.args.username)
			assert.NoError(t, err)
			name, err := DecodeJWT(token)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, name)
		})
	}
}
