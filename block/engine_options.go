package block

// EngineOption configures optional capabilities independently of business type.
type EngineOption func(*engineConfig)

type engineConfig struct {
	assetProductSelection bool
	chargePolicies        bool
	documentSeries        bool
	recoveryReports       bool
}

// WithAssetProductSelection enables product choices in asset add/edit drawers.
// The zero/default value leaves the product selector unavailable.
func WithAssetProductSelection(enabled bool) EngineOption {
	return func(c *engineConfig) { c.assetProductSelection = enabled }
}

// WithChargePolicies mounts the Ledger › Settings › Charge Policies pages
// (list, policy detail, version page, drawers and actions). Off by default so
// apps that do not opt in keep byte-identical compositions (UI-CP-10: only
// leasing-admin enables it).
func WithChargePolicies(enabled bool) EngineOption {
	return func(c *engineConfig) { c.chargePolicies = enabled }
}

// WithDocumentSeries mounts the Ledger › Settings › Document Series pages (list,
// add/edit drawer, retire). Off by default so apps that do not opt in keep
// byte-identical compositions.
func WithDocumentSeries(enabled bool) EngineOption {
	return func(c *engineConfig) { c.documentSeries = enabled }
}

// WithRecoveryReports mounts the recoverables aging and cost-source
// reconciliation reports. Off by default (same opt-in rule as charge policies).
func WithRecoveryReports(enabled bool) EngineOption {
	return func(c *engineConfig) { c.recoveryReports = enabled }
}
