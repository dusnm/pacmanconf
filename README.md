# pacmanconf

[![License: LGPL v3](https://img.shields.io/badge/License-LGPL_v3-blue.svg)](https://github.com/dusnm/pacmanconf/blob/main/LICENSE)
[![Tests](https://github.com/dusnm/pacmanconf/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/dusnm/pacmanconf/actions/workflows/test.yml)

A go library that parses pacman.conf for use in `libalpm` and similar purposes.

## Requirements

- Arch Linux
- `go` >= 1.27
- `pacman-conf` (Included with any Arch Linux installation.)

## Installation
```shell
go get github.com/dusnm/pacmanconf@v1.0.0
```

## Basic usage
```go
package main

import (
    "context"
    "log"

    "github.com/dusnm/pacmanconf"
)

func main() {
    ctx := context.Background()
    confPath := "/etc/pacman.conf" // The default path, if not provided.
    cfg, err := pacmanconf.Parse(ctx, confPath)
    if err != nil {
        log.Fatal(err)
    }

    log.Println(cfg.Options().RootDir)
    log.Println(cfg.Repositories()[0].Servers)
}
```

## Data types

Everything from `pacman.conf` is provided in the returned `*pacmanconf.Config`:

```go
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
```

## Licensing

This library is free and open-source software.

Licensed under the terms of the GNU Lesser General Public License, version 3.
