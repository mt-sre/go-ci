package regtest

type RegistryOption interface {
	ConfigureRegistry(*RegistryConfig)
}

type WithDomain string

func (d WithDomain) ConfigureRegistry(c *RegistryConfig) {
	c.Domain = string(d)
}
