package config

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetTemplatedEnv(t *testing.T) {
	t.Setenv("SOME_ENV_VAR", "tutu")
	result := getTemplatedEnv("{{SOME_ENV_VAR}}")
	if result != "tutu" {
		t.Error("Failed env var substitution")
	}
}

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name       string
		config     Config
		headersArg string
		result     http.Header
		error      error
	}{
		{"Empty headers", Config{}, "", http.Header{"User-Agent": {"Gmax76/urlcheck"}}, nil},
		{"Authorization header", Config{}, "Authorization: Basic YXplcnR5Cg==", http.Header{"Authorization": {"Basic YXplcnR5Cg=="}, "User-Agent": {"Gmax76/urlcheck"}}, nil},
		{"Several headers", Config{}, "Authorization: Basic YXplcnR5Cg==, Some-Header: itsvalue", http.Header{"Authorization": {"Basic YXplcnR5Cg=="}, "Some-Header": {"itsvalue"}, "User-Agent": {"Gmax76/urlcheck"}}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers, err := tt.config.parseHeaders(tt.headersArg)
			if !cmp.Equal(headers, tt.result) {
				t.Errorf("Result returned %v, wanted %v", headers, tt.result)
			}
			if err != tt.error {
				t.Errorf("Error should be %v, is %v", tt.error, err)
			}
		})
	}
}
