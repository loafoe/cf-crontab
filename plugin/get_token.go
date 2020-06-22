package plugin

import (
	"code.cloudfoundry.org/cli/plugin"
	plugin_models "code.cloudfoundry.org/cli/plugin/models"
	"strings"
)

func GetToken(cliConnection plugin.CliConnection) (string, error) {
	output, err := cliConnection.AccessToken()
	if err != nil {
		return "", err
	}
	parsed := strings.Split(strings.Trim(output, "\n"), " ")
	if len(parsed) != 2 || parsed[0] != "bearer" || parsed[1] == "" {
		return "", errMissingOrInvalidToken
	}
	return parsed[1], nil
}

func CrontabServerResolver(cliConnection plugin.CliConnection) (*plugin_models.GetAppsModel, error) {
	apps, err := cliConnection.GetApps()
	if err != nil {
		return nil, err
	}
	for _, app := range apps {
		if app.Name == "cf-crontab" {
			return &app, nil
		}
	}
	return nil, errNoDeployedCFCrontabFound
}

