package report

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func lynguaFile(rel string) string {
	return filepath.Join("..", "..", "..", "lyngua", "translations", "en", rel)
}

// The Go json tags must byte-match the Lyngua key tree (unknown keys disallowed).
func TestRecoveryLabelsMatchLynguaKeyTree(t *testing.T) {
	for _, rel := range []string{"general/recovery_report.json", "leasing/recovery_report.json"} {
		raw, err := os.ReadFile(lynguaFile(rel))
		if err != nil {
			t.Skipf("lyngua file not reachable (%s): %v", rel, err)
		}
		var doc struct {
			RecoveryReport RecoveryReportsLabels `json:"recovery_report"`
		}
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&doc); err != nil {
			t.Errorf("%s does not decode into RecoveryReportsLabels: %v", rel, err)
		}
	}
}

// C14: the English defaults equal the general Lyngua value of EVERY string field.
func TestDefaultRecoveryLabelsEqualGeneralLyngua(t *testing.T) {
	raw, err := os.ReadFile(lynguaFile("general/recovery_report.json"))
	if err != nil {
		t.Skipf("lyngua file not reachable: %v", err)
	}
	var doc struct {
		RecoveryReport RecoveryReportsLabels `json:"recovery_report"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	var walk func(path string, def, got reflect.Value)
	walk = func(path string, def, got reflect.Value) {
		if def.Kind() == reflect.Struct {
			for i := 0; i < def.NumField(); i++ {
				walk(path+"."+def.Type().Field(i).Name, def.Field(i), got.Field(i))
			}
			return
		}
		if def.String() == "" {
			t.Errorf("%s: empty Go default", path)
		}
		if def.String() != got.String() {
			t.Errorf("%s: default %q != general lyngua %q", path, def.String(), got.String())
		}
	}
	walk("Labels", reflect.ValueOf(DefaultRecoveryReportsLabels()), reflect.ValueOf(doc.RecoveryReport))
}

func TestRecoveryRouteMap(t *testing.T) {
	m := DefaultRecoveryReportsRoutes().RouteMap()
	want := map[string]string{
		"reports.recoverables_aging":                "/reports/recoverables-aging",
		"reports.recoverables_aging_export":         "/reports/recoverables-aging/export",
		"reports.cost_source_reconciliation":        "/reports/cost-source-reconciliation",
		"reports.cost_source_reconciliation_export": "/reports/cost-source-reconciliation/export",
	}
	if !reflect.DeepEqual(m, want) {
		t.Errorf("route map = %v", m)
	}
}
