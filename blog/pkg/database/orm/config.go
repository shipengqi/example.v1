package orm

type Config struct {
	dbType      string
	user        string
	password    string
	host        string
	name        string
	tablePrefix string
}

func (c *Config) DbType() string {
	return c.dbType
}

func (c *Config) SetDbType(dbType string) {
	c.dbType = dbType
}

func (c *Config) User() string {
	return c.user
}

func (c *Config) SetUser(user string) {
	c.user = user
}

func (c *Config) Password() string {
	return c.password
}

func (c *Config) SetPassword(password string) {
	c.password = password
}

func (c *Config) Host() string {
	return c.host
}

func (c *Config) SetHost(host string) {
	c.host = host
}

func (c *Config) Name() string {
	return c.name
}

func (c *Config) SetName(name string) {
	c.name = name
}

func (c *Config) TablePrefix() string {
	return c.tablePrefix
}

func (c *Config) SetTablePrefix(tablePrefix string) {
	c.tablePrefix = tablePrefix
}