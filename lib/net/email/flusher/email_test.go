package flusher

import (
	"github.com/emersion/go-message/mail"
	"github.com/walterjwhite/go-code/lib/net/email"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestEmailFlusher_Flush(t *testing.T) {
	flusher1 := (*EmailFlusher)(nil)
	err := flusher1.Flush([]byte("test"))
	if err == nil {
		t.Errorf("Expected error for nil flusher, got nil")
	}

	flusher2 := &EmailFlusher{}
	err = flusher2.Flush([]byte("test"))
	if err == nil {
		t.Errorf("Expected error for nil account, got nil")
	}

	account := &email.EmailAccount{
		EmailAddress: &mail.Address{Address: "test@example.com"},
	}
	flusher3 := &EmailFlusher{
		Account:         account,
		SubjectTemplate: "Test Subject",
		Recipients:      []*mail.Address{{Address: "recipient@example.com"}},
	}
	err = flusher3.Flush([]byte{})
	if err == nil {
		t.Errorf("Expected error for empty buffer, got nil")
	}

	flusher4 := &EmailFlusher{
		Account:         account,
		SubjectTemplate: "Test Subject",
		Recipients:      []*mail.Address{},
	}
	err = flusher4.Flush([]byte("test"))
	if err == nil {
		t.Errorf("Expected error for no recipients, got nil")
	}

	/*
		account := &email.EmailAccount{
			EmailAddress: &mail.Address{Address: "test@example.com"},
		}
		flusher5 := &EmailFlusher{
			Account:         account,
			SubjectTemplate: "Test Subject",
			Recipients:      []*mail.Address{{Address: "recipient@example.com"}},
		}
		err = flusher5.Flush([]byte("test body"))
		if err != nil {
			t.Errorf("Expected no error for successful flush (assuming write.Send succeeds), got %v", err)
		}
	*/
}

func TestEmailFlusher_toMessage(t *testing.T) {
	account := &email.EmailAccount{
		EmailAddress: &mail.Address{Address: "test@example.com"},
	}
	recipient := &mail.Address{Address: "recipient@example.com"}
	flusher := &EmailFlusher{
		Account:         account,
		SubjectTemplate: "Test Subject",
		Recipients:      []*mail.Address{recipient},
	}
	body := []byte("test body")

	actualMessage := flusher.toMessage(body)

	expectedFrom := account.EmailAddress
	expectedTo := []*mail.Address{recipient}
	expectedBody := "test body"
	expectedSubjectPrefix := "Test Subject - "

	if actualMessage.From.Address != expectedFrom.Address {
		t.Errorf("Expected From: %s, got %s", expectedFrom.Address, actualMessage.From.Address)
	}
	if !reflect.DeepEqual(actualMessage.To, expectedTo) {
		t.Errorf("Expected To: %v, got %v", expectedTo, actualMessage.To)
	}
	if actualMessage.Body != expectedBody {
		t.Errorf("Expected Body: %s, got %s", expectedBody, actualMessage.Body)
	}

	if !strings.HasPrefix(actualMessage.Subject, expectedSubjectPrefix) {
		t.Errorf("Subject does not have expected prefix. Expected '%s', Got: '%s'", expectedSubjectPrefix, actualMessage.Subject)
	}

	subjectParts := strings.Split(actualMessage.Subject, " - ")
	if len(subjectParts) != 3 {
		t.Errorf("Expected subject to split into 3 parts by ' - ', got %d", len(subjectParts))
	} else {
		_, err := time.Parse(time.RFC3339, subjectParts[1])
		if err != nil {
			t.Errorf("Could not parse subject timestamp '%s': %v", subjectParts[1], err)
		}
		secureID := subjectParts[2]
		if len(secureID) != 32 {
			t.Errorf("Expected secure ID to be 32 chars, got %d", len(secureID))
		}
	}
}

func TestGenerateSecureID(t *testing.T) {
	id1 := generateSecureID()
	id2 := generateSecureID()

	if len(id1) != 32 {
		t.Errorf("Expected secure ID to be 32 chars, got %d", len(id1))
	}

	if id1 == id2 {
		t.Errorf("Expected unique secure IDs, got duplicate: %s", id1)
	}
}
