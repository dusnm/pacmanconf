package pacmanconf

type entry struct {
	section, key, value string
}

type Options struct {
	RootDir                  string
	DBPath                   string
	CacheDir                 []string
	HookDir                  []string
	GPGDir                   string
	LogFile                  string
	HoldPkg                  []string
	IgnorePkg                []string
	IgnoreGroup              []string
	NoUpgrade                []string
	NoExtract                []string
	Architecture             []string
	XferCommand              string
	CleanMethod              []string
	SigLevel                 []string
	LocalFileSigLevel        []string
	RemoteFileSigLevel       []string
	ParallelDownloads        int
	DownloadUser             string
	UseSyslog                bool
	Color                    bool
	NoProgressBar            bool
	CheckSpace               bool
	VerbosePkgLists          bool
	DisableDownloadTimeout   bool
	DisableSandbox           bool
	DisableSandboxFilesystem bool
	DisableSandboxSyscalls   bool
	ILoveCandy               bool
}

type Repository struct {
	Name         string
	Servers      []string
	CacheServers []string
	SigLevel     []string
	Usage        []string
}

type Config struct {
	options      *Options
	repositories []*Repository
}

func (c *Config) Options() *Options {
	return c.options
}

func (c *Config) Repositories() []*Repository {
	return c.repositories
}

func newOptions() *Options {
	return &Options{
		CacheDir:           make([]string, 0),
		HookDir:            make([]string, 0),
		HoldPkg:            make([]string, 0),
		IgnorePkg:          make([]string, 0),
		IgnoreGroup:        make([]string, 0),
		NoUpgrade:          make([]string, 0),
		NoExtract:          make([]string, 0),
		Architecture:       make([]string, 0),
		CleanMethod:        make([]string, 0),
		SigLevel:           make([]string, 0),
		LocalFileSigLevel:  make([]string, 0),
		RemoteFileSigLevel: make([]string, 0),
	}
}

func newRepository() *Repository {
	return &Repository{
		Servers:      make([]string, 0),
		CacheServers: make([]string, 0),
		SigLevel:     make([]string, 0),
		Usage:        make([]string, 0),
	}
}
