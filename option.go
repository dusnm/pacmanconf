package pacmanconf

import "fmt"

func assembleOptions(entries []entry) (*Options, error) {
	o := newOptions()

	fields := map[string]any{
		KeyRootDir:                  &o.RootDir,
		KeyDBPath:                   &o.DBPath,
		KeyCacheDir:                 &o.CacheDir,
		KeyHookDir:                  &o.HookDir,
		KeyGPGDir:                   &o.GPGDir,
		KeyLogFile:                  &o.LogFile,
		KeyHoldPkg:                  &o.HoldPkg,
		KeyIgnorePkg:                &o.IgnorePkg,
		KeyIgnoreGroup:              &o.IgnoreGroup,
		KeyNoUpgrade:                &o.NoUpgrade,
		KeyNoExtract:                &o.NoExtract,
		KeyArchitecture:             &o.Architecture,
		KeyXferCommand:              &o.XferCommand,
		KeyCleanMethod:              &o.CleanMethod,
		KeySigLevel:                 &o.SigLevel,
		KeyLocalFileSigLevel:        &o.LocalFileSigLevel,
		KeyRemoteFileSigLevel:       &o.RemoteFileSigLevel,
		KeyParallelDownloads:        &o.ParallelDownloads,
		KeyDownloadUser:             &o.DownloadUser,
		KeyUseSyslog:                &o.UseSyslog,
		KeyColor:                    &o.Color,
		KeyNoProgressBar:            &o.NoProgressBar,
		KeyCheckSpace:               &o.CheckSpace,
		KeyVerbosePkgLists:          &o.VerbosePkgLists,
		KeyDisableDownloadTimeout:   &o.DisableDownloadTimeout,
		KeyDisableSandbox:           &o.DisableSandbox,
		KeyDisableSandboxFilesystem: &o.DisableSandboxFilesystem,
		KeyDisableSandboxSyscalls:   &o.DisableSandboxSyscalls,
		KeyILoveCandy:               &o.ILoveCandy,
	}

	for _, e := range entries {
		target, ok := fields[e.key]
		if !ok {
			// unknown key
			continue
		}

		if err := parseValueType(target, e.value); err != nil {
			return nil, fmt.Errorf("%w %s: %w", ErrParseValueType, e.key, err)
		}
	}

	return o, nil
}
