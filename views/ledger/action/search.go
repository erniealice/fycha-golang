package action

// search.go provides a JSON search handler for the account auto-complete widget
// in the journal entry form.
//
// Route: GET /action/ledger/accounts/search?q=<term>
// Returns: [{"value":"<id>","label":"<code> — <name>"}, ...]

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
)

const accountSearchLimit = 15

// AccountSearchOption is the JSON shape returned by the account search handler.
type AccountSearchOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// AccountSearchDeps holds dependencies for the account search handler.
type AccountSearchDeps struct {
	// GetAccountListPageData is used to fetch the full account list for filtering.
	// The search filters in Go after fetching all active accounts.
	GetAccountListPageData func(ctx context.Context, req *accountpb.GetAccountListPageDataRequest) (*accountpb.GetAccountListPageDataResponse, error)
}

// NewSearchAccountsHandler returns an http.HandlerFunc that searches accounts by
// code or name and returns JSON results for the auto-complete component.
//
// Query parameter: ?q=<search term>
// Response: [{"value":"<id>","label":"<code> — <name>"}, ...]
func NewSearchAccountsHandler(deps *AccountSearchDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		queryLower := strings.ToLower(query)

		results := searchAccounts(ctx, deps, queryLower)
		writeAccountJSON(w, results)
	}
}

// searchAccounts fetches all active accounts and filters by query.
func searchAccounts(ctx context.Context, deps *AccountSearchDeps, queryLower string) []AccountSearchOption {
	if deps == nil || deps.GetAccountListPageData == nil {
		return []AccountSearchOption{}
	}

	resp, err := deps.GetAccountListPageData(ctx, &accountpb.GetAccountListPageDataRequest{})
	if err != nil {
		log.Printf("account search: GetAccountListPageData error: %v", err)
		return []AccountSearchOption{}
	}
	if resp == nil || !resp.GetSuccess() {
		return []AccountSearchOption{}
	}

	var results []AccountSearchOption
	for _, a := range resp.GetAccountList() {
		if !a.GetActive() {
			continue
		}

		code := a.GetCode()
		name := a.GetName()

		if queryLower != "" {
			codeLower := strings.ToLower(code)
			nameLower := strings.ToLower(name)
			if !strings.Contains(codeLower, queryLower) && !strings.Contains(nameLower, queryLower) {
				continue
			}
		}

		label := code + " \u2014 " + name // "1110 — Cash on Hand"
		results = append(results, AccountSearchOption{
			Value: a.GetId(),
			Label: label,
		})

		if len(results) >= accountSearchLimit {
			break
		}
	}

	if results == nil {
		return []AccountSearchOption{}
	}
	return results
}

// writeAccountJSON marshals data as JSON and writes it to the response writer.
func writeAccountJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("account search: failed to encode JSON response: %v", err)
	}
}
