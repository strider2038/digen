package app

import (
	"github.com/muonsoft/errors"
	"golang.org/x/mod/semver"
)

const Version = "v0.1"

func checkVersion(configVersion, appVersion string) error {
	compare := semver.Compare(configVersion, appVersion)
	if compare < 0 {
		return errors.Errorf(
			"config version %s is outdated, please upgrade your config file to match application %s",
			configVersion, appVersion,
		)
	}
	if compare > 0 {
		return errors.Errorf(
			"application version %s is outdated (config requires %s), please update the application",
			appVersion, configVersion,
		)
	}

	return nil
}
