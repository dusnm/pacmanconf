package pacmanconf

import "fmt"

func assembleRepositories(entries []entry) ([]*Repository, error) {
	result := make([]*Repository, 0, 10)
	repositories := make(map[string]*Repository, 10)

	for _, entry := range entries {
		_, ok := repositories[entry.section]
		if !ok {
			repositories[entry.section] = newRepository()
		}
	}

	for _, e := range entries {
		repository := repositories[e.section]
		repository.Name = e.section
		fields := map[string]any{
			KeyServer:      &repository.Servers,
			KeyCacheServer: &repository.CacheServers,
			KeySigLevel:    &repository.SigLevel,
			KeyUsage:       &repository.Usage,
		}

		target, ok := fields[e.key]
		if !ok {
			// unknown key
			continue
		}

		if err := parseValueType(target, e.value); err != nil {
			return nil, fmt.Errorf("%w %s: %w", ErrParseValueType, e.key, err)
		}
	}

	appended := make(map[string]struct{}, len(repositories))
	for _, e := range entries {
		_, ok := appended[e.section]
		if !ok {
			appended[e.section] = struct{}{}
			result = append(result, repositories[e.section])
		}
	}

	return result, nil
}
