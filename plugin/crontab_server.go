package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	plugin_models "code.cloudfoundry.org/cli/plugin/models"
	"github.com/philips-labs/cf-crontab/crontab"
	"github.com/philips-labs/cf-crontab/server"
	signer "github.com/philips-software/go-hsdp-signer"
)

type CrontabServer struct {
	app        *plugin_models.GetAppsModel
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
	u, err := url.Parse("https://" + host + endpoint)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(context.Background(), method, endpoint, body)
	if err != nil {
		return nil, err
	}
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
	return http.DefaultClient.Do(req)
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

func (c CrontabServer) SaveCrontab() error {
	entries, err := c.GetEntries()
	if err != nil {
		return err
	}
	parts, err := crontab.EnvParts(entries)
	if err != nil {
		return err
	}
	for k, v := range parts {
		_, err := c.connection.CliCommandWithoutTerminalOutput("set-env", c.app.Name, k, v)
		if err != nil {
			return err
		}
	}
	return nil
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
	var entries []*crontab.Task
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &entries)
	return entries, err
}

func (c CrontabServer) DeleteEntry(index int) (bool, error) {
	resp, err := c.ServerRequest("DELETE", fmt.Sprintf("/entries/%d", index), nil)
	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, errUnexpectedResponse
	}
	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	}
	var errResponse server.ErrResponse
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(data, &errResponse)
	if err != nil {
		return false, err
	}
	return false, errors.New(errResponse.Message)
}

func (c CrontabServer) AddEntries(tasks []crontab.Task) ([]*crontab.Task, error) {
	data, err := json.Marshal(tasks)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(data)
	resp, err := c.ServerRequest("POST", "/entries", body)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		return nil, errUnexpectedResponse
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var entries []*crontab.Task
	err = json.Unmarshal(data, &entries)
	return entries, err
}

func CrontabServerResolver(cliConnection plugin.CliConnection) (*CrontabServer, error) {
	apps, err := cliConnection.GetApps()
	if err != nil {
		return nil, err
	}
	for _, app := range apps {
		if app.Name == crontab.DefaultAppName {
			return &CrontabServer{
				app:        &app,
				connection: cliConnection,
			}, nil
		}
	}
	return nil, errNoDeployedCFCrontabFound
}
