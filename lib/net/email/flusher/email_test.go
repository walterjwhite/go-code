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
	if err != nil {
		t.Errorf("Expected no error for nil flusher, got %v", err)
	}

	flusher2 := &EmailFlusher{}
	err = flusher2.Flush([]byte("test"))
	if err != nil {
		t.Errorf("Expected no error for nil account, got %v", err)
	}

	/*
		account := &email.EmailAccount{
			EmailAddress: &mail.Address{Address: "test@example.com"},
		}
		flusher3 := &EmailFlusher{
			Account:         account,
			SubjectTemplate: "Test Subject",
		}
		err = flusher3.Flush([]byte("test body"))
		if err != nil {
			t.Errorf("Expected no error for successful flush (assuming write.Send succeeds), got %v", err)
		}
	*/
}

func TestEmailFlusher_toMessage(t *testing.T) {
	account := &email.EmailAccount{
		EmailAddress: &mail.Address{Address: "test@example.com"},
	}
	flusher := &EmailFlusher{
		Account:         account,
		SubjectTemplate: "Test Subject",
	}
	body := []byte("test body")

	actualMessage := flusher.toMessage(body)

	expectedFrom := account.EmailAddress
	expectedTo := []*mail.Address{account.EmailAddress}
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

	subjectParts := strings.Split(actualMessage.Subject, expectedSubjectPrefix)
	if len(subjectParts) != 2 {
		t.Errorf("Expected subject to split into 2 parts by '%s', got %d", expectedSubjectPrefix, len(subjectParts))
	} else {
		_, err := time.Parse(time.RFC3339, subjectParts[1])
		if err != nil {
			t.Errorf("Could not parse subject timestamp '%s': %v", subjectParts[1], err)
		}
	}
}
