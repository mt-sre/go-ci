package gocmd

type WithBinPath string

func (w WithBinPath) ConfigureGoCmd(c *GoCmdConfig) {
	c.BinPath = string(w)
}

type WithWorkingDir string

func (w WithWorkingDir) ConfigureModule(c *ModuleConfig) {
	c.WorkingDir = string(w)
}
