package rclone

import (
	"encoding/json"
	"fmt"
)

// json struct:
// operations/copyurl
type ListremotesResponse struct {
	Remotes []string
}

func (sc *ServerConfig) ListRemotes() ([]string, error) {
	body, err := sc.Do("config/listremotes", nil)
	if err != nil {
		return nil, fmt.Errorf("listremotes: %v", err)
	}
	resp := ListremotesResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("listremotes parse resp body %s err: %v", removeNewlines(string(body)), err)
	}
	return resp.Remotes, nil
}
