// Package block — workspace-scoped reference helpers shared across domain wirings.
//
// This file holds small, stateless helpers that any wireXxxModule function
// (and Block() itself) can call. They take args, return values, and never
// mutate package-level state. Anything that grew here past ~150 lines should
// be considered for a more specific name (e.g. workspace_currency.go).
package block

import (
	"context"
	"os"

	consumer "github.com/erniealice/espyna-golang/consumer"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
)

// getDefaultWorkspaceID mirrors the entydad block convention: read from env,
// fall back to a fixed default ID. Used by helpers that need workspace-scoped
// reference data (e.g. functional_currency for money display).
func getDefaultWorkspaceID() string {
	if v := os.Getenv("DEFAULT_WORKSPACE_ID"); v != "" {
		return v
	}
	return "default-workspace"
}

// getFunctionalCurrency returns the workspace's functional_currency (ISO 4217)
// for use in money display strings. Returns empty string when the workspace use
// case is not wired or the read fails — types.FormatMoney handles empty
// currency by omitting the prefix, so the worst-case fallback is the bare
// number rather than a hardcoded peso glyph.
//
// Use this for money values that are denominated in workspace currency by
// definition (revaluation amounts, ledger postings, GL-side totals). For
// values that already carry a per-row Currency field (e.g. depreciation
// candidate periods), pass that field instead — assets in foreign-currency
// workspaces may differ from the workspace currency at the row level.
func getFunctionalCurrency(ctx context.Context, useCases *consumer.UseCases) string {
	if useCases == nil || useCases.Entity == nil || useCases.Entity.Workspace == nil ||
		useCases.Entity.Workspace.ReadWorkspace == nil {
		return ""
	}
	resp, err := useCases.Entity.Workspace.ReadWorkspace.Execute(ctx, &workspacepb.ReadWorkspaceRequest{
		Data: &workspacepb.Workspace{Id: getDefaultWorkspaceID()},
	})
	if err != nil || resp == nil {
		return ""
	}
	data := resp.GetData()
	if len(data) == 0 {
		return ""
	}
	return data[0].GetFunctionalCurrency()
}
