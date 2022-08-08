package rclone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ServerConfig holds the information used to connect to a Clone rcd http server.
type ServerConfig struct {
	// URL used to connect to the http server.
	// End with '/'. For example, "http://localhost:5572/"
	Server string
	// "User" and "Password" are used for basic authorization.
	// Keep them empty if --rc-no-auth are set.
	User     string
	Password string
}

// NewServerConfig returns a new ServerConfig.
func NewServerConfig(server string, user string, password string) *ServerConfig {
	if !strings.HasSuffix(server, "/") {
		server = server + "/"
	}
	return &ServerConfig{Server: server, User: user, Password: password}
}

// Submit a asynchronous job to the server.
//
// 'path' is the rclone API, such as 'core/stats'.
// 'input' is the input which is going to be encoded to JSON format and put in the request body.
func (sc *ServerConfig) Do(path string, input any) ([]byte, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("json marshal request error: %v", err)
	}
	questURL := sc.Server + path
	req, err := http.NewRequest("POST", questURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if sc.User != "" {
		req.SetBasicAuth(sc.User, sc.Password)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %v", err)
		}
		return body, nil
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("server err %s: parse body %v", resp.Status, err)
		}
		errorRespons := ErrorRespons{}
		if err := json.Unmarshal(body, &errorRespons); err != nil {
			return nil, fmt.Errorf("server err %s: parse error json %v", resp.Status, err)
		}
		return nil, fmt.Errorf("server err %d %s", errorRespons.Status, errorRespons.Error)
	}
}

// DoAsync returns the JobID for the submitted job.
// "_async" will be set in request url as a param.
func (sc *ServerConfig) DoAsync(path string, input any) (int, error) {
	path = addParamToURL(path, "_async=true")
	body, err := sc.Do(path, input)
	if err != nil {
		return -1, err
	}
	respJSON := &JobRespons{}
	err = json.Unmarshal(body, respJSON)
	if err != nil {
		return -1, fmt.Errorf("failed to parse response %s: %v", removeNewlines(string(body)), err)
	}
	return respJSON.JobID, nil
}
