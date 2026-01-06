package daily_activity

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/emersion/go-message/mail"
	"github.com/jmoiron/sqlx"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/net/email"
)

func TestWith(t *testing.T) {
	c := With(nil, nil)
	if c == nil {
		t.Fatal("With returned nil")
	}
}

func TestFetchHTTPRequests_DBNil(t *testing.T) {
	c := &Conf{db: nil}
	_, _, err := c.fetchHTTPRequests()
	if err == nil || !strings.Contains(err.Error(), "db is nil") {
		t.Fatalf("expected db is nil error, got: %v", err)
	}
}

func TestFetchHTTPRequests_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer closeTestResource(db)
	sqlxdb := sqlx.NewDb(db, "sqlmock")
	mock.ExpectQuery("SELECT \\* FROM http_requests").WillReturnError(fmt.Errorf("boom"))

	c := &Conf{db: sqlxdb}
	_, _, err = c.fetchHTTPRequests()
	if err == nil || !strings.Contains(err.Error(), "query http_requests") {
		t.Fatalf("expected query http_requests error, got: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestFetchHTTPRequests_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer closeTestResource(db)
	sqlxdb := sqlx.NewDb(db, "sqlmock")

	cols := []string{"id", "url", "body", "created_at"}
	now := time.Now().Truncate(time.Second)
	rows := sqlmock.NewRows(cols).AddRow(int64(1), "http://x", []byte("b"), now)
	mock.ExpectQuery("SELECT \\* FROM http_requests").WillReturnRows(rows)

	c := &Conf{db: sqlxdb}
	gotCols, records, err := c.fetchHTTPRequests()
	if err != nil {
		t.Fatalf("fetchHTTPRequests error: %v", err)
	}
	if !reflect.DeepEqual(gotCols, cols) {
		t.Fatalf("cols mismatch: got %v want %v", gotCols, cols)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	rec := records[0]
	if rec["url"] != "http://x" {
		t.Fatalf("unexpected url: %v", rec["url"])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestConvertRecordToStrings(t *testing.T) {
	cols := []string{"a", "b", "c", "d", "e"}
	now := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	rec := map[string]interface{}{
		"a": []byte("hello"),
		"b": "world",
		"c": now,
		"d": nil,
		"e": 42,
	}
	out := convertRecordToStrings(cols, rec)
	if out[0] != "hello" {
		t.Fatalf("a: expected hello, got %q", out[0])
	}
	if out[1] != "world" {
		t.Fatalf("b: expected world, got %q", out[1])
	}
	if out[2] != now.Format(time.RFC3339) {
		t.Fatalf("c: expected %q, got %q", now.Format(time.RFC3339), out[2])
	}
	if out[3] != "" {
		t.Fatalf("d: expected empty string, got %q", out[3])
	}
	if out[4] != "42" {
		t.Fatalf("e: expected 42, got %q", out[4])
	}
}

func TestTruncateHTTPRequests_DBNil(t *testing.T) {
	c := &Conf{db: nil}
	err := c.truncateHTTPRequests()
	if err == nil || !strings.Contains(err.Error(), "db is nil") {
		t.Fatalf("expected db is nil error, got: %v", err)
	}
}

func TestTruncateHTTPRequests_SuccessAndError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer closeTestResource(db)

	sqlxdb := sqlx.NewDb(db, "sqlmock")

	mock.ExpectExec("DELETE FROM http_requests").WillReturnResult(sqlmock.NewResult(0, 1))
	c := &Conf{db: sqlxdb}
	if err := c.truncateHTTPRequests(); err != nil {
		t.Fatalf("truncate should succeed, got: %v", err)
	}

	mock.ExpectExec("DELETE FROM http_requests").WillReturnError(fmt.Errorf("exec fail"))
	if err := c.truncateHTTPRequests(); err == nil || !strings.Contains(err.Error(), "truncate http_requests") {
		t.Fatalf("expected truncate http_requests error, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGenerateCSVBody(t *testing.T) {
	cols := []string{"id", "name"}
	records := []map[string]interface{}{{"id": 1, "name": "alice"}}
	body, err := generateCSVBody(cols, records)
	if err != nil {
		t.Fatalf("generateCSVBody error: %v", err)
	}
	if !strings.Contains(body, "id,name") || !strings.Contains(body, "1,alice") {
		t.Fatalf("unexpected csv body: %q", body)
	}
}

func TestBuildEmailMessage(t *testing.T) {
	addr := &mail.Address{Address: "me@example.org"}
	acct := &email.EmailAccount{EmailAddress: addr}
	msg := buildEmailMessage(acct, "body text")
	if msg.From == nil {
		t.Fatalf("From is nil")
	}
	if msg.From.String() != addr.String() {
		t.Fatalf("From mismatch: %v", msg.From)
	}
	if len(msg.To) != 1 {
		t.Fatalf("To length mismatch: %v", msg.To)
	}
	if msg.To[0].String() != addr.String() && msg.To[0].String() != addr.Address {
		t.Fatalf("To mismatch: %v", msg.To)
	}
	expectedDate := time.Now().Format("2006/01/02")
	if !strings.Contains(msg.Subject, expectedDate) {
		t.Fatalf("subject %q does not contain date %q", msg.Subject, expectedDate)
	}
}

func closeTestResource(db *sql.DB) {
	logging.Warn(db.Close(), "closeTestResource.db.Close")
}
