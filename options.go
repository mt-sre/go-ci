package regtest

type RegistryOption interface {
	ConfigureRegistry(*RegistryConfig)
}

type WithPort int

func (p WithPort) ConfigureRegistry(c *RegistryConfig) {
	c.Port = int(p)
}

type WithName string

func (n WithName) ConfigureRegistry(c *RegistryConfig) {
	c.Name = string(n)
}

type WithImage string

func (i WithImage) ConfigureRegistry(c *RegistryConfig) {
	c.Image = string(i)
}

type WithRuntime string

func (r WithRuntime) ConfigureRegistry(c *RegistryConfig) {
	c.Runtime = string(r)
}

type WithTLSEnabled bool

func (t WithTLSEnabled) ConfigureRegistry(c *RegistryConfig) {
	c.EnableTLS = bool(t)
}
