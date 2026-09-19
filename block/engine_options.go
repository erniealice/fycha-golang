package block

// EngineOption configures optional capabilities independently of business type.
type EngineOption func(*engineConfig)

type engineConfig struct{ assetProductSelection bool }

// WithAssetProductSelection enables product choices in asset add/edit drawers.
// The zero/default value leaves the product selector unavailable.
func WithAssetProductSelection(enabled bool) EngineOption {
	return func(c *engineConfig) { c.assetProductSelection = enabled }
}
