package fycha

import (
	"net/http"
	"testing"
)

func TestHTMXSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		tableID     string
		wantStatus  int
		wantTrigger string
	}{
		{
			name:        "basic table ID",
			tableID:     "invoices-table",
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"invoices-table"}`,
		},
		{
			name:        "empty table ID",
			tableID:     "",
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":""}`,
		},
		{
			name:        "table ID with special characters",
			tableID:     "table-123_abc",
			wantStatus:  http.StatusOK,
			wantTrigger: `{"formSuccess":true,"refreshTable":"table-123_abc"}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := HTMXSuccess(tt.tableID)

			if result.StatusCode != tt.wantStatus {
				t.Errorf("StatusCode = %d, want %d", result.StatusCode, tt.wantStatus)
			}

			trigger, ok := result.Headers["HX-Trigger"]
			if !ok {
				t.Fatal("HX-Trigger header not set")
			}
			if trigger != tt.wantTrigger {
				t.Errorf("HX-Trigger = %q, want %q", trigger, tt.wantTrigger)
			}

			if result.Template != "" {
				t.Errorf("Template should be empty for header-only response, got %q", result.Template)
			}
		})
	}
}

func TestHTMXError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		message     string
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "validation error",
			message:     "Invalid amount",
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: "Invalid amount",
		},
		{
			name:        "empty message",
			message:     "",
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: "",
		},
		{
			name:        "long error message",
			message:     "The account balance is insufficient for this transaction. Please review your entries.",
			wantStatus:  http.StatusUnprocessableEntity,
			wantMessage: "The account balance is insufficient for this transaction. Please review your entries.",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := HTMXError(tt.message)

			if result.StatusCode != tt.wantStatus {
				t.Errorf("StatusCode = %d, want %d", result.StatusCode, tt.wantStatus)
			}

			msg, ok := result.Headers["HX-Error-Message"]
			if !ok {
				t.Fatal("HX-Error-Message header not set")
			}
			if msg != tt.wantMessage {
				t.Errorf("HX-Error-Message = %q, want %q", msg, tt.wantMessage)
			}

			if result.Template != "" {
				t.Errorf("Template should be empty for header-only response, got %q", result.Template)
			}
		})
	}
}
