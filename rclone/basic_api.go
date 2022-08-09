package rclone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// RcloneServer holds the information used to connect to a Clone rcd http server.
type RcloneServer struct {
	// URL used to connect to the http server.
	// End with '/'. For example, "http://localhost:5572/"
	Server string
	// "User" and "Password" are used for basic authorization.
	// Keep them empty if --rc-no-auth are set.
	User     string
	Password string
}

// NewRcloneServer returns a new ServerConfig.
func NewRcloneServer(server string, user string, password string) *RcloneServer {
	if !strings.HasSuffix(server, "/") {
		server = server + "/"
	}
	return &RcloneServer{Server: server, User: user, Password: password}
}

// Do calls the server's rcd api and returns the response body.
func (rs *RcloneServer) Do(method string, payload any) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json marshal request error: %v", err)
	}
	questURL := rs.Server + method
	req, err := http.NewRequest("POST", questURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if rs.User != "" {
		req.SetBasicAuth(rs.User, rs.Password)
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

// DoAsync submits a asynchronous job using Rclone rcd api and returns the JobID for the submitted job.
// "_async" will be set in request url as a param.
func (rs *RcloneServer) DoAsync(method string, payload any) (int, error) {
	method = addParamToURL(method, "_async=true")
	body, err := rs.Do(method, payload)
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
