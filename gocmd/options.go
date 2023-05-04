package gocmd

type WithBinPath string

func (w WithBinPath) ConfigureGoCmd(c *GoCmdConfig) {
	c.BinPath = string(w)
}

type WithWorkingDir string

func (w WithWorkingDir) ConfigureModule(c *ModuleConfig) {
	c.WorkingDir = string(w)
}

type tidyOptionFunc func(*TidyConfig)

type tidyOption struct {
    f tidyOptionFunc
}

func (t *tidyOption) ConfigureTidy(c *TidyConfig) {
    t.f(c)
}

func WithGoVersion(version string) TidyOption {
    return &tidyOption{
        f: func(c *TidyConfig) {
            c.GoVersion = version
        },
    }
}

func WithBinWorkingDir(dir string) TidyOption {
    return &tidyOption{
        f: func(c *TidyConfig) {
            c.WorkingDir = dir
        },
    }
}
