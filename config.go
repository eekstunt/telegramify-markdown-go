package tgmd

// Option configures the converter.
type Option func(*config)

type config struct {
	headingSymbols  [6]string
	checkedMarker   string
	uncheckedMarker string
	bulletMarker    string
	maxMessageLen   int
}

func defaultConfig() config {
	return config{
		headingSymbols:  [6]string{"📌", "✏", "📚", "🔖", "", ""},
		checkedMarker:   "✅",
		uncheckedMarker: "☐",
		bulletMarker:    "⦁",
		maxMessageLen:   4096,
	}
}

func applyOptions(opts []Option) config {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}

// WithHeadingSymbols sets the emoji prefix for each heading level (h1 through h6).
func WithHeadingSymbols(symbols [6]string) Option {
	return func(c *config) { c.headingSymbols = symbols }
}

// WithTaskMarkers sets the unicode markers for checked/unchecked task list items.
func WithTaskMarkers(checked, unchecked string) Option {
	return func(c *config) {
		c.checkedMarker = checked
		c.uncheckedMarker = unchecked
	}
}

// WithBulletMarker sets the bullet character for unordered lists.
func WithBulletMarker(s string) Option {
	return func(c *config) { c.bulletMarker = s }
}

// WithMaxMessageLen sets the maximum message length in UTF-16 code units for splitting.
func WithMaxMessageLen(n int) Option {
	return func(c *config) { c.maxMessageLen = n }
}
