package config

import (
	"github.com/muonsoft/errors"
	"golang.org/x/mod/semver"
)

const Version = "v0.1"

func checkVersion(configVersion string) error {
	if configVersion == "" {
		configVersion = "(unknown)"
	}

	compare := semver.Compare(configVersion, Version)
	if compare < 0 {
		return errors.Errorf(
			"config version %s is outdated, please upgrade your config file to match application %s",
			configVersion, Version,
		)
	}
	if compare > 0 {
		return errors.Errorf(
			"application version %s is outdated (config requires %s), please update the application",
			Version, configVersion,
		)
	}

	return nil
}
