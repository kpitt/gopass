package config

// Pre1127 is a pre-1.12.7 config.
type Pre1127 struct {
	AutoClip      bool              `yaml:"autoclip"`      // decide whether passwords are automatically copied or not.
	AutoImport    bool              `yaml:"autoimport"`    // import missing public keys w/o asking.
	ClipTimeout   int               `yaml:"cliptimeout"`   // clear clipboard after seconds.
	ExportKeys    bool              `yaml:"exportkeys"`    // automatically export public keys of all recipients.
	NoColor       bool              `yaml:"nocolor"`       // do not use color when outputing text.
	NoPager       bool              `yaml:"nopager"`       // do not invoke a pager to display long lists.
	Notifications bool              `yaml:"notifications"` // enable desktop notifications.
	Parsing       bool              `yaml:"parsing"`       // allows to switch off all output parsing.
	Path          string            `yaml:"path"`
	SafeContent   bool              `yaml:"safecontent"` // avoid showing passwords in terminal.
	Mounts        map[string]string `yaml:"mounts"`

	ConfigPath string `yaml:"-"`

	// Catches all undefined files and must be empty after parsing.
	XXX map[string]any `yaml:",inline"`
}

// Config converts the Pre1127 config to the current config struct.
func (c *Pre1127) Config() *Config {
	cfg := &Config{
		AutoClip:      c.AutoClip,
		AutoImport:    c.AutoImport,
		ClipTimeout:   c.ClipTimeout,
		ExportKeys:    c.ExportKeys,
		NoPager:       c.NoPager,
		Notifications: c.Notifications,
		Parsing:       c.Parsing,
		Path:          c.Path,
		SafeContent:   c.SafeContent,
		Mounts:        make(map[string]string, len(c.Mounts)),
	}

	for k, v := range c.Mounts {
		cfg.Mounts[k] = v
	}

	return cfg
}

// CheckOverflow implements configer.
func (c *Pre1127) CheckOverflow() error {
	return checkOverflow(c.XXX)
}
