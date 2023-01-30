package docx

type Config struct {
	value Option
}

func newConfig(options ...Option) *Config {
	c := &Config{}
	for _, option := range options {
		c.value = c.value | option
	}
	return c
}

func (c *Config) checkOption(option Option) bool {
	return c.value&option == option
}

type Option uint64

const (
	WithDefaultStyle = Option(1 << iota)
)
