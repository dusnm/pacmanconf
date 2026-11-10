package pacmanconf_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	conf "github.com/dusnm/pacmanconf"
	"github.com/stretchr/testify/require"
)

// These tests run the real pacman-conf, so they check the parser against the
// exact output format of the installed pacman. They are skipped elsewhere.
func requirePacmanConf(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("pacman-conf"); err != nil {
		t.Skip("pacman-conf not found in PATH")
	}
}

// newRoot returns a temp dir (symlinks resolved, because pacman-conf prints
// real paths) containing the directories the fixtures refer to.
// pacman-conf refuses paths that do not exist.
func newRoot(t *testing.T) string {
	t.Helper()
	assert := require.New(t)

	root, err := filepath.EvalSymlinks(t.TempDir())
	assert.NoError(err)
	for _, d := range []string{"db", "cache1", "cache2", "gpg"} {
		assert.NoError(os.Mkdir(filepath.Join(root, d), 0o755))
	}
	return root
}

// writeConf writes content to root/name, replacing @ROOT@ with root.
func writeConf(t *testing.T, root, name, content string) string {
	t.Helper()
	path := filepath.Join(root, name)
	content = strings.ReplaceAll(content, "@ROOT@", root)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func mustParse(t *testing.T, path string) *conf.Config {
	t.Helper()
	assert := require.New(t)

	cfg, err := conf.Parse(context.Background(), path)
	assert.NoError(err)
	assert.NotNil(cfg)
	assert.NotNil(cfg.Options())
	return cfg
}

// clean normalizes paths: pacman-conf resolves them with realpath(3), so
// trailing slashes differ between keys.
func clean(paths ...string) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		out[i] = filepath.Clean(p)
	}
	return out
}

func TestParseOptions(t *testing.T) {
	requirePacmanConf(t)
	assert := require.New(t)

	root := newRoot(t)
	path := writeConf(t, root, "pacman.conf", `
[options]
RootDir     = @ROOT@
DBPath      = @ROOT@/db/
CacheDir    = @ROOT@/cache1/
CacheDir    = @ROOT@/cache2/
GPGDir      = @ROOT@/gpg/
HoldPkg     = pacman glibc
IgnorePkg   = foo bar
IgnoreGroup = grp
NoUpgrade   = etc/a.conf
NoExtract   = usr/share/doc/*
Architecture = x86_64 x86_64_v3
XferCommand = /usr/bin/curl -o %o %u
CleanMethod = KeepInstalled KeepCurrent
SigLevel    = Required DatabaseOptional
LocalFileSigLevel  = Optional
RemoteFileSigLevel = Required
ParallelDownloads  = 7
DownloadUser = alpm
UseSyslog
Color
NoProgressBar
CheckSpace
VerbosePkgLists
DisableDownloadTimeout
DisableSandboxFilesystem
DisableSandboxSyscalls
ILoveCandy
`)
	o := mustParse(t, path).Options()

	// Single-value paths
	assert.Equal(filepath.Clean(root), filepath.Clean(o.RootDir), "RootDir")
	assert.Equal(root+"/db", filepath.Clean(o.DBPath), "DBPath")
	assert.Equal(root+"/gpg", filepath.Clean(o.GPGDir), "GPGDir")
	assert.Equal(root+"/var/log/pacman.log", o.LogFile, "LogFile defaults to under RootDir")

	// Repeated keys keep their order
	assert.Equal(clean(root+"/cache1", root+"/cache2"), clean(o.CacheDir...), "CacheDir")

	// Multi-word values are split into one entry each
	assert.Equal([]string{"pacman", "glibc"}, o.HoldPkg, "HoldPkg")
	assert.Equal([]string{"foo", "bar"}, o.IgnorePkg, "IgnorePkg")
	assert.Equal([]string{"grp"}, o.IgnoreGroup, "IgnoreGroup")
	assert.Equal([]string{"etc/a.conf"}, o.NoUpgrade, "NoUpgrade")
	assert.Equal([]string{"usr/share/doc/*"}, o.NoExtract, "NoExtract")
	assert.Equal([]string{"x86_64", "x86_64_v3"}, o.Architecture, "Architecture")
	assert.Equal([]string{"KeepInstalled", "KeepCurrent"}, o.CleanMethod, "CleanMethod")

	// A value containing spaces stays whole
	assert.Equal("/usr/bin/curl -o %o %u", o.XferCommand, "XferCommand")

	// SigLevel shorthands are expanded by pacman-conf
	assert.Equal([]string{
		"PackageRequired", "PackageTrustedOnly", "DatabaseOptional", "DatabaseTrustedOnly",
	}, o.SigLevel, "SigLevel")
	assert.Equal([]string{"PackageOptional", "PackageTrustedOnly"}, o.LocalFileSigLevel, "LocalFileSigLevel")
	assert.Equal([]string{"PackageRequired", "PackageTrustedOnly"}, o.RemoteFileSigLevel, "RemoteFileSigLevel")

	// Scalars
	assert.Equal(7, o.ParallelDownloads, "ParallelDownloads")
	assert.Equal("alpm", o.DownloadUser, "DownloadUser")

	// Flags
	assert.True(o.UseSyslog, "UseSyslog")
	assert.True(o.Color, "Color")
	assert.True(o.NoProgressBar, "NoProgressBar")
	assert.True(o.CheckSpace, "CheckSpace")
	assert.True(o.VerbosePkgLists, "VerbosePkgLists")
	assert.True(o.DisableDownloadTimeout, "DisableDownloadTimeout")
	assert.True(o.DisableSandboxFilesystem, "DisableSandboxFilesystem")
	assert.True(o.DisableSandboxSyscalls, "DisableSandboxSyscalls")
	assert.True(o.ILoveCandy, "ILoveCandy")
}

func TestParseMinimalUsesPacmanDefaults(t *testing.T) {
	requirePacmanConf(t)
	assert := require.New(t)

	root := newRoot(t)
	cfg := mustParse(t, writeConf(t, root, "pacman.conf", "[options]\n"))
	o := cfg.Options()

	assert.Empty(cfg.Repositories())
	assert.Equal(1, o.ParallelDownloads, "ParallelDownloads")
	assert.Equal([]string{"KeepInstalled"}, o.CleanMethod, "CleanMethod")
	assert.Equal([]string{
		"PackageRequired", "PackageTrustedOnly", "DatabaseRequired", "DatabaseTrustedOnly",
	}, o.SigLevel, "SigLevel")

	assert.False(o.Color, "Color")
	assert.False(o.CheckSpace, "CheckSpace")
	assert.False(o.ILoveCandy, "ILoveCandy")
	assert.False(o.DisableSandboxFilesystem, "DisableSandboxFilesystem")

	assert.Empty(o.IgnorePkg, "IgnorePkg")
	assert.Empty(o.DownloadUser, "DownloadUser")
}

// pacman-conf never prints DisableSandbox; it expands it into both specific flags.
func TestParseDisableSandboxExpands(t *testing.T) {
	requirePacmanConf(t)
	assert := require.New(t)

	root := newRoot(t)
	o := mustParse(t, writeConf(t, root, "pacman.conf", "[options]\nDisableSandbox\n")).Options()

	assert.True(o.DisableSandboxFilesystem, "DisableSandboxFilesystem")
	assert.True(o.DisableSandboxSyscalls, "DisableSandboxSyscalls")
}

func TestParseRepositories(t *testing.T) {
	requirePacmanConf(t)
	assert := require.New(t)

	root := newRoot(t)
	writeConf(t, root, "mirrorlist", `
Server = https://mirror.a/$repo/os/$arch
Server = https://mirror.b/$repo/os/$arch
`)
	path := writeConf(t, root, "pacman.conf", `
[options]
Architecture = x86_64

[core]
Include = @ROOT@/mirrorlist

[extra]
Include = @ROOT@/mirrorlist

[custom]
SigLevel    = Optional TrustAll
Usage       = Sync Search
CacheServer = https://cache.example/$repo
Server      = file:///srv/$repo
Server      = https://x.example/get?repo=$repo&a=b

[empty]
`)
	repos := mustParse(t, path).Repositories()

	// Section order is repository priority
	var names []string
	for _, r := range repos {
		names = append(names, r.Name)
	}
	assert.Equal([]string{"core", "extra", "custom", "empty"}, names, "repository order")
	core, extra, custom, empty := repos[0], repos[1], repos[2], repos[3]

	// Include and $repo/$arch are expanded; server order is kept
	assert.Equal([]string{
		"https://mirror.a/core/os/x86_64", "https://mirror.b/core/os/x86_64",
	}, core.Servers, "core.Servers")
	assert.Equal([]string{
		"https://mirror.a/extra/os/x86_64", "https://mirror.b/extra/os/x86_64",
	}, extra.Servers, "extra.Servers")

	assert.Empty(core.SigLevel, "core.SigLevel")
	assert.Equal([]string{
		"PackageOptional", "PackageTrustAll", "DatabaseOptional", "DatabaseTrustAll",
	}, custom.SigLevel, "custom.SigLevel")

	// Usage defaults to All; explicit values are split
	assert.Equal([]string{"All"}, core.Usage, "core.Usage")
	assert.Equal([]string{"Sync", "Search"}, custom.Usage, "custom.Usage")

	assert.Equal([]string{"https://cache.example/custom"}, custom.CacheServers, "custom.CacheServers")
	// Only the first '=' separates key and value
	assert.Equal([]string{
		"file:///srv/custom", "https://x.example/get?repo=custom&a=b",
	}, custom.Servers, "custom.Servers")

	// A repository with no servers is still listed
	assert.Empty(empty.Servers, "empty.Servers")
	assert.Empty(empty.CacheServers, "empty.CacheServers")
}

func TestParseIgnoresUnrecognizedDirectives(t *testing.T) {
	requirePacmanConf(t)
	assert := require.New(t)

	root := newRoot(t)
	// pacman-conf only warns about these (exit 0) and leaves them out of its output.
	path := writeConf(t, root, "pacman.conf", `
[options]
Color = yes
SomeFutureOption = 1
Bogu
CheckSpace
`)
	o := mustParse(t, path).Options()

	assert.False(o.Color, "Color = yes is invalid and must not enable Color")
	assert.True(o.CheckSpace, "CheckSpace after unknown lines must still be parsed")
}

func TestParseErrors(t *testing.T) {
	requirePacmanConf(t)
	root := newRoot(t)

	tests := []struct {
		name    string
		path    string
		wantMsg string // must appear in the error, i.e. pacman-conf's stderr is passed on
	}{
		{
			name:    "missing file",
			path:    filepath.Join(root, "does-not-exist.conf"),
			wantMsg: "could not be read",
		},
		{
			name:    "invalid ParallelDownloads",
			path:    writeConf(t, root, "bad-int.conf", "[options]\nParallelDownloads = abc\n"),
			wantMsg: "ParallelDownloads",
		},
		{
			name:    "nonexistent DBPath",
			path:    writeConf(t, root, "bad-path.conf", "[options]\nDBPath = @ROOT@/nope/\n"),
			wantMsg: "DBPath",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)

			cfg, err := conf.Parse(context.Background(), tt.path)
			assert.ErrorContains(err, tt.wantMsg)
			assert.Nil(cfg, "config must be nil on error")
		})
	}
}

func TestParseEmptyPathUsesDefault(t *testing.T) {
	requirePacmanConf(t)
	assert := require.New(t)

	if _, err := os.Stat("/etc/pacman.conf"); err != nil {
		t.Skip("/etc/pacman.conf not present")
	}

	got, err := conf.Parse(context.Background(), "")
	assert.NoError(err)
	assert.Equal(mustParse(t, "/etc/pacman.conf"), got)
}

func TestParseCanceledContext(t *testing.T) {
	requirePacmanConf(t)
	assert := require.New(t)

	root := newRoot(t)
	path := writeConf(t, root, "pacman.conf", "[options]\n")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := conf.Parse(ctx, path)
	assert.Error(err)
}
