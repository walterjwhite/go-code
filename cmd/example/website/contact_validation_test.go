package main

import (
	"strings"
	"testing"
)

func TestValidateEmailAddress(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{name: "valid simple address", email: "user@example.com", want: true},
		{name: "invalid format", email: "not-an-email", want: false},
		{name: "display name rejected", email: "User <user@example.com>", want: false},
		{name: "whitespace mismatch rejected", email: " user@example.com ", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateEmailAddress(tt.email)
			if got != tt.want {
				t.Fatalf("validateEmailAddress(%q)=%v want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestValidateContactRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *ContactRequest
		wantError string
	}{
		{
			name: "valid request with trimming",
			request: &ContactRequest{
				Name:    " Alice ",
				Email:   "alice@example.com",
				Subject: " Hello ",
				Message: " Hi there ",
			},
		},
		{
			name:      "nil request",
			request:   nil,
			wantError: "invalid request payload",
		},
		{
			name: "missing required field",
			request: &ContactRequest{
				Name:    "",
				Email:   "alice@example.com",
				Subject: "Hello",
				Message: "Hi",
			},
			wantError: "all fields are required",
		},
		{
			name: "message exceeds max length",
			request: &ContactRequest{
				Name:    "Alice",
				Email:   "alice@example.com",
				Subject: "Hello",
				Message: strings.Repeat("a", maxMessageLength+1),
			},
			wantError: "message exceeds maximum length of 5000 characters",
		},
		{
			name: "invalid email",
			request: &ContactRequest{
				Name:    "Alice",
				Email:   "invalid",
				Subject: "Hello",
				Message: "Hi",
			},
			wantError: "invalid email address",
		},
		{
			name: "header injection in subject",
			request: &ContactRequest{
				Name:    "Alice",
				Email:   "alice@example.com",
				Subject: "Hello\nX-Injected: yes",
				Message: "Hi",
			},
			wantError: "invalid characters in request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateContactRequest(tt.request)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("validateContactRequest() unexpected error: %v", err)
				}
				if tt.request.Name != strings.TrimSpace(tt.request.Name) {
					t.Fatalf("name was not normalized")
				}
				if tt.request.Subject != strings.TrimSpace(tt.request.Subject) {
					t.Fatalf("subject was not normalized")
				}
				if tt.request.Message != strings.TrimSpace(tt.request.Message) {
					t.Fatalf("message was not normalized")
				}
				return
			}

			if err == nil {
				t.Fatalf("validateContactRequest() expected error %q", tt.wantError)
			}
			if err.Error() != tt.wantError {
				t.Fatalf("validateContactRequest() error=%q want %q", err.Error(), tt.wantError)
			}
		})
	}
}

func TestHasCRLF(t *testing.T) {
	if hasCRLF("hello") {
		t.Fatalf("expected false for normal text")
	}
	if !hasCRLF("hello\nworld") {
		t.Fatalf("expected true for newline")
	}
	if !hasCRLF("hello\rworld") {
		t.Fatalf("expected true for carriage return")
	}
}
