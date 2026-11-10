// Package pacmanconf
//
// A set of functions to parse the output
// of pacman-conf into go structs.
package pacmanconf

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"strings"
)

func makeEntries(rawConfig io.Reader) ([]entry, error) {
	sc := bufio.NewScanner(rawConfig)
	sc.Split(bufio.ScanLines)

	var (
		section string
		entries = make([]entry, 0)
	)

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case line == "":
			continue
		case strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]"):
			section = line[1 : len(line)-1]
		default:
			key, value, _ := strings.Cut(line, "=")
			entries = append(
				entries,
				entry{
					section: section,
					key:     normalize(key),
					value:   normalize(value),
				},
			)
		}
	}

	return entries, sc.Err()
}

func Parse(ctx context.Context, configPath string) (*Config, error) {
	rawConfig, err := readConfig(ctx, configPath)
	if err != nil {
		return nil, err
	}

	entries, err := makeEntries(bytes.NewReader(rawConfig))
	if err != nil {
		return nil, err
	}

	pf := func(e entry) bool { return e.section == "options" }
	oEntries, rEntries := partition(entries, pf)

	options, err := assembleOptions(oEntries)
	if err != nil {
		return nil, err
	}

	repositories, err := assembleRepositories(rEntries)
	if err != nil {
		return nil, err
	}

	return &Config{
		options:      options,
		repositories: repositories,
	}, nil
}
