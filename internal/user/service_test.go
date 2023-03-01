package user

import (
	"testing"
)

func TestUser_IsEmailValid(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "valid email",
			email:    "dilaragorum@gmail.com",
			expected: true,
		},
		{
			name:     "invalid email",
			email:    "dilaragorum",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Email: tt.email}
			if got := u.IsEmailValid(); got != tt.expected {
				t.Errorf("IsEmailValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}
