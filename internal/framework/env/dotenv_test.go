package env

import "testing"

func TestInitAllowedOrigins(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectError bool
	}{
		{
			name:        "Valid origins",
			input:       "https://localhost:3000;http://example.com",
			expected:    []string{"https://localhost:3000", "http://example.com"},
			expectError: false,
		},
		{
			name:        "Invalid URL",
			input:       "https://localhost:3000;invalid-url",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Empty input",
			input:       "",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Unsupported scheme",
			input:       "ftp://localhost:3000;https://example.com",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Whitespace trimming",
			input:       " https://localhost:3000 ; https://example.com ",
			expected:    []string{"https://localhost:3000", "https://example.com"},
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := initAllowedOrigins(test.input)

			if test.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}

				if len(result) != len(test.expected) {
					t.Errorf("expected %v, got %v", test.expected, result)
				}
				for i := range result {
					if result[i] != test.expected[i] {
						t.Errorf("expected %v, got %v", test.expected, result)
					}
				}
			}
		})
	}
}
