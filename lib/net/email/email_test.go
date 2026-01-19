package email

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/emersion/go-message/mail"
)

func TestEmailAccount_String(t *testing.T) {
	tests := []struct {
		name     string
		account  *EmailAccount
		expected string
	}{
		{
			name: "Full EmailAccount",
			account: &EmailAccount{
				Username: "testuser",
				Password: "testpassword",
				Domain:   "example.com",
				EmailAddress: &mail.Address{
					Name:    "Test User",
					Address: "testuser@example.com",
				},
				ImapServer: &EmailServer{Host: "imap.example.com", Port: 993},
				SmtpServer: &EmailServer{Host: "smtp.example.com", Port: 587},
			},
			expected: "EmailAccount{Username:testuser, Password:********, Domain:example.com, EmailAddress:\"Test User\" <testuser@example.com>, ImapServer:imap.example.com:993, SmtpServer:smtp.example.com:587}",
		},
		{
			name: "EmailAccount with no password",
			account: &EmailAccount{
				Username: "testuser",
				Password: "",
				Domain:   "example.com",
				EmailAddress: &mail.Address{
					Name:    "Test User",
					Address: "testuser@example.com",
				},
				ImapServer: &EmailServer{Host: "imap.example.com", Port: 993},
				SmtpServer: &EmailServer{Host: "smtp.example.com", Port: 587},
			},
			expected: "EmailAccount{Username:testuser, Password:, Domain:example.com, EmailAddress:\"Test User\" <testuser@example.com>, ImapServer:imap.example.com:993, SmtpServer:smtp.example.com:587}",
		},
		{
			name: "EmailAccount with nil EmailAddress",
			account: &EmailAccount{
				Username:     "testuser",
				Password:     "testpassword",
				Domain:       "example.com",
				EmailAddress: nil,
				ImapServer:   &EmailServer{Host: "imap.example.com", Port: 993},
				SmtpServer:   &EmailServer{Host: "smtp.example.com", Port: 587},
			},
			expected: "EmailAccount{Username:testuser, Password:********, Domain:example.com, EmailAddress:, ImapServer:imap.example.com:993, SmtpServer:smtp.example.com:587}",
		},
		{
			name:     "Nil EmailAccount",
			account:  nil,
			expected: "<nil EmailAccount>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.account.String()
			if actual != tt.expected {
				t.Errorf("For %s, expected:\n%q\nGot:\n%q", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestEmailMessage_String(t *testing.T) {
	dateSent := time.Date(2026, time.January, 11, 10, 0, 0, 0, time.UTC)
	expectedDateSent := dateSent.Format(time.RFC3339)

	tests := []struct {
		name     string
		message  *EmailMessage
		expected string
	}{
		{
			name: "Full EmailMessage",
			message: &EmailMessage{
				From: &mail.Address{Name: "Sender", Address: "sender@example.com"},
				To: []*mail.Address{
					{Name: "Recipient 1", Address: "rec1@example.com"},
					{Name: "Recipient 2", Address: "rec2@example.com"},
				},
				Cc: []*mail.Address{
					{Name: "CC 1", Address: "cc1@example.com"},
				},
				Bcc: []*mail.Address{
					{Name: "BCC 1", Address: "bcc1@example.com"},
				},
				Subject:   "Test Subject",
				Body:      "Test Body",
				DateSent:  dateSent,
				MessageId: "message-id-123",
			},
			expected: fmt.Sprintf("EmailMessage{From:\"Sender\" <sender@example.com>, To:[\"Recipient 1\" <rec1@example.com> \"Recipient 2\" <rec2@example.com>], Cc:[\"CC 1\" <cc1@example.com>], Bcc:[\"BCC 1\" <bcc1@example.com>], Subject:\"Test Subject\", DateSent:%s, MessageId:message-id-123}", expectedDateSent),
		},
		{
			name: "EmailMessage with no recipients",
			message: &EmailMessage{
				From:      &mail.Address{Name: "Sender", Address: "sender@example.com"},
				Subject:   "No Recipients",
				Body:      "Test Body",
				DateSent:  dateSent,
				MessageId: "message-id-456",
			},
			expected: fmt.Sprintf("EmailMessage{From:\"Sender\" <sender@example.com>, To:[], Cc:[], Bcc:[], Subject:\"No Recipients\", DateSent:%s, MessageId:message-id-456}", expectedDateSent),
		},
		{
			name: "EmailMessage with nil From address",
			message: &EmailMessage{
				From: nil,
				To: []*mail.Address{
					{Name: "Recipient", Address: "rec@example.com"},
				},
				Subject:   "Nil From",
				Body:      "Test Body",
				DateSent:  dateSent,
				MessageId: "message-id-789",
			},
			expected: fmt.Sprintf("EmailMessage{From:, To:[\"Recipient\" <rec@example.com>], Cc:[], Bcc:[], Subject:\"Nil From\", DateSent:%s, MessageId:message-id-789}", expectedDateSent),
		},
		{
			name:     "Nil EmailMessage",
			message:  nil,
			expected: "<nil EmailMessage>",
		},
		{
			name: "EmailMessage with empty To/Cc/Bcc slices",
			message: &EmailMessage{
				From:      &mail.Address{Name: "Sender", Address: "sender@example.com"},
				To:        []*mail.Address{},
				Cc:        []*mail.Address{},
				Bcc:       []*mail.Address{},
				Subject:   "Empty Slices",
				Body:      "Test Body",
				DateSent:  dateSent,
				MessageId: "message-id-abc",
			},
			expected: fmt.Sprintf("EmailMessage{From:\"Sender\" <sender@example.com>, To:[], Cc:[], Bcc:[], Subject:\"Empty Slices\", DateSent:%s, MessageId:message-id-abc}", expectedDateSent),
		},
		{
			name: "EmailMessage with nil entries in recipient lists",
			message: &EmailMessage{
				From: &mail.Address{Name: "Sender", Address: "sender@example.com"},
				To: []*mail.Address{
					nil,
					{Name: "Recipient 2", Address: "rec2@example.com"},
				},
				Cc: []*mail.Address{
					{Name: "CC 1", Address: "cc1@example.com"},
					nil,
				},
				Bcc: []*mail.Address{
					nil,
					{Name: "BCC 2", Address: "bcc2@example.com"},
				},
				Subject:   "Nil Entries",
				Body:      "Test Body",
				DateSent:  dateSent,
				MessageId: "message-id-def",
			},
			expected: fmt.Sprintf("EmailMessage{From:\"Sender\" <sender@example.com>, To:[\"Recipient 2\" <rec2@example.com>], Cc:[\"CC 1\" <cc1@example.com>], Bcc:[\"BCC 2\" <bcc2@example.com>], Subject:\"Nil Entries\", DateSent:%s, MessageId:message-id-def}", expectedDateSent),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.message.String()
			if actual != tt.expected {
				t.Errorf("For %s, expected:\n%q\nGot:\n%q", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestEmailAccount_SecretFields(t *testing.T) {
	account := &EmailAccount{} // The values of the fields don't matter for this test.
	expected := []string{"Username", "Password", "Domain", "EmailAddress.Address"}
	actual := account.SecretFields()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("SecretFields() got %v, want %v", actual, expected)
	}
}
