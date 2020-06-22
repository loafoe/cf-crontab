package plugin

import (
	"code.cloudfoundry.org/cli/plugin"
	plugin_models "code.cloudfoundry.org/cli/plugin/models"
	"encoding/json"
	"fmt"
	"github.com/philips-labs/cf-crontab/crontab"
	"strings"
)

type CrontabServer struct {
	app *plugin_models.GetAppsModel
	connection plugin.CliConnection
}

func (c CrontabServer) GetToken() (string, error) {
	output, err := c.connection.AccessToken()
	if err != nil {
		return "", err
	}
	parsed := strings.Split(strings.Trim(output, "\n"), " ")
	if len(parsed) != 2 || parsed[0] != "bearer" || parsed[1] == "" {
		return "", errMissingOrInvalidToken
	}
	return parsed[1], nil
}

type appEnv struct {
	Environment map[string]string `json:"environment_json"`
}

func (c CrontabServer) GetSecret() (string, error) {
	url := fmt.Sprintf("/v2/apps/%s/env", c.app.Guid)
	out, err := c.connection.CliCommandWithoutTerminalOutput("curl", url)
	if err != nil {
		return "", err
	}
	data := strings.Join(out, "")
	err = json.Unmarshal([]byte(data), &env)
	if err != nil {
		return "", err
	}
	return env.Environment["CF_CRONTAB_SECRET"], nil
}

func (c CrontabServer) Host() (string, error) {
	for _, r := range c.app.Routes {
		if r.Domain.Name != "apps.internal" {
			return r.Host + "." + r.Domain.Name, nil
		}
	}
	return "", errMissingRoute
}

func (c CrontabServer) GetEntries() ([]*crontab.Task, error) {
	return nil, errMissingOrInvalidToken
}

func CrontabServerResolver(cliConnection plugin.CliConnection) (*CrontabServer, error) {
	apps, err := cliConnection.GetApps()
	if err != nil {
		return nil, err
	}
	for _, app := range apps {
		if app.Name == "cf-crontab" {
			return &CrontabServer{
				app: &app,
				connection: cliConnection,
			}, nil
		}
	}
	return nil, errNoDeployedCFCrontabFound
}