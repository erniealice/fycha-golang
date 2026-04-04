package fycha

import (
	"testing"
	"time"
)

func TestParsePeriodPreset(t *testing.T) {
	t.Parallel()

	now := time.Now()
	year, month, _ := now.Date()
	loc := now.Location()

	tests := []struct {
		preset    string
		wantStart time.Time
		wantEnd   time.Time
		approxEnd bool // true = end ≈ now, so we only check it's recent
	}{
		{
			preset:    "thisMonth",
			wantStart: time.Date(year, month, 1, 0, 0, 0, 0, loc),
			approxEnd: true,
		},
		{
			preset:    "lastMonth",
			wantStart: time.Date(year, month-1, 1, 0, 0, 0, 0, loc),
			wantEnd:   time.Date(year, month-1, 1, 0, 0, 0, 0, loc).AddDate(0, 1, 0).Add(-time.Second),
		},
		{
			preset:    "thisQuarter",
			wantStart: time.Date(year, time.Month(((int(month)-1)/3)*3+1), 1, 0, 0, 0, 0, loc),
			approxEnd: true,
		},
		{
			preset:    "lastQuarter",
			wantStart: time.Date(year, time.Month(((int(month)-1)/3)*3+1)-3, 1, 0, 0, 0, 0, loc),
			wantEnd:   time.Date(year, time.Month(((int(month)-1)/3)*3+1), 1, 0, 0, 0, 0, loc).Add(-time.Second),
		},
		{
			preset:    "thisYear",
			wantStart: time.Date(year, 1, 1, 0, 0, 0, 0, loc),
			approxEnd: true,
		},
		{
			preset:    "lastYear",
			wantStart: time.Date(year-1, 1, 1, 0, 0, 0, 0, loc),
			wantEnd:   time.Date(year, 1, 1, 0, 0, 0, 0, loc).Add(-time.Second),
		},
		{
			// unknown preset falls through to default (thisMonth behavior)
			preset:    "unknownPreset",
			wantStart: time.Date(year, month, 1, 0, 0, 0, 0, loc),
			approxEnd: true,
		},
		{
			// empty string falls through to default
			preset:    "",
			wantStart: time.Date(year, month, 1, 0, 0, 0, 0, loc),
			approxEnd: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.preset+"_preset", func(t *testing.T) {
			t.Parallel()

			start, end := ParsePeriodPreset(tt.preset)

			if !start.Equal(tt.wantStart) {
				t.Errorf("start = %v, want %v", start, tt.wantStart)
			}

			if tt.approxEnd {
				// end should be approximately now (within 2 seconds)
				diff := time.Since(end)
				if diff < 0 {
					diff = -diff
				}
				if diff > 2*time.Second {
					t.Errorf("end = %v, expected approximately now (%v), diff = %v", end, time.Now(), diff)
				}
			} else {
				if !end.Equal(tt.wantEnd) {
					t.Errorf("end = %v, want %v", end, tt.wantEnd)
				}
			}
		})
	}
}

func TestParsePeriodPreset_AdversarialInputs(t *testing.T) {
	t.Parallel()

	now := time.Now()
	year, month, _ := now.Date()
	loc := now.Location()

	// All adversarial inputs should fall through to the default case (thisMonth behavior)
	wantStart := time.Date(year, month, 1, 0, 0, 0, 0, loc)

	tests := []struct {
		name   string
		preset string
	}{
		{name: "SQL injection SELECT", preset: "'; DROP TABLE users; --"},
		{name: "SQL injection UNION", preset: "thisMonth' UNION SELECT * FROM users --"},
		{name: "script tag", preset: "<script>alert('xss')</script>"},
		{name: "null byte", preset: "thisMonth\x00"},
		{name: "newline injection", preset: "thisMonth\nDROP TABLE"},
		{name: "very long string", preset: string(make([]byte, 10000))},
		{name: "unicode emoji", preset: "\U0001f4a9"},
		{name: "path traversal", preset: "../../../etc/passwd"},
		{name: "LDAP injection", preset: "*)(&"},
		{name: "tab character", preset: "\t"},
		{name: "carriage return", preset: "\r\n"},
		{name: "backslash", preset: "\\"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			start, end := ParsePeriodPreset(tt.preset)

			// All unknown presets should fall to default (thisMonth)
			if !start.Equal(wantStart) {
				t.Errorf("start = %v, want %v (thisMonth default)", start, wantStart)
			}

			// end should be approximately now
			diff := time.Since(end)
			if diff < 0 {
				diff = -diff
			}
			if diff > 2*time.Second {
				t.Errorf("end = %v, expected approximately now, diff = %v", end, diff)
			}
		})
	}
}

func TestActiveFilterCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		filter FilterState
		want   int
	}{
		{
			name:   "all defaults returns 0",
			filter: FilterState{ActivePreset: "thisMonth", GroupBy: "product"},
			want:   0,
		},
		{
			name:   "empty preset and groupBy returns 0",
			filter: FilterState{},
			want:   0,
		},
		{
			name:   "non-default preset only",
			filter: FilterState{ActivePreset: "lastMonth", GroupBy: "product"},
			want:   1,
		},
		{
			name:   "non-default groupBy only",
			filter: FilterState{ActivePreset: "thisMonth", GroupBy: "customer"},
			want:   1,
		},
		{
			name:   "both non-default",
			filter: FilterState{ActivePreset: "thisYear", GroupBy: "status"},
			want:   2,
		},
		{
			name:   "empty preset with non-default groupBy",
			filter: FilterState{ActivePreset: "", GroupBy: "vendor"},
			want:   1,
		},
		{
			name:   "empty groupBy with non-default preset",
			filter: FilterState{ActivePreset: "lastQuarter", GroupBy: ""},
			want:   1,
		},
		{
			name:   "default preset string with empty groupBy",
			filter: FilterState{ActivePreset: "thisMonth", GroupBy: ""},
			want:   0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ActiveFilterCount(tt.filter)
			if got != tt.want {
				t.Errorf("ActiveFilterCount(%+v) = %d, want %d", tt.filter, got, tt.want)
			}
		})
	}
}

func TestActiveFilterCount_WhitespaceAndEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		filter FilterState
		want   int
	}{
		{
			name:   "whitespace-only preset counts as non-default",
			filter: FilterState{ActivePreset: "   ", GroupBy: "product"},
			want:   1,
		},
		{
			name:   "whitespace-only groupBy counts as non-default",
			filter: FilterState{ActivePreset: "thisMonth", GroupBy: "   "},
			want:   1,
		},
		{
			name:   "both whitespace-only",
			filter: FilterState{ActivePreset: "   ", GroupBy: "   "},
			want:   2,
		},
		{
			name:   "tab character preset",
			filter: FilterState{ActivePreset: "\t", GroupBy: "product"},
			want:   1,
		},
		{
			name:   "newline character groupBy",
			filter: FilterState{ActivePreset: "thisMonth", GroupBy: "\n"},
			want:   1,
		},
		{
			name:   "very long preset string",
			filter: FilterState{ActivePreset: string(make([]byte, 1000)), GroupBy: "product"},
			want:   1, // non-empty, non-default
		},
		{
			name:   "SQL injection in preset",
			filter: FilterState{ActivePreset: "'; DROP TABLE --", GroupBy: "product"},
			want:   1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ActiveFilterCount(tt.filter)
			if got != tt.want {
				t.Errorf("ActiveFilterCount(%+v) = %d, want %d", tt.filter, got, tt.want)
			}
		})
	}
}

func TestDefaultPeriodPresets(t *testing.T) {
	t.Parallel()

	labels := PeriodLabels{
		ThisMonth:   "This Month",
		LastMonth:   "Last Month",
		ThisQuarter: "This Quarter",
		LastQuarter: "Last Quarter",
		ThisYear:    "This Year",
		LastYear:    "Last Year",
		Custom:      "Custom",
	}

	tests := []struct {
		name             string
		active           string
		wantLen          int
		wantSelected     string // value that should be selected
		wantNoneSelected bool   // if true, no option should be selected
	}{
		{
			name:         "thisMonth active",
			active:       "thisMonth",
			wantLen:      7,
			wantSelected: "thisMonth",
		},
		{
			name:         "lastYear active",
			active:       "lastYear",
			wantLen:      7,
			wantSelected: "lastYear",
		},
		{
			name:         "custom active",
			active:       "custom",
			wantLen:      7,
			wantSelected: "custom",
		},
		{
			name:             "no match selects nothing",
			active:           "nonexistent",
			wantLen:          7,
			wantNoneSelected: true,
		},
		{
			name:             "empty active selects nothing",
			active:           "",
			wantLen:          7,
			wantNoneSelected: true,
		},
	}

	expectedOrder := []struct {
		value string
		label string
	}{
		{"thisMonth", "This Month"},
		{"lastMonth", "Last Month"},
		{"thisQuarter", "This Quarter"},
		{"lastQuarter", "Last Quarter"},
		{"thisYear", "This Year"},
		{"lastYear", "Last Year"},
		{"custom", "Custom"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := DefaultPeriodPresets(labels, tt.active)
			if len(opts) != tt.wantLen {
				t.Fatalf("len(opts) = %d, want %d", len(opts), tt.wantLen)
			}

			// Verify order and labels
			for i, expected := range expectedOrder {
				if opts[i].Value != expected.value {
					t.Errorf("opts[%d].Value = %q, want %q", i, opts[i].Value, expected.value)
				}
				if opts[i].Label != expected.label {
					t.Errorf("opts[%d].Label = %q, want %q", i, opts[i].Label, expected.label)
				}
			}

			// Verify selection
			selectedCount := 0
			for _, opt := range opts {
				if opt.Selected {
					selectedCount++
					if tt.wantNoneSelected {
						t.Errorf("option %q should not be selected", opt.Value)
					} else if opt.Value != tt.wantSelected {
						t.Errorf("selected option = %q, want %q", opt.Value, tt.wantSelected)
					}
				}
			}

			if tt.wantNoneSelected && selectedCount != 0 {
				t.Errorf("expected no selected options, got %d", selectedCount)
			}
			if !tt.wantNoneSelected && selectedCount != 1 {
				t.Errorf("expected exactly 1 selected option, got %d", selectedCount)
			}
		})
	}
}
