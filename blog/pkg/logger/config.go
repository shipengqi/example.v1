package logger

type Config struct {
	level  string
	output string
	prefix string
}

func (c *Config) Prefix() string {
	return c.prefix
}

func (c *Config) SetPrefix(prefix string) {
	c.prefix = prefix
}

func (c *Config) Output() string {
	return c.output
}

func (c *Config) SetOutput(output string) {
	c.output = output
}

func (c *Config) Level() string {
	return c.level
}

func (c *Config) SetLevel(level string) {
	c.level = level
}