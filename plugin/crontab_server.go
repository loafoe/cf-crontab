package plugin

import (
	"code.cloudfoundry.org/cli/plugin"
	plugin_models "code.cloudfoundry.org/cli/plugin/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/philips-labs/cf-crontab/crontab"
	signer "github.com/philips-software/go-hsdp-signer"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CrontabServer struct {
	app *plugin_models.GetAppsModel
	connection plugin.CliConnection
}

func (c CrontabServer) ServerRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	secret, err := c.GetSecret()
	if err != nil {
		return nil, err
	}
	host, err := c.Host()
	if err != nil {
		return nil, err
	}
	s, err := signer.New(crontab.SharedKey, secret)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse("https://"+host+endpoint)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(context.Background(), method, endpoint, body)
	req.URL = u
	req.Host = u.Host
	req.Proto = "HTTP/1.1"
	req.ProtoMajor = 1
	req.ProtoMinor = 1
	req.Header = make(http.Header)
	req.Header.Set("Content-Type", "application/json")

	err = s.SignRequest(req)
	if err != nil {
		return nil, err
	}
	return  http.DefaultClient.Do(req)
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
	envURL := fmt.Sprintf("/v2/apps/%s/env", c.app.Guid)
	out, err := c.connection.CliCommandWithoutTerminalOutput("curl", envURL)
	if err != nil {
		return "", err
	}
	var env appEnv
	data := strings.Join(out, "")
	err = json.Unmarshal([]byte(data), &env)
	if err != nil {
		return "", err
	}
	return env.Environment[crontab.EnvironmentSecret], nil
}

func (c CrontabServer) Host() (string, error) {
	for _, r := range c.app.Routes {
		if r.Domain.Name != crontab.InternalDomain {
			return r.Host + "." + r.Domain.Name, nil
		}
	}
	return "", errMissingRoute
}

func (c CrontabServer) GetEntries() ([]*crontab.Task, error) {
	resp, err := c.ServerRequest("GET", "/entries", nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", resp)
	return nil, errMissingOrInvalidToken
}

func CrontabServerResolver(cliConnection plugin.CliConnection) (*CrontabServer, error) {
	apps, err := cliConnection.GetApps()
	if err != nil {
		return nil, err
	}
	for _, app := range apps {
		if app.Name == crontab.DefaultAppName {
			return &CrontabServer{
				app: &app,
				connection: cliConnection,
			}, nil
		}
	}
	return nil, errNoDeployedCFCrontabFound
}