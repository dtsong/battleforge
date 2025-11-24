package main

import (
	"os"
	"strings"
	"testing"
)

func TestGetAddr(t *testing.T) {
	tests := []struct {
		name     string
		envVar   string
		envValue string
		expected string
	}{
		{
			name:     "default address",
			envVar:   "SERVER_PORT",
			envValue: "",
			expected: ":8080",
		},
		{
			name:     "custom port from env",
			envVar:   "SERVER_PORT",
			envValue: "9000",
			expected: ":9000",
		},
		{
			name:     "port 3000",
			envVar:   "SERVER_PORT",
			envValue: "3000",
			expected: ":3000",
		},
		{
			name:     "port 80",
			envVar:   "SERVER_PORT",
			envValue: "80",
			expected: ":80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.envValue != "" {
				_ = os.Setenv(tt.envVar, tt.envValue)
				defer func() { _ = os.Unsetenv(tt.envVar) }()
			} else {
				_ = os.Unsetenv(tt.envVar)
			}

			result := getAddr()

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetDBConnString(t *testing.T) {
	tests := []struct {
		name            string
		setupEnv        map[string]string
		expectedContent []string
	}{
		{
			name:     "default values",
			setupEnv: map[string]string{},
			expectedContent: []string{
				"postgres://",
				"localhost",
				"5432",
				"vgccorner",
				"sslmode=disable",
			},
		},
		{
			name: "custom host",
			setupEnv: map[string]string{
				"DB_HOST": "db.example.com",
			},
			expectedContent: []string{
				"db.example.com",
			},
		},
		{
			name: "custom port",
			setupEnv: map[string]string{
				"DB_PORT": "5433",
			},
			expectedContent: []string{
				":5433",
			},
		},
		{
			name: "custom database name",
			setupEnv: map[string]string{
				"DB_NAME": "custom_db",
			},
			expectedContent: []string{
				"/custom_db",
			},
		},
		{
			name: "all custom values",
			setupEnv: map[string]string{
				"DB_HOST":     "production.db",
				"DB_PORT":     "5432",
				"DB_USER":     "produser",
				"DB_PASSWORD": "prodpass",
				"DB_NAME":     "proddb",
				"DB_SSL_MODE": "require",
			},
			expectedContent: []string{
				"production.db",
				"produser",
				"prodpass",
				"proddb",
				"sslmode=require",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env vars
			envKeys := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE"}
			savedEnv := make(map[string]string)

			for _, key := range envKeys {
				if val, exists := os.LookupEnv(key); exists {
					savedEnv[key] = val
				}
			}

			// Clear env vars
			for _, key := range envKeys {
				_ = os.Unsetenv(key)
			}

			// Set test env vars
			for key, val := range tt.setupEnv {
				_ = os.Setenv(key, val)
			}

			// Restore original env vars
			defer func() {
				for _, key := range envKeys {
					_ = os.Unsetenv(key)
				}
				for key, val := range savedEnv {
					_ = os.Setenv(key, val)
				}
			}()

			result := getDBConnString()

			// Check that expected content is in the result
			for _, expected := range tt.expectedContent {
				if !strings.Contains(result, expected) {
					t.Errorf("expected connection string to contain %q, got %q", expected, result)
				}
			}

			// Should always start with postgres://
			if !strings.Contains(result, "postgres://") {
				t.Errorf("expected postgres:// scheme, got %q", result)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		envValue   string
		defaultVal string
		expected   string
		setEnv     bool
	}{
		{
			name:       "environment variable set",
			key:        "TEST_VAR",
			envValue:   "from_env",
			defaultVal: "default",
			expected:   "from_env",
			setEnv:     true,
		},
		{
			name:       "environment variable not set",
			key:        "TEST_VAR_NOT_SET",
			defaultVal: "default_value",
			expected:   "default_value",
			setEnv:     false,
		},
		{
			name:       "empty string from env returns default",
			key:        "EMPTY_VAR",
			envValue:   "",
			defaultVal: "fallback",
			expected:   "fallback",
			setEnv:     true,
		},
		{
			name:       "empty default value",
			key:        "TEST_KEY",
			envValue:   "value",
			defaultVal: "",
			expected:   "value",
			setEnv:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				_ = os.Setenv(tt.key, tt.envValue)
				defer func() { _ = os.Unsetenv(tt.key) }()
			} else {
				_ = os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultVal)

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
