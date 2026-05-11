// Package block — option types and toggles for fycha.Block.
//
// This file owns the public BlockOption surface and the package-private
// blockConfig flags. Adding a new optional fycha module means: (1) add a
// bool field to blockConfig, (2) add a `WithX() BlockOption` here, (3) add
// a `wantX()` accessor here, (4) check `cfg.wantX()` in block.go and wire
// the module. Nothing else in this file is load-bearing — it is a flat list
// by design so a reader can scan every option in one screen.
package block

// BlockOption enables specific fycha sub-modules within Block().
type BlockOption func(*blockConfig)

type blockConfig struct {
	enableAll bool
	reports   bool
	asset     bool
	ledger    bool
	loans     bool
	equity    bool
	payroll   bool
	cash      bool
	expenses  bool
	financial bool
	// Tax integration modules (Phase 2)
	taxRate                bool
	forexRate              bool
	withholdingCertificate bool
	// assetDepreciationRunURL is the resolved run-detail URL template plumbed into
	// the Surface A drawer so toast links point to the correct run-detail page.
	// Set via WithAssetDepreciationRunURL (Wave 2 hard requirement).
	assetDepreciationRunURL string
	// useCases is the typed use-case container supplied via WithUseCases().
	useCases *UseCases
}

// WithUseCases supplies the typed use-case closures to Block().
// Required: Block() returns an error if this option is not provided.
// Service-admin constructs the *UseCases via an adapter function that
// bridges espyna's consumer container to fycha's typed shape.
func WithUseCases(uc *UseCases) BlockOption {
	return func(c *blockConfig) { c.useCases = uc }
}

func WithReports() BlockOption   { return func(c *blockConfig) { c.reports = true } }
func WithAsset() BlockOption     { return func(c *blockConfig) { c.asset = true } }
func WithLedger() BlockOption    { return func(c *blockConfig) { c.ledger = true } }
func WithLoans() BlockOption     { return func(c *blockConfig) { c.loans = true } }
func WithEquity() BlockOption    { return func(c *blockConfig) { c.equity = true } }
func WithPayroll() BlockOption   { return func(c *blockConfig) { c.payroll = true } }
func WithCash() BlockOption      { return func(c *blockConfig) { c.cash = true } }
func WithExpenses() BlockOption  { return func(c *blockConfig) { c.expenses = true } }
func WithFinancial() BlockOption { return func(c *blockConfig) { c.financial = true } }
func WithTaxRate() BlockOption   { return func(c *blockConfig) { c.taxRate = true } }
func WithForexRate() BlockOption { return func(c *blockConfig) { c.forexRate = true } }
func WithWithholdingCertificate() BlockOption {
	return func(c *blockConfig) { c.withholdingCertificate = true }
}

// WithAssetDepreciationRunURL injects the run-detail URL into the block so
// the Surface A drawer can include a resolved link in its success toast payload.
// Wave 2 hard requirement — must be called before routes register.
//
// Example:
//
//	fychablock.Block(
//	    fychablock.WithAsset(),
//	    fychablock.WithAssetDepreciationRunURL(fycha.DepreciationRunDetailURL),
//	)
func WithAssetDepreciationRunURL(url string) BlockOption {
	return func(c *blockConfig) { c.assetDepreciationRunURL = url }
}

func (c *blockConfig) wantReports() bool   { return c.enableAll || c.reports }
func (c *blockConfig) wantAsset() bool     { return c.enableAll || c.asset }
func (c *blockConfig) wantLedger() bool    { return c.enableAll || c.ledger }
func (c *blockConfig) wantLoans() bool     { return c.enableAll || c.loans }
func (c *blockConfig) wantEquity() bool    { return c.enableAll || c.equity }
func (c *blockConfig) wantPayroll() bool   { return c.enableAll || c.payroll }
func (c *blockConfig) wantCash() bool      { return c.enableAll || c.cash }
func (c *blockConfig) wantExpenses() bool  { return c.enableAll || c.expenses }
func (c *blockConfig) wantFinancial() bool { return c.enableAll || c.financial }
func (c *blockConfig) wantTaxRate() bool   { return c.enableAll || c.taxRate }
func (c *blockConfig) wantForexRate() bool { return c.enableAll || c.forexRate }
func (c *blockConfig) wantWithholdingCertificate() bool {
	return c.enableAll || c.withholdingCertificate
}
